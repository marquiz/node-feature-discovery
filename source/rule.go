/*
Copyright 2020-2021 The Kubernetes Authors.

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
	"regexp"
	"strconv"
	"strings"

	"k8s.io/klog/v2"
)

type CustomRule interface {
	// Match on rule
	Match() (bool, error)
}

type MatchExpression struct {
	Op     MatchOp
	Values []string `json:",omitempty"`
}

type MatchExpressionSet map[string]*MatchExpression

type MatchOp string

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

var matchOps = map[MatchOp]struct{}{
	MatchAny:          struct{}{},
	MatchIn:           struct{}{},
	MatchNotIn:        struct{}{},
	MatchInRegexp:     struct{}{},
	MatchExists:       struct{}{},
	MatchDoesNotExist: struct{}{},
	MatchGt:           struct{}{},
	MatchLt:           struct{}{},
	MatchIsTrue:       struct{}{},
	MatchIsFalse:      struct{}{},
}

func (m *MatchOp) UnmarshalJSON(data []byte) error {
	var raw string

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	if _, ok := matchOps[MatchOp(raw)]; !ok {
		return fmt.Errorf("invalid Op %q", raw)
	}
	*m = MatchOp(raw)
	return nil
}

func NewMatchExpression(op MatchOp, values ...string) *MatchExpression {
	return &MatchExpression{
		Op:     op,
		Values: values,
	}
}

func (m *MatchExpression) Validate() error {
	if _, ok := matchOps[m.Op]; !ok {
		return fmt.Errorf("invalid Op %q", m.Op)
	}
	return nil
}

func (m *MatchExpression) Match(valid bool, value interface{}) (bool, error) {
	switch m.Op {
	case MatchAny:
		return true, nil
	case MatchExists:
		return valid, nil
	case MatchDoesNotExist:
		return !valid, nil
	}

	if valid {
		value := fmt.Sprintf("%v", value)
		switch m.Op {
		case MatchIn:
			for _, v := range m.Values {
				if value == v {
					return true, nil
				}
			}
		case MatchNotIn:
			for _, v := range m.Values {
				if value == v {
					return false, nil
				}
			}
			return true, nil
		case MatchInRegexp:
			for _, v := range m.Values {
				re, err := regexp.Compile(v)
				if err != nil {
					return false, fmt.Errorf("invalid regexp %q in %v", v, m)
				}
				if re.MatchString(value) {
					return true, nil
				}
			}
		case MatchGt, MatchLt:
			i, err := strconv.Atoi(value)
			if err != nil {
				return false, fmt.Errorf("not a number %q", value)
			}
			for _, v := range m.Values {
				j, err := strconv.Atoi(v)
				if err != nil {
					return false, fmt.Errorf("not a number %q in %v", v, m)
				}
				if (i < j && m.Op == MatchLt) || (i > j && m.Op == MatchGt) {
					return true, nil
				}
			}
		case MatchIsTrue:
			return value == "true", nil
		case MatchIsFalse:
			return value == "false", nil
		}
	}
	return false, nil
}

func (m *MatchExpression) MatchKeys(name string, keys map[string]struct{}) (bool, error) {
	klog.V(3).Infof("matching %q %q against %v", name, m.Op, keys)

	_, ok := keys[name]
	switch m.Op {
	case MatchAny:
		return true, nil
	case MatchExists:
		return ok, nil
	case MatchDoesNotExist:
		return !ok, nil
	default:
		return false, fmt.Errorf("invalid Op %q when matching keys", m.Op)
	}

	return false, nil
}

func (m *MatchExpression) MatchValues(name string, values map[string]string) (bool, error) {
	klog.V(3).Infof("matching %q %q %v against %v", name, m.Op, m.Values, values)
	v, ok := values[name]
	return m.Match(ok, v)
}

func (m *MatchExpression) UnmarshalJSON(data []byte) error {
	raw := new(interface{})

	err := json.Unmarshal(data, raw)
	if err != nil {
		return err
	}

	switch v := (*raw).(type) {
	case string:
		*m = *NewMatchExpression(MatchIn, v)
	case bool:
		*m = *NewMatchExpression(MatchIn, strconv.FormatBool(v))
	case float64:
		*m = *NewMatchExpression(MatchIn, strconv.FormatFloat(v, 'f', -1, 64))
	case []interface{}:
		values := make([]string, len(v))
		for i, value := range v {
			str, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid value %v in %v", value, v)
			}
			values[i] = str
		}
		*m = *NewMatchExpression(MatchIn, values...)
	case map[string]interface{}:
		op, ok := v["op"]
		if !ok {
			return fmt.Errorf("missing 'op' field in %v", v)
		}
		opStr, ok := op.(string)
		if !ok {
			return fmt.Errorf("invalid 'op' %v in %v", op, v)
		}
		m.Op = MatchOp(opStr)

		if values, ok := v["values"]; ok {
			vSlice, ok := values.([]interface{})
			if !ok {
				return fmt.Errorf("invalid 'values' %v in %v", values, v)
			}

			m.Values = make([]string, len(vSlice))
			for i, value := range vSlice {
				str, ok := value.(string)
				if !ok {
					return fmt.Errorf("invalid value %v in %v", value, v)
				}
				m.Values[i] = str
			}
		}
	default:
		return fmt.Errorf("invalid rule '%v' (%T)", v, v)
	}

	return m.Validate()
}

func (m *MatchExpressionSet) MatchKeys(keys map[string]struct{}) (bool, error) {
	for n, e := range *m {
		match, err := e.MatchKeys(n, keys)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func (m *MatchExpressionSet) MatchValues(values map[string]string) (bool, error) {
	for n, e := range *m {
		match, err := e.MatchValues(n, values)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func (m *MatchExpressionSet) UnmarshalJSON(data []byte) error {
	*m = make(MatchExpressionSet)

	names := make([]string, 0)
	if err := json.Unmarshal(data, &names); err == nil {
		// Simplified slice form
		for _, name := range names {
			split := strings.SplitN(name, "=", 2)
			if len(split) == 1 {
				(*m)[split[0]] = NewMatchExpression(MatchExists)
			} else {
				(*m)[split[0]] = NewMatchExpression(MatchIn, split[1])
			}
		}
	} else {
		// Unmarshal the full map form
		klog.V(2).Info("failed ot unmarshal to simple MatchExpression slice, trying map...")
		expressions := make(map[string]*MatchExpression)
		if err := json.Unmarshal(data, &expressions); err != nil {
			return err
		} else {
			for k, v := range expressions {
				if v != nil {
					(*m)[k] = v
				} else {
					(*m)[k] = NewMatchExpression(MatchExists)
				}
			}
		}
	}

	return nil
}
