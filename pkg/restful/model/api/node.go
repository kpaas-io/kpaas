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
	NodeData struct {
		Name           string      `json:"name" binding:"required" minLength:"1" maxLength:"64"` // node name
		Description    string      `json:"description"`                                          // node description
		MachineRole    MachineRole `json:"role" default:"master" enums:"master,worker,etcd"`     // machine role, like: master, worker, etcd. Master and worker roles are mutually exclusive.
		Labels         []Label     `json:"labels"`                                               // Node labels
		Taints         []Taint     `json:"taint"`                                                // Node taints
		ConnectionData `json:",inline"`
	}

	ConnectionData struct {
		IP                 string             `json:"ip" binding:"required" minLength:"1" maxLength:"15"`               // node ip
		Port               uint16             `json:"port" binding:"required" minimum:"1" maximum:"65535" default:"22"` // ssh port
		Username           string             `json:"username" binding:"required" maxLength:"128"`                      // ssh username
		AuthenticationType AuthenticationType `json:"authorizationType" enums:"password,privateKey"`                    // type of authorization
		Password           string             `json:"password"`                                                         // login password
		PrivateKeyName     string             `json:"privateKeyName"`                                                   // the private key name of login
	}

	Taint struct {
		Key    string      `json:"key" binding:"required" minimum:"1" maximum:"63"`
		Value  string      `json:"value" binding:"required" minimum:"1" maximum:"63"`
		Effect TaintEffect `json:"effect" enums:"NoSchedule,NoExecute,PreferNoSchedule"`
	}

	MachineRole string // Machine Role, master or worker

	AuthenticationType string // Type of authorization,  password or privateKey

	TaintEffect string // Taint Effect, NoSchedule, NoExecute or PreferNoSchedule
)

const (
	AuthenticationTypePassword   AuthenticationType = "password"   // Use Password to authorize
	AuthenticationTypePrivateKey AuthenticationType = "privateKey" // Use RSA PrivateKey to authorize

	MachineRoleMaster MachineRole = "master" // master node
	MachineRoleWorker MachineRole = "worker" // worker node
	MachineRoleEtcd   MachineRole = "etcd"   // etcd node

	TaintEffectNoSchedule       TaintEffect = "NoSchedule"
	TaintEffectNoExecute        TaintEffect = "NoExecute"
	TaintEffectPreferNoSchedule TaintEffect = "PreferNoSchedule"
)
