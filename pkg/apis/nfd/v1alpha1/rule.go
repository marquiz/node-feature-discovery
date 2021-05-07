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

	"sigs.k8s.io/node-feature-discovery/pkg/api/feature"
	"sigs.k8s.io/node-feature-discovery/pkg/utils"
	"sigs.k8s.io/node-feature-discovery/source/cpu"
	"sigs.k8s.io/node-feature-discovery/source/kernel"
	"sigs.k8s.io/node-feature-discovery/source/pci"
	"sigs.k8s.io/node-feature-discovery/source/system"
	"sigs.k8s.io/node-feature-discovery/source/usb"
)

// Match exercises a Rule against a set of features.
func (r *Rule) Match(features map[string]*feature.DomainFeatures) (map[string]string, error) {
	ret := make(map[string]string)

	if len(r.MatchOn) > 0 {
		// Logical OR over the legacy rules
		matched := false
		for _, matchRule := range r.MatchOn {
			if m, err := matchRule.match(features); err != nil {
				return nil, err
			} else if m {
				matched = true

				if err := r.executeTemplate(nil, ret); err != nil {
					return nil, err
				}

				break
			}
		}
		if !matched {
			return nil, nil
		}
	}

	if len(r.MatchAny) > 0 {
		// Logical OR over the matchAny domains rules
		matched := false
		for _, matchRule := range r.MatchAny {
			if m, err := matchRule.match(features); err != nil {
				return nil, err
			} else if m != nil {
				matched = true
				utils.KlogDump(4, "matches for matchAny "+r.Name, "  ", m)

				if err := r.executeTemplate(m, ret); err != nil {
					return nil, err
				}

				if isTemplate, err := r.NameIsTemplate(); err != nil {
					return nil, nil
				} else if !isTemplate {
					// No templating so we stop here (further matches would just
					// produce the same labels)
					break
				}
			}
		}
		if !matched {
			return nil, nil
		}
	}

	if len(r.MatchAll) > 0 {
		// Logical AND over the matchAny domains rules
		for _, matchRule := range r.MatchAll {
			if m, err := matchRule.match(features); err != nil {
				return nil, err
			} else if m == nil {
				return nil, nil
			} else {

				utils.KlogDump(4, "matches for matchAll "+r.Name, "  ", m)
				if err := r.executeTemplate(m, ret); err != nil {
					return nil, err
				}
			}
		}
	}

	// We have a match
	return ret, nil
}

func (r *Rule) executeTemplate(in matchedFeatures, out map[string]string) error {
	expandedName, err := r.ExpandName(in)
	if err != nil {
		return nil
	}

	// Split out individual labels
	for _, item := range strings.Split(expandedName, "\n") {
		// Remove leading/trailing whitespace and skip empty lines
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			n, v := r.getNameValue(trimmed)
			out[n] = v
		}
	}
	return nil
}

func (r *Rule) getNameValue(name string) (string, string) {
	// Value can be overridden in Name with "key=value". This is useful for
	// templates.
	if strings.ContainsRune(name, '=') {
		split := strings.SplitN(name, "=", 2)
		return split[0], split[1]
	}

	if r.Value != nil {
		return name, *r.Value
	}

	return name, "true"
}

func (r *Rule) ExpandName(data interface{}) (string, error) {
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

type matchedFeatures map[string]domainMatchedFeatures

type domainMatchedFeatures map[string]interface{}

func (r *MatchRule) match(features map[string]*feature.DomainFeatures) (map[string]domainMatchedFeatures, error) {
	ret := make(matchedFeatures, len(*r))

	for domain, rules := range *r {
		domainFeatures, ok := features[domain]
		if !ok {
			return nil, fmt.Errorf("unknown feature source/domain %q", domain)
		}

		// Matched features
		matched := make(domainMatchedFeatures)

		for featureName, featureRules := range rules {
			var m bool
			var e error

			// Ignore case
			featureName = strings.ToLower(featureName)

			if f, ok := domainFeatures.Keys[featureName]; ok {
				v, err := featureRules.MatchGetKeys(f.Features)
				m = len(v) > 0
				e = err
				matched[featureName] = v
			} else if f, ok := domainFeatures.Values[featureName]; ok {
				v, err := featureRules.MatchGetValues(f.Features)
				m = len(v) > 0
				e = err
				matched[featureName] = v
			} else if f, ok := domainFeatures.Instances[featureName]; ok {
				v, err := featureRules.MatchGetInstances(f.Features)
				m = len(v) > 0
				e = err
				matched[featureName] = v
			} else {
				return nil, fmt.Errorf("%q feature of source/domain %q not available", featureName, domain)
			}
			if e != nil {
				return nil, e
			} else if !m {
				return nil, nil
			}
		}
		ret[domain] = matched
	}
	return ret, nil
}

func (rules *LegacyRule) match(features map[string]*feature.DomainFeatures) (bool, error) {
	if rules.CpuID != nil {
		if f, ok := features[cpu.Name].Keys[cpu.CpuidFeature]; !ok {
			return false, fmt.Errorf("cpuid information not available")
		} else if match, err := rules.CpuID.MatchKeys(f.Features); !match || err != nil {
			return match, err
		}
	}

	if rules.LoadedKMod != nil {
		if f, ok := features[kernel.Name].Keys[kernel.LoadedModuleFeature]; !ok {
			return false, fmt.Errorf("information about loaded modules not available")
		} else if match, err := rules.LoadedKMod.MatchKeys(f.Features); !match || err != nil {
			return match, err
		}
	}

	if rules.Kconfig != nil {
		if f, ok := features[kernel.Name].Values[kernel.ConfigFeature]; !ok {
			return false, fmt.Errorf("kernel config options not available")
		} else if match, err := rules.Kconfig.MatchValues(f.Features); !match || err != nil {
			return match, err
		}
	}

	if rules.PciID != nil {
		if f, ok := features[pci.Name].Instances[pci.DeviceFeature]; !ok {
			return false, fmt.Errorf("pci device information not available")
		} else if match, err := rules.PciID.MatchInstances(f.Features); !match || err != nil {
			return match, err
		}
	}

	if rules.UsbID != nil {
		if f, ok := features[usb.Name].Instances[usb.DeviceFeature]; !ok {
			return false, fmt.Errorf("usb device information not available")
		} else if match, err := rules.UsbID.MatchInstances(f.Features); !match || err != nil {
			return match, err
		}
	}

	if rules.Nodename != nil {
		if f, ok := features[system.Name].Values[system.NameFeature]; !ok {
			return false, fmt.Errorf("system name information not available")
		} else {
			n, ok := f.Features["nodename"]
			if match, err := rules.Nodename.MatchExpression.Match(ok, n); !match || err != nil {
				return match, err
			}
		}
	}
	return true, nil
}

func (r *Rule) getNameExpander() (*nameExpander, error) {
	if r.name == nil {
		n, err := newNameExpander(r.Name)
		if err != nil {
			return nil, err
		}
		r.name = n
	}
	return r.name, nil
}

type nameExpander struct {
	name         string
	nameTemplate *template.Template
}

func newNameExpander(name string) (*nameExpander, error) {
	e := nameExpander{name: name}

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
func (in *nameExpander) DeepCopy() *nameExpander {
	if in == nil {
		return nil
	}
	out := new(nameExpander)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is a stub to augment the auto-generated code
func (in *nameExpander) DeepCopyInto(out *nameExpander) {
	// HACK: just re-use the template
	out.nameTemplate = in.nameTemplate
}

func (e *nameExpander) expand(data interface{}) (string, error) {
	if e.nameTemplate == nil {
		return e.name, nil
	}

	var tmp bytes.Buffer
	if err := e.nameTemplate.Execute(&tmp, data); err != nil {
		return "", err
	}
	return tmp.String(), nil
}
