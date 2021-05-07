/*
Copyright 2020-2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package custom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"k8s.io/klog/v2"

	nfdv1alpha1 "sigs.k8s.io/node-feature-discovery/pkg/apis/nfd/v1alpha1"
	"sigs.k8s.io/node-feature-discovery/pkg/utils"
	"sigs.k8s.io/node-feature-discovery/source"
	"sigs.k8s.io/node-feature-discovery/source/custom/rules"
)

const Name = "custom"

type Domains nfdv1alpha1.MatchRule

// Custom Features Configurations
type MatchRule struct {
	Domains
	Legacy
}

// Legacy rules
type Legacy struct {
	PciID      *rules.PciIDRule      `json:"pciId,omitempty"`
	UsbID      *rules.UsbIDRule      `json:"usbId,omitempty"`
	LoadedKMod *rules.LoadedKModRule `json:"loadedKMod,omitempty"`
	CpuID      *rules.CpuIDRule      `json:"cpuId,omitempty"`
	Kconfig    *rules.KconfigRule    `json:"kConfig,omitempty"`
	Nodename   *rules.NodenameRule   `json:"nodename,omitempty"`
}

type FeatureSpec struct {
	Name    string      `json:"name"`
	Value   *string     `json:"value,omitempty"`
	MatchOn []MatchRule `json:"matchOn"`

	nameTemplate *template.Template
}

type config []FeatureSpec

// newDefaultConfig returns a new config with pre-populated defaults
func newDefaultConfig() *config {
	return &config{}
}

// customSource implements the LabelSource and ConfigurableSource interfaces.
type customSource struct {
	config *config
}

// Singleton source instance
var (
	src customSource
	_   source.LabelSource        = &src
	_   source.ConfigurableSource = &src
)

// Return name of the feature source
func (s *customSource) Name() string { return Name }

// NewConfig method of the LabelSource interface
func (s *customSource) NewConfig() source.Config { return newDefaultConfig() }

// GetConfig method of the LabelSource interface
func (s *customSource) GetConfig() source.Config { return s.config }

// SetConfig method of the LabelSource interface
func (s *customSource) SetConfig(c source.Config) {
	switch c.(type) {
	case *config:
	default:
		klog.Fatalf("invalid config type: %T", c)
	}

	// Parse template rules
	conf := c.(*config)
	for i, spec := range *conf {
		if strings.Contains(spec.Name, "{{") {
			(*conf)[i].nameTemplate = template.Must(template.New("").Option("missingkey=error").Parse(spec.Name))
		}
	}

	s.config = conf
}

// Priority method of the LabelSource interface
func (s *customSource) Priority() int { return 10 }

// GetLabels method of the LabelSource interface
func (s *customSource) GetLabels() (source.FeatureLabels, error) {
	// Get raw features from all sources
	domainFeatures := make(map[string]source.Features)
	for n, s := range source.GetAllFeatureSources() {
		domainFeatures[n] = s.GetFeatures()
	}

	labels := source.FeatureLabels{}
	allFeatureConfig := append(getStaticFeatureConfig(), *s.config...)
	allFeatureConfig = append(allFeatureConfig, getDirectoryFeatureConfig()...)
	utils.KlogDump(2, "custom features configuration:", "  ", allFeatureConfig)
	// Iterate over features
	for _, spec := range allFeatureConfig {
		features, err := spec.discover(domainFeatures)
		if err != nil {
			klog.Errorf("failed to discover feature: %q: %s", spec.Name, err.Error())
			continue
		}
		for k, v := range features {
			labels[k] = v
		}
	}
	return labels, nil
}

// Process a single feature spec by Matching on the defined rules.
// A feature is present if all defined Rules in a MatchRule return a match.
func (s *FeatureSpec) discover(features map[string]source.Features) (map[string]string, error) {
	ret := make(map[string]string)

	for _, matchRules := range s.MatchOn {
		if match, err := matchRules.Legacy.match(); err != nil {
			return nil, err
		} else if match {
			matched, err := matchRules.Domains.match(features)
			if err != nil {
				return nil, err
			} else if matched == nil {
				continue
			}

			// We have a match
			expandedName := s.Name
			if s.nameTemplate != nil {
				// Execute template to produce an array of labels
				var tmp bytes.Buffer
				if err := s.nameTemplate.Execute(&tmp, matched); err != nil {
					return nil, err
				}
				expandedName = tmp.String()
			}

			// Split out individual labels
			for _, item := range strings.Split(expandedName, "\n") {
				// Remove leading/trailing whitespace and skip empty lines
				if trimmed := strings.TrimSpace(item); trimmed != "" {
					n, v := s.getNameValue(trimmed)
					ret[n] = v
				}
			}

			if s.nameTemplate == nil {
				// No templating so we stop here (further matches would just
				// produce the same labels)
				break
			}
		}
	}
	return ret, nil
}

func (s *FeatureSpec) getNameValue(name string) (string, string) {
	// Value can be overridden in Name with "key=value". This is useful for
	// templates.
	if strings.ContainsRune(name, '=') {
		split := strings.SplitN(name, "=", 2)
		return split[0], split[1]
	}

	if s.Value != nil {
		return name, *s.Value
	}

	return name, "true"
}

type domainMatchedFeatures map[string]interface{}

func (r *Domains) match(features map[string]source.Features) (map[string]domainMatchedFeatures, error) {
	ret := make(map[string]domainMatchedFeatures, len(*r))

	for domain, rules := range *r {
		domainFeatures, ok := features[domain]
		if !ok {
			return nil, fmt.Errorf("unknown feature source/domain %q", domain)
		}
		for featureName, featureRules := range rules {
			var m bool
			var e error

			// Ignore case
			featureName = strings.ToLower(featureName)

			// Matched features
			matched := make(map[string]interface{})

			if f, ok := domainFeatures.Keys[featureName]; ok {
				v, err := featureRules.MatchGetKeys(f)
				m = len(v) > 0
				e = err
				matched[featureName] = v
			} else if f, ok := domainFeatures.Values[featureName]; ok {
				v, err := featureRules.MatchGetValues(f)
				m = len(v) > 0
				e = err
				matched[featureName] = v
			} else if f, ok := domainFeatures.Instances[featureName]; ok {
				v, err := featureRules.MatchGetInstances(f)
				m = len(v) > 0
				e = err
				matched[featureName] = v
			} else {
				return nil, fmt.Errorf("%q feature of source/domain %q not available", featureName, domain)
			}
			if e != nil {
				return nil, e
			} else if !m {
				return nil, nil
			}

			ret[domain] = matched
		}
	}
	return ret, nil
}

func (r *Legacy) match() (bool, error) {
	allRules := []source.CustomRule{
		r.PciID,
		r.UsbID,
		r.LoadedKMod,
		r.CpuID,
		r.Kconfig,
		r.Nodename,
	}

	// return true, nil if all rules match
	matchRules := func(rules []source.CustomRule) (bool, error) {
		for _, rule := range rules {
			if reflect.ValueOf(rule).IsNil() {
				continue
			}
			if match, err := rule.Match(); err != nil {
				return false, err
			} else if !match {
				return false, nil
			}
		}
		return true, nil
	}

	return matchRules(allRules)
}

func (m *MatchRule) UnmarshalJSON(data []byte) error {
	rule := &MatchRule{Domains: make(Domains)}

	// First, unmarshal legacy rules
	if err := json.Unmarshal(data, &rule.Legacy); err != nil {
		return err
	}

	// Next, unmarshal per-domain rules.
	// Start with unmarshalling into a map without trying to decode values
	raw := map[string]json.RawMessage{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	// Continue with decoding one domain at a time, skipping names that are
	// "registered" as legacy rules
	for k, v := range raw {
		k = strings.ToLower(k)
		if _, ok := legacyRuleNames[k]; ok {
			continue
		}

		r := make(nfdv1alpha1.DomainRule)
		if err := json.Unmarshal(v, &r); err != nil {
			return err
		}
		rule.Domains[k] = r
	}

	*m = *rule
	return nil
}

var legacyRuleNames map[string]struct{}

func init() {
	source.Register(&src)

	// Get fields names of Legacy
	v := reflect.ValueOf(Legacy{})
	legacyRuleNames = make(map[string]struct{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		legacyRuleNames[strings.ToLower(name)] = struct{}{}
	}
}
