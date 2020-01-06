// Copyright 2020 Shanghai JingDuo Information Technology co., Ltd.
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

type NetworkType string

const (
	NetworkTypeCalico NetworkType = "calico"
)

type NetworkOptions struct {
	NetworkType   NetworkType    `json:"networkType" enums:"calico"`
	CalicoOptions *CalicoOptions `json:"calicoOptions,omitempty"`
}

type EncapsulationMode string

const (
	EncapsulationVxlan = "vxlan"
	EncapsulationIpip  = "ipip"
	EncapsulationNone  = "none"
)

type IPDetectionMethod string

const (
	IPDetectionMethodInterface      = "interface"
	IPDetectionMethodFirstFound     = "first-found"
	IPDetectionMethodFromKubernetes = "from-kubernetes"
)

const (
	DefaultVxlanPort = 4789
)

type CalicoOptions struct {
	EncapsulationMode    EncapsulationMode `json:"encapsulationMode,omitempty" enums:"vxlan, ipip, none"`
	VxlanPort            int               `json:"vxlanPort,omitempty"`
	InitialPodIPs        string            `json:"initialPodIPs,omitempty"`
	VethMtu              int               `json:"vethMtu,omitempty"`
	IPDetectionMethod    IPDetectionMethod `json:"ipDetectionMethod,omitempty" enums:"from-kubernetes,first-found,interface"`
	IPDetectionInterface string            `json:"ipDetectionInterface,omitempty"`
}
