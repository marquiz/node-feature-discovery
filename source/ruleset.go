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
package source

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CustomRuleSet is a collection of rules which to match
type CustomRuleSet map[string]CustomRule

func (r *CustomRuleSet) UnmarshalJSON(data []byte) error {
	*r = make(CustomRuleSet)

	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		name := strings.ToLower(k)
		if f, ok := rules[name]; ok {
			rule, err := f(v)
			if err != nil {
				return err
			}
			(*r)[name] = rule
		} else {
			supported := make([]string, 0, len(rules))
			for k := range rules {
				supported = append(supported, k)
			}
			return fmt.Errorf("unknown rule type %q, must be one of %s", k, strings.Join(supported, ", "))
		}
	}

	return nil
}
