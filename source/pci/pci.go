/*
Copyright 2018-2021 The Kubernetes Authors.

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

package pci

import (
	"fmt"
	"strings"

	"k8s.io/klog/v2"

	"sigs.k8s.io/node-feature-discovery/pkg/utils"
	"sigs.k8s.io/node-feature-discovery/source"
)

type Config struct {
	DeviceClassWhitelist []string `json:"deviceClassWhitelist,omitempty"`
	DeviceLabelFields    []string `json:"deviceLabelFields,omitempty"`
}

// newDefaultConfig returns a new config with pre-populated defaults
func newDefaultConfig() *Config {
	return &Config{
		DeviceClassWhitelist: []string{"03", "0b40", "12"},
		DeviceLabelFields:    []string{"class", "vendor"},
	}
}

// features contains all discovered features
type features struct {
	Devices map[deviceClass][]deviceInfo
}

type deviceInfo map[string]string

type deviceClass string

// pciSource implements the FeatureSource, LabelSource and ConfigurableSource interfaces
type pciSource struct {
	config   *Config
	features *features
}

// Singleton source instance
var (
	src pciSource
	_   source.FeatureSource      = &src
	_   source.LabelSource        = &src
	_   source.ConfigurableSource = &src
)

// Return name of the feature source
func (s *pciSource) Name() string { return "pci" }

// NewConfig method of the LabelSource interface
func (s *pciSource) NewConfig() source.Config { return newDefaultConfig() }

// GetConfig method of the LabelSource interface
func (s *pciSource) GetConfig() source.Config { return s.config }

// SetConfig method of the LabelSource interface
func (s *pciSource) SetConfig(conf source.Config) error {
	switch v := conf.(type) {
	case *Config:
		s.config = v
	default:
		return fmt.Errorf("invalid config type: %T", conf)
	}
	return nil
}

// Priority method of the LabelSource interface
func (s *pciSource) Priority() int { return 0 }

// GetLabels method of the LabelSource interface
func (s *pciSource) GetLabels() (source.FeatureLabels, error) {
	labels := source.FeatureLabels{}

	// Construct a device label format, a sorted list of valid attributes
	deviceLabelFields := make([]string, 0)
	configLabelFields := make(map[string]struct{}, len(s.config.DeviceLabelFields))
	for _, field := range s.config.DeviceLabelFields {
		configLabelFields[field] = struct{}{}
	}

	for _, attr := range mandatoryDevAttrs {
		if _, ok := configLabelFields[attr]; ok {
			deviceLabelFields = append(deviceLabelFields, attr)
			delete(configLabelFields, attr)
		}
	}
	if len(configLabelFields) > 0 {
		keys := []string{}
		for key := range configLabelFields {
			keys = append(keys, key)
		}
		klog.Warningf("invalid fields (%s) in deviceLabelFields, ignoring...", strings.Join(keys, ", "))
	}
	if len(deviceLabelFields) == 0 {
		klog.Warningf("no valid fields in deviceLabelFields defined, using the defaults")
		deviceLabelFields = []string{"class", "vendor"}
	}

	// Iterate over all device classes
	for class, classDevs := range s.features.Devices {
		for _, white := range s.config.DeviceClassWhitelist {
			if strings.HasPrefix(string(class), strings.ToLower(white)) {
				for _, dev := range classDevs {
					devLabel := ""
					for i, attr := range deviceLabelFields {
						devLabel += dev[attr]
						if i < len(deviceLabelFields)-1 {
							devLabel += "_"
						}
					}
					labels[devLabel+".present"] = true

					if _, ok := dev["sriov_totalvfs"]; ok {
						labels[devLabel+".sriov.capable"] = true
					}
				}
			}
		}
	}
	return labels, nil
}

// Discover method of the FeatureSource interface
func (s *pciSource) Discover() error {
	s.features = &features{
		Devices: make(map[deviceClass][]deviceInfo),
	}

	devs, err := detectPci()
	if err != nil {
		return fmt.Errorf("failed to detect PCI devices: %s", err.Error())
	}
	s.features.Devices = devs

	if klog.V(3).Enabled() {
		klog.Info("discovered pci features:\n", utils.Dump(s.features))
	}
	return nil
}

func init() {
	source.Register(&src)
}
