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

package pci

import (
	"encoding/json"

	"sigs.k8s.io/node-feature-discovery/source"
)

type PciIDRule struct {
	source.MatchExpressionSet
}

// Match PCI devices on provided PCI device attributes
func (r *PciIDRule) Match() (bool, error) {
	for _, classDevs := range src.features.Devices {
		for _, dev := range classDevs {
			// match rule on a single device
			if match, err := r.MatchValues(dev); err != nil {
				return false, err
			} else if match {
				return true, nil
			}
		}
	}
	return false, nil
}

func NewPciIDRule(ruleConfig []byte) (source.CustomRule, error) {
	r := new(PciIDRule)
	return r, json.Unmarshal(ruleConfig, r)
}

func init() {
	source.RegisterCustomRule("pciId", NewPciIDRule)
}
