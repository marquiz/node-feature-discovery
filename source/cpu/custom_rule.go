/*
Copyright 2021 The Kubernetes Authors.

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

package cpu

import (
	"k8s.io/klog/v2"

	"sigs.k8s.io/node-feature-discovery/source"
)

// CustomRule implements the source.CustomRule interface
type CustomRule struct {
	Cpuid  source.MatchExpressionSet
	Cstate source.MatchExpressionSet
	Pstate source.MatchExpressionSet
	Rdt    source.MatchExpressionSet
	Sst    struct {
		BaseFrequency struct {
			Enabled source.MatchExpression
		}
	}
	Topology struct {
		HardwareMultithreading source.MatchExpression
	}
}

func (r *CustomRule) Match() (bool, error) {
	if m, err := r.Cpuid.MatchKeys(src.features.Cpuid); err != nil || m == false {
		klog.V(2).Infof("cpu CustomRule: failed to match cpuid")
		return m, err
	}
	if m, err := r.Cstate.MatchValues(src.features.Cstate); err != nil || m == false {
		klog.V(2).Infof("cpu CustomRule: failed to match cstate")
		return m, err
	}
	if m, err := r.Pstate.MatchValues(src.features.Pstate); err != nil || m == false {
		klog.V(2).Infof("cpu CustomRule: failed to match pstate")
		return m, err
	}
	if m, err := r.Rdt.MatchKeys(src.features.Rdt); err != nil || m == false {
		klog.V(2).Infof("cpu CustomRule: failed to match rdt")
		return m, err
	}

	ok, v := src.features.Sst.BaseFrequency.Enabled.Get()
	if m, err := r.Sst.BaseFrequency.Enabled.Match(ok, v); err != nil || m == false {
		klog.V(2).Infof("cpu CustomRule: failed to match sst")
		return m, err
	}
	return true, nil
}
