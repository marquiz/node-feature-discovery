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
	"text/template"

	"bytes"
	"fmt"
	"strings"
)

func (r *Rule) executeNameTemplate(data interface{}) (string, error) {
	n, err := r.getNameExpander()
	if err != nil {
		return "", err
	}
	return n.expand(data)
}

func (r *Rule) NameIsTemplate() (bool, error) {
	n, err := r.getNameExpander()
	if err != nil {
		return true, err
	}
	return n.nameTemplate != nil, nil
}

func (r *Rule) getNameExpander() (*nameTemplateHelper, error) {
	if r.name == nil {
		n, err := newNameTemplateHelper(r.Name)
		if err != nil {
			return nil, err
		}
		r.name = n
	}
	return r.name, nil
}

type nameTemplateHelper struct {
	name         string
	nameTemplate *template.Template
}

func newNameTemplateHelper(name string) (*nameTemplateHelper, error) {
	e := nameTemplateHelper{name: name}

	if strings.Contains(name, "{{") {

		tmpl, err := template.New("").Option("missingkey=error").Parse(name)
		if err != nil {
			return nil, fmt.Errorf("invalid template in rule name: %w", err)
		}
		e.nameTemplate = tmpl
	}
	return &e, nil
}

// DeepCopy is a stub to augment the auto-generated code
func (in *nameTemplateHelper) DeepCopy() *nameTemplateHelper {
	if in == nil {
		return nil
	}
	out := new(nameTemplateHelper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is a stub to augment the auto-generated code
func (in *nameTemplateHelper) DeepCopyInto(out *nameTemplateHelper) {
	// HACK: just re-use the template
	out.nameTemplate = in.nameTemplate
}

func (e *nameTemplateHelper) expand(data interface{}) (string, error) {
	if e.nameTemplate == nil {
		return e.name, nil
	}

	var tmp bytes.Buffer
	if err := e.nameTemplate.Execute(&tmp, data); err != nil {
		return "", err
	}
	return tmp.String(), nil
}
