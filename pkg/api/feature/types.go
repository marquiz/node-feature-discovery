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

package feature

// DomainFeatures is the collection of all discovered features of one domain.
type DomainFeatures struct {
	Keys      map[string]KeyFeatures      `protobuf:"bytes,1,rep,name=keys"`
	Values    map[string]ValueFeatures    `protobuf:"bytes,2,rep,name=values"`
	Instances map[string]InstanceFeatures `protobuf:"bytes,3,rep,name=instances"`
}

// KeyFeatures is a set of simple features only containing names without values.
type KeyFeatures struct {
	Features map[string]Nil `protobuf:"bytes,1,rep,name=features"`
}

// ValueFeatures is a set of features having string value.
type ValueFeatures struct {
	Features map[string]string `protobuf:"bytes,1,rep,name=features"`
}

// InstanceFeatures is a set of features each of which is an instance having multiple attributes.
type InstanceFeatures struct {
	Features []InstanceFeature `protobuf:"bytes,1,rep,name=features"`
}

// InstanceFeature represents one instance of a complex features, e.g. a device.
type InstanceFeature struct {
	Attributes map[string]string `protobuf:"bytes,1,rep,name=attributes"`
}

// Nil is a  summy empty struct for protobuf compatibility
type Nil struct{}

// NewDomainFeatures creates a new instance of Features, initializing specified
// features to empty values
func NewDomainFeatures() *DomainFeatures {
	return &DomainFeatures{
		Keys:      make(map[string]KeyFeatures),
		Values:    make(map[string]ValueFeatures),
		Instances: make(map[string]InstanceFeatures)}
}

func NewKeyFeatures() *KeyFeatures { return &KeyFeatures{Features: make(map[string]Nil)} }

func NewValueFeatures() *ValueFeatures { return &ValueFeatures{Features: make(map[string]string)} }

func NewInstanceFeatures() *InstanceFeatures { return &InstanceFeatures{} }

func NewInstanceFeature() *InstanceFeature {
	return &InstanceFeature{Attributes: make(map[string]string)}
}
