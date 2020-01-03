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

import (
	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
)

var DefaultNetworkOptions api.NetworkOptions = api.NetworkOptions{
	NetworkType: api.NetworkTypeCalico,
	CalicoOptions: &api.CalicoOptions{
		EncapsulationMode: api.EncapsulationVxlan,
		VxlanPort:         api.DefaultVxlanPort,
		InitialPodIPs:     constant.DefaultPodSubnet,
		VethMtu:           1400,
		IPDetectionMethod: api.IPDetectionMethodFromKubernetes,
	},
}

func (cluster *Cluster) SetNetworkOptions(options *api.NetworkOptions) {
	cluster.lock.Lock()
	defer cluster.lock.Unlock()
	// TODO: process situation where options == nil specially?
	cluster.NetworkOptions = options
}

func (cluster *Cluster) GetNetworkOptions() *api.NetworkOptions {
	cluster.lock.RLock()
	defer cluster.lock.RUnlock()
	if cluster.NetworkOptions == nil {
		// TODO: what to return if cluster.NetworkOptions == nil?
		return &DefaultNetworkOptions
	}
	return cluster.NetworkOptions
}
