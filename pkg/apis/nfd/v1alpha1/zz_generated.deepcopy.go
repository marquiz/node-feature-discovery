// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Expressions) DeepCopyInto(out *Expressions) {
	{
		in := &in
		*out = make(Expressions, len(*in))
		for key, val := range *in {
			var outVal *MatchExpression
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(MatchExpression)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Expressions.
func (in Expressions) DeepCopy() Expressions {
	if in == nil {
		return nil
	}
	out := new(Expressions)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FeatureRule) DeepCopyInto(out *FeatureRule) {
	*out = *in
	in.MatchExpressionSet.DeepCopyInto(&out.MatchExpressionSet)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FeatureRule.
func (in *FeatureRule) DeepCopy() *FeatureRule {
	if in == nil {
		return nil
	}
	out := new(FeatureRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LabelRule) DeepCopyInto(out *LabelRule) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LabelRule.
func (in *LabelRule) DeepCopy() *LabelRule {
	if in == nil {
		return nil
	}
	out := new(LabelRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LabelRule) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LabelRuleList) DeepCopyInto(out *LabelRuleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LabelRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LabelRuleList.
func (in *LabelRuleList) DeepCopy() *LabelRuleList {
	if in == nil {
		return nil
	}
	out := new(LabelRuleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LabelRuleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LabelRuleSpec) DeepCopyInto(out *LabelRuleSpec) {
	*out = *in
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]Rule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LabelRuleSpec.
func (in *LabelRuleSpec) DeepCopy() *LabelRuleSpec {
	if in == nil {
		return nil
	}
	out := new(LabelRuleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchExpression) DeepCopyInto(out *MatchExpression) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = make(MatchValue, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchExpression.
func (in *MatchExpression) DeepCopy() *MatchExpression {
	if in == nil {
		return nil
	}
	out := new(MatchExpression)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchExpressionSet) DeepCopyInto(out *MatchExpressionSet) {
	*out = *in
	if in.Expressions != nil {
		in, out := &in.Expressions, &out.Expressions
		*out = make(Expressions, len(*in))
		for key, val := range *in {
			var outVal *MatchExpression
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(MatchExpression)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchExpressionSet.
func (in *MatchExpressionSet) DeepCopy() *MatchExpressionSet {
	if in == nil {
		return nil
	}
	out := new(MatchExpressionSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in MatchRule) DeepCopyInto(out *MatchRule) {
	{
		in := &in
		*out = make(MatchRule, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchRule.
func (in MatchRule) DeepCopy() MatchRule {
	if in == nil {
		return nil
	}
	out := new(MatchRule)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in MatchValue) DeepCopyInto(out *MatchValue) {
	{
		in := &in
		*out = make(MatchValue, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchValue.
func (in MatchValue) DeepCopy() MatchValue {
	if in == nil {
		return nil
	}
	out := new(MatchValue)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in MatchedInstance) DeepCopyInto(out *MatchedInstance) {
	{
		in := &in
		*out = make(MatchedInstance, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchedInstance.
func (in MatchedInstance) DeepCopy() MatchedInstance {
	if in == nil {
		return nil
	}
	out := new(MatchedInstance)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchedKey) DeepCopyInto(out *MatchedKey) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchedKey.
func (in *MatchedKey) DeepCopy() *MatchedKey {
	if in == nil {
		return nil
	}
	out := new(MatchedKey)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchedValue) DeepCopyInto(out *MatchedValue) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchedValue.
func (in *MatchedValue) DeepCopy() *MatchedValue {
	if in == nil {
		return nil
	}
	out := new(MatchedValue)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rule) DeepCopyInto(out *Rule) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(string)
		**out = **in
	}
	if in.MatchAny != nil {
		in, out := &in.MatchAny, &out.MatchAny
		*out = make([]MatchRule, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(MatchRule, len(*in))
				for key, val := range *in {
					(*out)[key] = *val.DeepCopy()
				}
			}
		}
	}
	if in.MatchAll != nil {
		in, out := &in.MatchAll, &out.MatchAll
		*out = make([]MatchRule, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(MatchRule, len(*in))
				for key, val := range *in {
					(*out)[key] = *val.DeepCopy()
				}
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rule.
func (in *Rule) DeepCopy() *Rule {
	if in == nil {
		return nil
	}
	out := new(Rule)
	in.DeepCopyInto(out)
	return out
}
