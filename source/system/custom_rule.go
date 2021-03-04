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

package system

import (
	"k8s.io/klog/v2"

	"sigs.k8s.io/node-feature-discovery/source"
)

// CustomRule implements the source.CustomRule interface
type CustomRule struct {
	NodeName  source.MatchExpression
	OsRelease source.MatchExpressionSet
}

func (r *CustomRule) Match() (bool, error) {
	if m, err := r.OsRelease.MatchValues(src.features.OsRelease); err != nil || m == false {
		klog.V(2).Infof("system CustomRule: failed to match osRelease")
		return m, err
	}

	n := src.features.NodeName
	if m, err := r.NodeName.Match(n != "", n); err != nil || m == false {
		klog.V(2).Infof("system CustomRule: failed to match nodeName")
		return m, err
	}
	return true, nil
}
