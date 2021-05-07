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

package v1alpha1

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LabelRuleList contains a list of LabelRule objects.
// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LabelRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []LabelRule `json:"items"`
}

// LabelRule resource specifies a configuration for custom node labeling.
// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
type LabelRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec LabelRuleSpec `json:"spec"`
}

// RuleSpec defines a set of labeling rules.
type LabelRuleSpec struct {
	Rules []Rule `json:"rules"`
}

// Rule defines a rule for creating on feature label.
type Rule struct {
	// Name of the label to be generated.
	Name string `json:"name"`

	// Value of the label, optional.
	// +optional
	Value *string `json:"value,omitempty"`

	// MatchOn specifies a list of alternative expression sets
	MatchOn []MatchRule `json:"matchOn"`

	// nameExpander is a private helper/cache for handling golang templates
	name *nameExpander `json:"-"`
}

// MatchRule defines one complete set of rules to satisfy a successful match.
type MatchRule struct {
	Domains map[string]DomainRule `json:",inline"`

	// Legacy
	Legacy Legacy `json:"-"`
}

// DomainRule defines per-feature rules for one domain.
type DomainRule map[string]FeatureRule

// FeatureRule defines rules for one feature, matching against its attributes.
type FeatureRule struct {
	MatchExpressionSet `json:",inline"`
}

type nameExpander struct {
	nameTemplate *template.Template
}

func newNameExpander(name string) (*nameExpander, error) {
	e := nameExpander{}

	if strings.Contains(name, "{{") {

		tmpl, err := template.New("").Option("missingkey=error").Parse(name)
		if err != nil {
			return nil, fmt.Errorf("invalid template in rule name: %w", err)
		}
		e.nameTemplate = tmpl
	}
	return &e, nil
}

func (e *nameExpander) expand(data interface{}) (string, error) {
	if e.nameTemplate == nil {
		return "", fmt.Errorf("not a template")
	}
	var tmp bytes.Buffer
	if err := e.nameTemplate.Execute(&tmp, data); err != nil {
		return "", err
	}
	return tmp.String(), nil
}

func (r *Rule) ExpandName(data interface{}) (string, error) {
	if isTemplate, err := r.NameIsTemplate(); err != nil {
		return "", err
	} else if isTemplate {
		return r.name.expand(data)
	}
	return r.Name, nil
}

func (r *Rule) NameIsTemplate() (bool, error) {
	if r.name == nil {
		n, err := newNameExpander(r.Name)
		if err != nil {
			return true, err
		}
		r.name = n
	}

	return r.name.nameTemplate != nil, nil
}
