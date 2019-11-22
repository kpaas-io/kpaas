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

package wizard

import (
	"github.com/kpaas-io/kpaas/pkg/service/model/common"
)

type (
	Node struct {
		Name             string              // node name
		Description      string              // node description
		MachineRole      []MachineRole       // machine role, like: master, worker, etcd. Master and worker roles are mutually exclusive.
		Labels           []*Label            // Node labels
		Taints           []*Taint            // Node taints
		CheckReport      []*CheckReport      // Overall inspection status
		DeploymentReport []*DeploymentReport // Deployment report for each role
		ConnectionData
	}

	ConnectionData struct {
		IP                 string             // node ip
		Port               uint16             // ssh port
		Username           string             // ssh username
		AuthenticationType AuthenticationType // type of authorization
		Password           string             // login password
		PrivateKeyName     string             // the private key name of login
	}

	DeploymentReport struct {
		Role   MachineRole
		Status DeployStatus
		Error  *common.FailureDetail
	}

	CheckReport struct {
		Role       MachineRole
		CheckItems []*CheckItem
	}

	CheckItem struct {
		ItemName    string // Check Item Name
		CheckResult CheckResult
		Error       *common.FailureDetail
	}

	Annotation struct {
		Key   string
		Value string
	}

	Label struct {
		Key   string
		Value string
	}

	Taint struct {
		Key    string
		Value  string
		Effect TaintEffect
	}

	MachineRole string // Machine Role, master or worker

	AuthenticationType string // Type of authorization,  password or privateKey

	TaintEffect string // Taint Effect, NoSchedule, NoExecute or PreferNoSchedule

	CheckResult string // Check node result

	DeployStatus string // Deploy node status
)

const (
	CheckResultNotRunning CheckResult = "notRunning"
	CheckResultChecking   CheckResult = "checking"
	CheckResultPassed     CheckResult = "passed"
	CheckResultFailed     CheckResult = "failed"

	AuthenticationTypePassword   AuthenticationType = "password"   // Use Password to authorize
	AuthenticationTypePrivateKey AuthenticationType = "privateKey" // Use RSA PrivateKey to authorize

	MachineRoleMaster MachineRole = "master" // master node
	MachineRoleWorker MachineRole = "worker" // worker node
	MachineRoleEtcd   MachineRole = "etcd"   // etcd node

	TaintEffectNoSchedule       TaintEffect = "NoSchedule"
	TaintEffectNoExecute        TaintEffect = "NoExecute"
	TaintEffectPreferNoSchedule TaintEffect = "PreferNoSchedule"

	DeployStatusPending   DeployStatus = "pending"
	DeployStatusDeploying DeployStatus = "deploying"
	DeployStatusCompleted DeployStatus = "completed"
	DeployStatusFailed    DeployStatus = "failed"
	DeployStatusAborted   DeployStatus = "aborted"
)

func NewNode() *Node {

	node := new(Node)
	node.init()
	return node
}

func (node *Node) init() {

	node.MachineRole = make([]MachineRole, 0, 2)
	node.DeploymentReport = make([]*DeploymentReport, 0, 2)
	node.CheckReport = make([]*CheckReport, 0, 0)
	node.Labels = make([]*Label, 0, 0)
	node.Taints = make([]*Taint, 0, 0)
	node.ConnectionData.AuthenticationType = AuthenticationTypePassword
}

func NewDeploymentReport() *DeploymentReport {

	report := new(DeploymentReport)
	report.init()
	return report
}

func (report *DeploymentReport) init() {

	report.Status = DeployStatusPending
}

func NewCheckReport() *CheckReport {

	report := new(CheckReport)
	report.init()
	return report
}

func (report *CheckReport) init() {

	report.CheckItems = make([]*CheckItem, 0, 0)
}

func NewCheckItem() *CheckItem {

	item := new(CheckItem)
	item.init()
	return item
}

func (item *CheckItem) init() {

	item.CheckResult = CheckResultNotRunning
}
