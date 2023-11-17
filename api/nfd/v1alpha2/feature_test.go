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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeature(t *testing.T) {
	f1 := Feature{}
	f2 := Feature{}
	var expectedElems []FeatureElement = nil

	f2.MergeInto(&f1)
	assert.Equal(t, expectedElems, f1.Elements)

	f2 = NewFeature("feat")
	expectedElems = []FeatureElement{}
	f2.MergeInto(&f1)
	assert.Equal(t, expectedElems, f1.Elements)

	f2 = NewFeature("feat", FeatureElement{})
	expectedElems = append(expectedElems, FeatureElement{})
	f2.MergeInto(&f1)
	assert.Equal(t, expectedElems, f1.Elements)

	f2 = NewFeature("feat",
		FeatureElement{
			Name: "e1",
			Attributes: map[string]string{
				"a1": "v1",
				"a2": "v2",
			},
		})
	expectedElems = append(expectedElems, *NewFeatureElement("e1", map[string]string{"a1": "v1", "a2": "v2"}))
	f2.MergeInto(&f1)
	assert.Equal(t, expectedElems, f1.Elements)

	f2.Elements[0].Attributes["a2"] = "v2.2"
	expectedElems = append(expectedElems, *NewFeatureElement("e1", map[string]string{"a1": "v1", "a2": "v2.2"}))
	f2.MergeInto(&f1)
	assert.Equal(t, expectedElems, f1.Elements)
}

func TestFeatureSpec(t *testing.T) {
	// Test FeatureExists() and InserFeature()
	f := NodeFeatureSpec{}
	assert.Empty(t, f.FeatureExists("dom.name"), "empty features shouldn't contain anything")

	f.InsertFeature("dom.inst", "i1", map[string]string{"k1": "v1", "k2": "v2"})
	expectedAttributes := map[string]string{"k1": "v1", "k2": "v2"}
	assert.True(t, f.FeatureExists("dom.inst"), "feature should exist")
	assert.Equal(t, expectedAttributes, f.Features[0].Elements[0].Attributes)

	f.InsertFeature("dom.inst", "i1", map[string]string{"k2": "v2.override", "k3": "v3"})
	expectedAttributes = map[string]string{"k2": "v2.override", "k3": "v3"}
	assert.Equal(t, expectedAttributes, f.Features[0].Elements[0])

	// Test merging
	f2 := NodeFeatureSpec{}
	expectedFeatures := NodeFeatureSpec{}

	f2.MergeInto(&f)
	assert.Equal(t, expectedFeatures, f)

	f2.Labels = map[string]string{"l1": "v1", "l2": "v2"}
	f2.InsertFeatures("dom.flags",
		FeatureElement{Name: "k1"},
		FeatureElement{Name: "k2"})
	f2.InsertFeatures("dom.attr",
		FeatureElement{Name: "k1", Attributes: map[string]string{"val": "v1"}},
		FeatureElement{Name: "k2", Attributes: map[string]string{"val": "v1"}})
	f2.InsertFeatures("dom.inst",
		FeatureElement{Name: "i1", Attributes: map[string]string{"a1": "v1.1", "a2": "v1.2"}},
		FeatureElement{Name: "i2", Attributes: map[string]string{"a1": "v2.1", "a2": "v2.2"}})

	f2.MergeInto(&f)
	assert.Equal(t, f2, f)

	f2.InsertFeature("dom.flag", "k3", nil)
	f2.InsertFeature("dom.attr", "k1", map[string]string{"val": "v1.override"})
	f2.InsertFeature("dom.inst", "i2", map[string]string{"a1": "v2.2"})
	f2.InsertFeature("dom.inst", "i3", map[string]string{"a1": "v3.1", "a3": "v3.3"})

	f2.MergeInto(&f)

	expectedFeatures = NodeFeatureSpec{
		Features: []Feature{
			NewFeature("dom.flag", FeatureElement{Name: "k1"}, FeatureElement{Name: "k2"}),
			NewFeature("dom.attr",
				FeatureElement{Name: "k1", Attributes: map[string]string{"val": "v1.override"}},
				FeatureElement{Name: "k2", Attributes: map[string]string{"val": "v2"}}),
			NewFeature("dom.inst",
				FeatureElement{Name: "i1", Attributes: map[string]string{"a1": "v1.1", "a2": "v1.2"}},
				FeatureElement{Name: "i2", Attributes: map[string]string{"a1": "v2.2"}},
				FeatureElement{Name: "i3", Attributes: map[string]string{"a1": "v3.1", "a3": "v3.3"}}),
		},
		Labels: map[string]string{"l1": "v1", "l2": "v2"},
	}
	assert.Equal(t, expectedFeatures, f)

	// Check that second merge updates the object correctly
	f2 = *NewNodeFeatureSpec()
	f2.Labels = map[string]string{"l1": "v1.override", "l3": "v3"}
	f2.InsertFeature("dom.flag2", "k3", nil)

	expectedFeatures.Labels["l1"] = "v1.override"
	expectedFeatures.Labels["l3"] = "v3"
	expectedFeatures.Features = append(expectedFeatures.Features,
		Feature{
			Name:     "dom.flag2",
			Elements: []FeatureElement{FeatureElement{Name: "k3"}},
		})

	f2.MergeInto(&f)
	assert.Equal(t, expectedFeatures, f)
}
