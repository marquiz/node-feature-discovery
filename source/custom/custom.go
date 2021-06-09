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
	"fmt"
	"strings"

	"k8s.io/klog/v2"

	"sigs.k8s.io/node-feature-discovery/pkg/api/feature"
	nfdv1alpha1 "sigs.k8s.io/node-feature-discovery/pkg/apis/nfd/v1alpha1"
	"sigs.k8s.io/node-feature-discovery/pkg/utils"
	"sigs.k8s.io/node-feature-discovery/source"
	"sigs.k8s.io/node-feature-discovery/source/cpu"
	"sigs.k8s.io/node-feature-discovery/source/kernel"
	"sigs.k8s.io/node-feature-discovery/source/pci"
	"sigs.k8s.io/node-feature-discovery/source/system"
	"sigs.k8s.io/node-feature-discovery/source/usb"
)

const Name = "custom"

type config []nfdv1alpha1.Rule

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
	s.config = conf
}

// Priority method of the LabelSource interface
func (s *customSource) Priority() int { return 10 }

// GetLabels method of the LabelSource interface
func (s *customSource) GetLabels() (source.FeatureLabels, error) {
	// Get raw features from all sources
	domainFeatures := make(map[string]*feature.DomainFeatures)
	for n, s := range source.GetAllFeatureSources() {
		domainFeatures[n] = s.GetFeatures()
	}

	labels := source.FeatureLabels{}
	allFeatureConfig := append(getStaticFeatureConfig(), *s.config...)
	allFeatureConfig = append(allFeatureConfig, getDirectoryFeatureConfig()...)
	utils.KlogDump(2, "custom features configuration:", "  ", allFeatureConfig)
	// Iterate over features
	for _, spec := range allFeatureConfig {
		features, err := MatchRule(&spec, domainFeatures)
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

// MatchRule processes a single feature rule.
func MatchRule(rule *nfdv1alpha1.Rule, features map[string]*feature.DomainFeatures) (map[string]string, error) {
	ret := make(map[string]string)

	for _, matchRules := range rule.MatchOn {
		if matchedLegacy, err := matchLegacy(matchRules.Legacy, features); err != nil {
			return nil, err
		} else if matchedLegacy {
			matched, err := matchDomains(matchRules.Domains, features)
			if err != nil {
				return nil, err
			} else if matched == nil {
				continue
			}

			// We have a match
			name, err := rule.ExpandName(matched)
			if err != nil {
				return nil, err
			}

			// Split out individual labels
			for _, item := range strings.Split(name, "\n") {
				// Remove leading/trailing whitespace and skip empty lines
				if trimmed := strings.TrimSpace(item); trimmed != "" {
					n, v := getNameValue(rule, trimmed)
					ret[n] = v
				}
			}

			if isTemplate, err := rule.NameIsTemplate(); err != nil {
				return nil, err
			} else if !isTemplate {
				// No templating so we stop here (further matches would just
				// produce the same labels)
				break
			}
		}
	}
	return ret, nil
}

func getNameValue(rule *nfdv1alpha1.Rule, name string) (string, string) {
	// Value can be overridden in Name with "key=value". This is useful for
	// templates.
	if strings.ContainsRune(name, '=') {
		split := strings.SplitN(name, "=", 2)
		return split[0], split[1]
	}

	if rule.Value != nil {
		return name, *rule.Value
	}

	return name, "true"
}

type domainMatchedFeatures map[string]interface{}

func matchDomains(domains map[string]nfdv1alpha1.DomainRule, features map[string]*feature.DomainFeatures) (map[string]domainMatchedFeatures, error) {
	ret := make(map[string]domainMatchedFeatures, len(domains))

	for domain, rules := range domains {
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
				v, err := featureRules.MatchGetKeys(f.Features)
				m = len(v) > 0
				e = err
				matched[featureName] = v
			} else if f, ok := domainFeatures.Values[featureName]; ok {
				v, err := featureRules.MatchGetValues(f.Features)
				m = len(v) > 0
				e = err
				matched[featureName] = v
			} else if f, ok := domainFeatures.Instances[featureName]; ok {
				v, err := featureRules.MatchGetInstances(f.Features)
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

func matchLegacy(rules nfdv1alpha1.Legacy, features map[string]*feature.DomainFeatures) (bool, error) {
	if rules.CpuID != nil {
		if f, ok := features[cpu.Name].Keys[cpu.CpuidFeature]; !ok {
			return false, fmt.Errorf("cpuid information not available")
		} else if match, err := rules.CpuID.MatchKeys(f.Features); !match || err != nil {
			return match, err
		}
	}

	if rules.LoadedKMod != nil {
		if f, ok := features[kernel.Name].Keys[kernel.LoadedModuleFeature]; !ok {
			return false, fmt.Errorf("information about loaded modules not available")
		} else if match, err := rules.LoadedKMod.MatchKeys(f.Features); !match || err != nil {
			return match, err
		}
	}

	if rules.Kconfig != nil {
		if f, ok := features[kernel.Name].Values[kernel.ConfigFeature]; !ok {
			return false, fmt.Errorf("kernel config options not available")
		} else if match, err := rules.Kconfig.MatchValues(f.Features); !match || err != nil {
			return match, err
		}
	}

	if rules.PciID != nil {
		if f, ok := features[pci.Name].Instances[pci.DeviceFeature]; !ok {
			return false, fmt.Errorf("pci device information not available")
		} else if match, err := rules.PciID.MatchInstances(f.Features); !match || err != nil {
			return match, err
		}
	}

	if rules.UsbID != nil {
		if f, ok := features[usb.Name].Instances[usb.DeviceFeature]; !ok {
			return false, fmt.Errorf("usb device information not available")
		} else if match, err := rules.UsbID.MatchInstances(f.Features); !match || err != nil {
			return match, err
		}
	}

	if rules.Nodename != nil {
		if f, ok := features[system.Name].Values[system.NameFeature]; !ok {
			return false, fmt.Errorf("system name information not available")
		} else {
			n, ok := f.Features["nodename"]
			if match, err := rules.Nodename.MatchExpression.Match(ok, n); !match || err != nil {
				return match, err
			}
		}
	}
	return true, nil
}

func init() {
	source.Register(&src)
}
