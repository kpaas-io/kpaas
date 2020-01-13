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
	"sync"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/service/model/common"
)

type (
	Node struct {
		ConnectionData

		Name                string                                    // node name
		Description         string                                    // node description
		MachineRoles        []constant.MachineRole                    // machine role, like: master, worker, etcd. Master and worker roles are mutually exclusive.
		Labels              []*Label                                  // Node labels
		Taints              []*Taint                                  // Node taints
		CheckReport         *CheckReport                              // Check node report
		DeploymentReports   map[constant.DeployItem]*DeploymentReport // Deployment report for each role
		DockerRootDirectory string                                    // Docker Root Directory
		rwLock              sync.RWMutex                              // Read write lock
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
		DeployItem constant.DeployItem
		Status     DeployStatus
		Error      *common.FailureDetail
	}

	CheckReport struct {
		CheckItems   []*CheckItem          // Check item list
		CheckResult  constant.CheckResult  // Overall inspection status
		CheckedError *common.FailureDetail // Checked failure detail
	}

	CheckItem struct {
		ItemName    string // Check Item Name
		CheckResult constant.CheckResult
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

	AuthenticationType string // Type of authorization,  password or privateKey

	TaintEffect string // Taint Effect, NoSchedule, NoExecute or PreferNoSchedule

	DeployStatus string // Deploy node status
)

const (
	AuthenticationTypePassword   AuthenticationType = "password"   // Use Password to authorize
	AuthenticationTypePrivateKey AuthenticationType = "privateKey" // Use RSA PrivateKey to authorize

	TaintEffectNoSchedule       TaintEffect = "NoSchedule"
	TaintEffectNoExecute        TaintEffect = "NoExecute"
	TaintEffectPreferNoSchedule TaintEffect = "PreferNoSchedule"

	DeployStatusPending    DeployStatus = "pending"
	DeployStatusRunning    DeployStatus = "running"
	DeployStatusSuccessful DeployStatus = "successful"
	DeployStatusFailed     DeployStatus = "failed"
	DeployStatusAborted    DeployStatus = "aborted"

	DefaultDockerRootDirectory = "/var/lib/docker"
	DefaultUsername            = "root"
)

func NewNode() *Node {

	node := new(Node)
	node.init()
	return node
}

func (node *Node) init() {

	node.MachineRoles = make([]constant.MachineRole, 0, 2)
	node.initDeploymentReports()
	node.Labels = make([]*Label, 0, 0)
	node.Taints = make([]*Taint, 0, 0)
	node.ConnectionData.Port = uint16(22)
	node.ConnectionData.Username = DefaultUsername
	node.ConnectionData.AuthenticationType = AuthenticationTypePassword
	node.DockerRootDirectory = DefaultDockerRootDirectory
	node.CheckReport = new(CheckReport)
	node.CheckReport.init()
	node.rwLock = sync.RWMutex{}
}

func (node *Node) initDeploymentReports() {
	node.DeploymentReports = make(map[constant.DeployItem]*DeploymentReport)
}

func (node *Node) SetCheckResult(result constant.CheckResult, detail *common.FailureDetail) {

	node.rwLock.Lock()
	defer node.rwLock.Unlock()

	node.CheckReport.CheckResult = result
	if detail != nil {
		node.CheckReport.CheckedError = detail.Clone()
	}
}

func (node *Node) SetCheckItem(itemName string, result constant.CheckResult, detail *common.FailureDetail) {

	node.rwLock.Lock()
	defer node.rwLock.Unlock()

	item := NewCheckItem()
	var isFound bool
	for _, iterateItem := range node.CheckReport.CheckItems {

		if iterateItem.ItemName == itemName {
			item = iterateItem
			isFound = true
		}
	}

	if !isFound {
		item.ItemName = itemName
		node.CheckReport.CheckItems = append(node.CheckReport.CheckItems, item)
	}

	item.CheckResult = result
	item.Error = detail
}

func (node *Node) SetDeployResult(deployItem constant.DeployItem, status DeployStatus, detail *common.FailureDetail) {

	node.rwLock.Lock()
	defer node.rwLock.Unlock()

	if _, exist := node.DeploymentReports[deployItem]; !exist {
		node.DeploymentReports[deployItem] = NewDeploymentReport()
		node.DeploymentReports[deployItem].DeployItem = deployItem
	}

	node.DeploymentReports[deployItem].Status = status
	node.DeploymentReports[deployItem].Error = detail
}

func (node *Node) IsMatchMachineRole(role constant.MachineRole) bool {

	node.rwLock.RLock()
	defer node.rwLock.RUnlock()

	for _, iterateRole := range node.MachineRoles {
		if iterateRole == role {
			return true
		}
	}

	return false
}

func NewDeploymentReport() *DeploymentReport {

	report := new(DeploymentReport)
	report.init()
	return report
}

func (report *DeploymentReport) init() {

	report.Status = DeployStatusPending
}

func NewCheckItem() *CheckItem {

	item := new(CheckItem)
	item.init()
	return item
}

func (item *CheckItem) init() {

	item.CheckResult = constant.CheckResultPending
}

func (report *CheckReport) init() {

	report.CheckResult = constant.CheckResultPending
	report.CheckedError = nil
	report.CheckItems = make([]*CheckItem, 0, 0)
}
