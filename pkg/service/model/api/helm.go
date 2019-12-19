// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

type (
	HelmValues map[string]interface{}
	// HelmRelease presents a helm release.
	HelmRelease struct {
		Cluster string `json:"cluster"`
		// empty Name means to generate a name, only used in installing a release.
		Name         string     `json:"name,omitempty"`
		Namespace    string     `json:"namespace"`
		Chart        string     `json:"chart,omitempty"`
		ChartRepo    string     `json:"chartRepo,omitempty"`
		ChartVersion string     `json:"chartVersion,omitempty"`
		Values       HelmValues `json:"values,omitempty"`
		Revision     uint32     `json:"revision,omitempty"`
		Manifest     string     `json:"manifest,omitempty"`
	}
)
