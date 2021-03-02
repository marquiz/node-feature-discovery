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
	"encoding/json"

	"sigs.k8s.io/node-feature-discovery/source"
)

// Rule that matches on nodenames configured in a ConfigMap
type NodenameRule struct {
	source.MatchExpression
}

// Force implementation of Rule
var _ source.CustomRule = &NodenameRule{}

func (r *NodenameRule) Match() (bool, error) {
	return r.MatchExpression.Match(src.features.NodeName != "", src.features.NodeName)
}

func (r *NodenameRule) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &r.MatchExpression); err != nil {
		return err
	}
	// Force regexp matching
	if r.Op == source.MatchIn {
		r.Op = source.MatchInRegexp
	}
	return nil
}

func NewNodenameRule(ruleConfig []byte) (source.CustomRule, error) {
	r := new(NodenameRule)
	return r, json.Unmarshal(ruleConfig, r)
}

func init() {
	source.RegisterCustomRule("nodename", NewNodenameRule)
}
