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

	// MatchOn specifies a list of alternative legacy rule sets
	// +optional
	MatchOn []LegacyRule `json:"matchOn"`

	// MatchAny specifies a list of expression sets one of which must match
	// +optional
	MatchAny []MatchRule `json:"matchAny"`

	// MatchAll specifies a list of expression sets all of which must match
	// +optional
	MatchAll []MatchRule `json:"matchAll"`
}

// MatchRule defines one complete set of rules to satisfy a successful match.
type MatchRule map[string]DomainRule

// DomainRule defines per-feature rules for one domain.
type DomainRule map[string]FeatureRule

// FeatureRule defines rules for one feature, matching against its attributes.
type FeatureRule struct {
	MatchExpressionSet `json:",inline"`
}

// LegacyRule defines a set of legacy rules, supported for backwards compatibility.
type LegacyRule struct {
	PciID      *MatchExpressionSet `json:"pciId,omitempty"`
	UsbID      *MatchExpressionSet `json:"usbId,omitempty"`
	LoadedKMod *MatchExpressionSet `json:"loadedKMod,omitempty"`
	CpuID      *MatchExpressionSet `json:"cpuId,omitempty"`
	Kconfig    *MatchExpressionSet `json:"kConfig,omitempty"`
	Nodename   *NodenameRule       `json:"nodename,omitempty"`
}

type NodenameRule struct {
	MatchExpression `json:",inline"`
}

// MatchExpression defines a set of MatchExpressions
type MatchExpressionSet map[string]*MatchExpression

// MatchExpression defines the expression to use for matching.
type MatchExpression struct {
	// Op is the operand to be applied.
	Op MatchOp `json:"op"`

	// Value is the list of values that the operand is applied on.
	// +optional
	Value MatchValue `json:"value,omitempty"`
}

// MatchOp is the operator to user for matching.
// +kubebuilder:validation:Enum="In";"NotIn";"InRegexp";"Exists";"DoesNotExist";"Gt";"Lt";"IsTrue";"IsFalse"
type MatchOp string

// MatchValue defines an array of values to use for matching.
type MatchValue []string

const (
	MatchAny          MatchOp = ""
	MatchIn           MatchOp = "In"
	MatchNotIn        MatchOp = "NotIn"
	MatchInRegexp     MatchOp = "InRegexp"
	MatchExists       MatchOp = "Exists"
	MatchDoesNotExist MatchOp = "DoesNotExist"
	MatchGt           MatchOp = "Gt"
	MatchLt           MatchOp = "Lt"
	MatchIsTrue       MatchOp = "IsTrue"
	MatchIsFalse      MatchOp = "IsFalse"
)
