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

import (
	"regexp"

	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/utils/validator"
)

type (
	NodeBaseData struct {
		Name                string                 `json:"name" binding:"required" minLength:"1" maxLength:"64"` // node name
		Description         string                 `json:"description"`                                          // node description
		MachineRoles        []constant.MachineRole `json:"roles" default:"" enums:"master,worker,etcd"`          // machine role, Master and worker roles are mutually exclusive.
		Labels              []Label                `json:"labels"`                                               // Node labels
		Taints              []Taint                `json:"taints"`                                               // Node taints
		DockerRootDirectory string                 `json:"dockerRootDirectory" default:"/var/lib/docker"`        // Docker Root Directory
	}

	NodeData struct {
		NodeBaseData   `json:",inline"`
		ConnectionData `json:",inline"`
	}

	ConnectionData struct {
		SSHLoginData `json:",inline"`

		IP   string `json:"ip" binding:"required" minLength:"1" maxLength:"15"`               // node ip
		Port uint16 `json:"port" binding:"required" minimum:"1" maximum:"65535" default:"22"` // ssh port
	}

	UpdateNodeData struct {
		NodeBaseData `json:",inline"`
		SSHLoginData `json:",inline"`

		Port uint16 `json:"port" binding:"required" minimum:"1" maximum:"65535" default:"22"` // ssh port
	}

	SSHLoginData struct {
		Username           string             `json:"username" binding:"required" maxLength:"128"`   // ssh username
		AuthenticationType AuthenticationType `json:"authorizationType" enums:"password,privateKey"` // type of authorization
		Password           string             `json:"password,omitempty"`                            // login password
		PrivateKeyName     string             `json:"privateKeyName,omitempty"`                      // the private key name of login
	}

	Taint struct {
		Key    string      `json:"key" binding:"required" minimum:"1" maximum:"63"`
		Value  string      `json:"value" binding:"required" minimum:"1" maximum:"63"`
		Effect TaintEffect `json:"effect" enums:"NoSchedule,NoExecute,PreferNoSchedule"`
	}

	GetNodeListResponse struct {
		Nodes []NodeData `json:"nodes"` // node list
	}

	AuthenticationType string // Type of authorization,  password or privateKey

	TaintEffect string // Taint Effect, NoSchedule, NoExecute or PreferNoSchedule
)

const (
	AuthenticationTypePassword   AuthenticationType = "password"   // Use Password to authorize
	AuthenticationTypePrivateKey AuthenticationType = "privateKey" // Use RSA PrivateKey to authorize

	TaintEffectNoSchedule       TaintEffect = "NoSchedule"
	TaintEffectNoExecute        TaintEffect = "NoExecute"
	TaintEffectPreferNoSchedule TaintEffect = "PreferNoSchedule"

	NodeNameLengthLimit           = 64
	NodeDescriptionLengthLimit    = 100
	TaintKeyLengthLimit           = 253
	TaintValueLengthLimit         = 63
	NodeUsernameRegularExpression = `^[A-Za-z]([\w\-.]+)?$`

	NodeSSHPortMinimum = 1
	NodeSSHPortMaximum = 65535
)

func (node *NodeBaseData) Validate() error {

	rolesNames := make([]string, 0, len(node.MachineRoles))

	for _, role := range node.MachineRoles {
		rolesNames = append(rolesNames, string(role))
	}

	wrapper := validator.NewWrapper(
		validator.ValidateString(node.Name, "name", validator.ItemNotEmptyLimit, NodeNameLengthLimit),
		validator.ValidateRegexp(regexp.MustCompile(`[A-Za-z][\w\-]*\w?`), node.Name, "name"),
		validator.ValidateString(node.Description, "description", validator.ItemNoLimit, NodeDescriptionLengthLimit),
		validator.ValidateStringArrayOptions(rolesNames, "role", []string{string(constant.MachineRoleMaster), string(constant.MachineRoleWorker), string(constant.MachineRoleEtcd)}),
	)

	for _, label := range node.Labels {

		wrapper.AddValidateFunc(
			func() error {
				return label.Validate()
			},
		)
	}

	for _, taint := range node.Taints {

		wrapper.AddValidateFunc(
			func() error {
				return taint.Validate()
			},
		)
	}

	return wrapper.Validate()
}

func (login *SSHLoginData) Validate() error {

	wrapper := validator.NewWrapper(
		validator.ValidateString(login.Username, "username", validator.ItemNotEmptyLimit, validator.ItemNoLimit),
		validator.ValidateRegexp(regexp.MustCompile(NodeUsernameRegularExpression), login.Username, "username"),
		validator.ValidateStringOptions(string(login.AuthenticationType), "authorizationType", []string{string(AuthenticationTypePassword), string(AuthenticationTypePrivateKey)}),
	)

	switch login.AuthenticationType {
	case AuthenticationTypePassword:
		wrapper.AddValidateFunc(
			validator.ValidateString(login.Password, "password", validator.ItemNotEmptyLimit, validator.ItemNoLimit),
		)
	case AuthenticationTypePrivateKey:
		wrapper.AddValidateFunc(
			validator.ValidateString(login.PrivateKeyName, "privateKeyName", validator.ItemNotEmptyLimit, validator.ItemNoLimit),
		)
	}

	return wrapper.Validate()
}

func (node *ConnectionData) Validate() error {

	return validator.NewWrapper(
		validator.ValidateIP(node.IP, "ip"),
		validator.ValidateIntRange(int(node.Port), "port", NodeSSHPortMinimum, NodeSSHPortMaximum),
		func() error {
			return node.SSHLoginData.Validate()
		},
	).Validate()
}

func (node *NodeData) Validate() error {

	return validator.NewWrapper(
		func() error {
			return node.NodeBaseData.Validate()
		},
		func() error {
			return node.ConnectionData.Validate()
		},
	).Validate()
}

func (taint *Taint) Validate() error {

	return validator.NewWrapper(
		validator.ValidateString(taint.Key, "key", validator.ItemNotEmptyLimit, TaintKeyLengthLimit),
		validator.ValidateString(taint.Value, "value", validator.ItemNotEmptyLimit, TaintValueLengthLimit),
		ValidateStringFunctionReturnErrorMessages(validation.IsQualifiedName, taint.Key, "taint.key"),
		ValidateStringFunctionReturnErrorMessages(validation.IsValidLabelValue, taint.Value, "taint.value"),
		validator.ValidateStringOptions(string(taint.Effect), "taint.effect",
			[]string{string(TaintEffectNoExecute), string(TaintEffectNoSchedule), string(TaintEffectPreferNoSchedule)}),
	).Validate()
}

func (node *UpdateNodeData) Validate() error {

	return validator.NewWrapper(
		func() error {
			return node.NodeBaseData.Validate()
		},
		validator.ValidateIntRange(int(node.Port), "port", NodeSSHPortMinimum, NodeSSHPortMaximum),
		func() error {
			return node.SSHLoginData.Validate()
		},
	).Validate()
}
