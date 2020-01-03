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

package wizard

import "github.com/kpaas-io/kpaas/pkg/deploy/protos"

var DefaultNetworkOptions protos.NetworkOptions = protos.NetworkOptions{
	NetworkType: "calico",
	CalicoOptions: &protos.CalicoOptions{
		EncapsulationMode: "vxlan",
		VxlanPort:         4789,
		InitialPodIps:     "",
		VethMtu:           1400,
		IpDetectionMethod: "fromKubernetes",
	},
}

func (cluster *Cluster) SetNetworkOptions(options *protos.NetworkOptions) {
	cluster.lock.Lock()
	defer cluster.lock.Unlock()

	cluster.NetworkOptions = options
}

func (cluster *Cluster) GetNetworkOptions() *protos.NetworkOptions {
	cluster.lock.Lock()
	defer cluster.lock.Unlock()
	if cluster.NetworkOptions == nil {
		return &DefaultNetworkOptions
	}
	return cluster.NetworkOptions
}
