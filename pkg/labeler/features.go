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
package labeler

import (
	"sigs.k8s.io/node-feature-discovery/source"
)

// ToFeatures converts DomainFeatures message into a NFD internal format
func (x *DomainFeatures) ToFeatures() *source.Features {
	f := source.NewFeatures()

	for k, v := range x.Keys {
		features := make(source.KeyAttributes, len(v.Attributes))
		for feature := range v.Attributes {
			features[feature] = struct{}{}
		}
		f.Keys[k] = features
	}

	for k, v := range x.Values {
		f.Values[k] = v.Attributes
	}

	for k, v := range x.Instances {
		features := make(source.InstanceAttributes, len(v.Attributes))
		for i, instance := range v.Attributes {
			features[i] = instance.Info
		}
		f.Instances[k] = features
	}

	return f
}
