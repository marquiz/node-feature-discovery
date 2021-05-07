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
	"encoding/json"
	"reflect"
	"strings"
)

// Legacy rules
type Legacy struct {
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

func (r *NodenameRule) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &r.MatchExpression); err != nil {
		return err
	}
	// Force regexp matching
	if r.Op == MatchIn {
		r.Op = MatchInRegexp
	}
	return nil
}

var legacyRuleNames map[string]struct{}

func init() {
	// Get fields names of Legacy
	v := reflect.ValueOf(Legacy{})
	legacyRuleNames = make(map[string]struct{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		legacyRuleNames[strings.ToLower(name)] = struct{}{}
	}
}
