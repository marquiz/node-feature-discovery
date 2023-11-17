/*
Copyright 2024 The Kubernetes Authors.

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

package v1alpha2

// NewNodeFeatureSpec creates a new emprty instance of NodeFeatureSpec type,
// initializing all fields to proper empty values.
func NewNodeFeatureSpec() *NodeFeatureSpec {
	return &NodeFeatureSpec{
		Features: make([]Feature, 0),
		Labels:   make(map[string]string),
	}
}

// NewFeature creates a new instance of Feature.
func NewFeature(name string, elements ...FeatureElement) Feature {
	return Feature{Name: name, Elements: elements}
}

// NewFeatureElement creates a new FeatureElement instance.
func NewFeatureElement(name string, attrs map[string]string) *FeatureElement {
	if attrs == nil {
		attrs = make(map[string]string)
	}
	return &FeatureElement{Name: name, Attributes: attrs}
}

// Exists returns a non-empty string if a feature exists.
func (f *NodeFeatureSpec) FeatureExists(name string) bool {
	for i := range f.Features {
		if f.Features[i].Name == name {
			return true
		}
	}
	return false
}

// InsertFeature inserts one feature into NodeFeatureSpec.
func (f *NodeFeatureSpec) InsertFeature(featureName, elementName string, attrs map[string]string) {
	for i := range f.Features {
		if f.Features[i].Name == featureName {
			mergeElement(NewFeatureElement(elementName, attrs), &f.Features[i].Elements)
			return
		}
	}
	f.Features = append(f.Features, NewFeature(featureName, *NewFeatureElement(elementName, attrs)))
}

// InsertFeatures inserts multiple feature elements into NodeFeatureSpec.
func (f *NodeFeatureSpec) InsertFeatures(featureName string, elems ...FeatureElement) {
	for i := range f.Features {
		if f.Features[i].Name == featureName {
			for j := range elems {
				mergeElement(&elems[j], &f.Features[i].Elements)
			}
			return
		}
	}
	f.Features = append(f.Features, NewFeature(featureName, elems...))
}

// MergeInto merges two NodeFeatureSpecs into one. Data in the input object takes
// precedence (overwrite) over data of the existing object we're merging into.
func (in *NodeFeatureSpec) MergeInto(out *NodeFeatureSpec) {
	for i := range in.Features {
		for j := range out.Features {
			if in.Features[i].Name == out.Features[j].Name {
				in.Features[i].MergeInto(&out.Features[j])
				break
			}
		}
	}

	if in.Labels != nil {
		if out.Labels == nil {
			out.Labels = make(map[string]string, len(in.Labels))
		}
		for key, val := range in.Labels {
			out.Labels[key] = val
		}
	}
}

// MergeInto merges two sets of instance featues.
func (in *Feature) MergeInto(out *Feature) {
	if in.Elements != nil {
		if out.Elements == nil {
			out.Elements = make([]FeatureElement, 0, len(in.Elements))
		}
		for _, e := range in.Elements {
			mergeElement(&e, &out.Elements)
		}
	}
}

func mergeElement(in *FeatureElement, out *[]FeatureElement) {
	for i := range *out {
		if (*out)[i].Name == in.Name {
			(*out)[i] = *in.DeepCopy()
			return
		}
	}
	*out = append(*out, *in.DeepCopy())
}
