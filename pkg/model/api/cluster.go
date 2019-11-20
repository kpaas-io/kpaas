// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

type (
	Cluster struct {
		ShortName                string                   `json:"shortName" binding:"required" minLength:"1" maxLength:"20"`
		Name                     string                   `json:"name" binding:"required"`
		KubeAPIServerConnectType KubeAPIServerConnectType `json:"kubeAPIServerConnectType" binding:"required" enums:"firstMasterIP,keepalived,loadbalancer"`
		VIP                      string                   `json:"vip" maxLength:"15"`
		NetInterfaceName         string                   `json:"netInterfaceName" maxLength:"30"`
		NodePortMinimum          uint16                   `json:"nodePort" minimum:"1" default:"30000"`
		NodePortMaximum          uint16                   `json:"nodePort" maximum:"65535" default:"32767"`
		Labels                   []Label                  `json:"labels"`
		Annotations              []Annotation             `json:"annotations"`
	}

	KubeAPIServerConnectType string

	Label struct {
		Key   string `json:"key" binding:"required" minimum:"1" maximum:"253"`
		Value string `json:"value" binding:"required" minimum:"1"`
	}

	Annotation struct {
		Key   string `json:"key" binding:"required" minimum:"1" maximum:"253"`
		Value string `json:"value" binding:"required" minimum:"1"`
	}
)

const (
	KubeAPIServerConnectTypeFirstMasterIP KubeAPIServerConnectType = "firstMasterIP"
	KubeAPIServerConnectTypeKeepalived    KubeAPIServerConnectType = "keepalived"
	KubeAPIServerConnectTypeLoadBalancer  KubeAPIServerConnectType = "loadbalancer"
)
