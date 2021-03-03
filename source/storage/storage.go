/*
Copyright 2018-2021 The Kubernetes Authors.

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

package storage

import (
	"fmt"
	"io/ioutil"

	"sigs.k8s.io/node-feature-discovery/pkg/api/feature"
	"sigs.k8s.io/node-feature-discovery/pkg/utils"
	"sigs.k8s.io/node-feature-discovery/source"
)

const Name = "storage"

const (
	DiskFeature = "disk"

	DiskFeatureNonrotationalPresent = "nonrotational.present"
)

// storageSource implements the FeatureSource and LabelSource interfaces.
type storageSource struct {
	features *feature.DomainFeatures
}

// Singleton source instance
var (
	src storageSource
	_   source.FeatureSource = &src
	_   source.LabelSource   = &src
)

// Name returns an identifier string for this feature source.
func (s *storageSource) Name() string { return Name }

// Priority method of the LabelSource interface
func (s *storageSource) Priority() int { return 0 }

// GetLabels method of the LabelSource interface
func (s *storageSource) GetLabels() (source.FeatureLabels, error) {
	features := source.FeatureLabels{}

	if s.features.Values[DiskFeature].Features[DiskFeatureNonrotationalPresent] == "true" {
		features["nonrotationaldisk"] = true
	}

	return features, nil
}

// Discover method of the FeatureSource interface
func (s *storageSource) Discover() error {
	s.features = feature.NewDomainFeatures()

	// Check if there is any non-rotational block devices attached to the node
	blockdevices, err := ioutil.ReadDir(source.SysfsDir.Path("block"))
	if err != nil {
		return fmt.Errorf("failed to list block devices: %w", err)
	}

	disk := feature.NewValueFeatures()
	disk.Features[DiskFeatureNonrotationalPresent] = "false"
	for _, bdev := range blockdevices {
		fname := source.SysfsDir.Path("block", bdev.Name(), "queue/rotational")
		bytes, err := ioutil.ReadFile(fname)
		if err != nil {
			return fmt.Errorf("can't read rotational status: %w", err)
		}
		if bytes[0] == byte('0') {
			// Non-rotational storage is present, add label.
			disk.Features[DiskFeatureNonrotationalPresent] = "true"
			break
		}
	}
	s.features.Values[DiskFeature] = *disk

	utils.KlogDump(3, "discovered system features:", "  ", s.features)

	return nil
}

// GetFeatures method of the FeatureSource Interface.
func (s *storageSource) GetFeatures() *feature.DomainFeatures {
	return s.features
}

func init() {
	source.Register(&src)
}
