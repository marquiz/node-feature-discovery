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

package kernel

import (
	"k8s.io/klog/v2"

	"sigs.k8s.io/node-feature-discovery/source"
)

// CustomRule implements the source.CustomRule interface
type CustomRule struct {
	Config    source.MatchExpressionSet
	LoadedMod source.MatchExpressionSet
	Version   source.MatchExpressionSet
	Selinux   struct {
		Enabled source.MatchExpression
	}
}

func (r *CustomRule) Match() (bool, error) {
	if m, err := r.Config.MatchValues(src.features.Config); err != nil || m == false {
		klog.V(2).Infof("kernel CustomRule: failed to match config")
		return m, err
	}
	if m, err := r.LoadedMod.MatchKeys(src.features.LoadedModules); err != nil || m == false {
		klog.V(2).Infof("kernel CustomRule: failed to match loadedModules")
		return m, err
	}
	if m, err := r.Version.MatchValues(src.features.Version); err != nil || m == false {
		klog.V(2).Infof("kernel CustomRule: failed to match version")
		return m, err
	}

	ok, v := src.features.Selinux.Enabled.Get()
	if m, err := r.Selinux.Enabled.Match(ok, v); err != nil || m == false {
		klog.V(2).Infof("kernel CustomRule: failed to match selinux")
		return m, err
	}
	return true, nil
}
