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
	"k8s.io/klog/v2"

	"sigs.k8s.io/node-feature-discovery/pkg/api/feature"
	nfdv1alpha1 "sigs.k8s.io/node-feature-discovery/pkg/apis/nfd/v1alpha1"
	"sigs.k8s.io/node-feature-discovery/pkg/utils"
	"sigs.k8s.io/node-feature-discovery/source"
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
	for _, rule := range allFeatureConfig {
		features, err := rule.Match(domainFeatures)
		if err != nil {
			klog.Errorf("failed to discover feature: %q: %s", rule.Name, err.Error())
			continue
		}
		for k, v := range features {
			labels[k] = v
		}
	}
	return labels, nil
}

func init() {
	source.Register(&src)
}
