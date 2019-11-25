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

type (
	Cluster struct {
		Info             *ClusterInfo
		Nodes            []*Node
		DeploymentStatus DeployClusterStatus
		Wizard           *WizardData
	}

	ClusterInfo struct {
		ShortName               string
		Name                    string
		KubeAPIServerConnection *KubeAPIServerConnectionData
		NodePortMinimum         uint16
		NodePortMaximum         uint16
		Labels                  []*Label
		Annotations             []*Annotation
	}

	KubeAPIServerConnectionData struct {
		KubeAPIServerConnectType KubeAPIServerConnectType
		VIP                      string
		NetInterfaceName         string
		LoadbalancerIP           string
		LoadbalancerPort         uint16
	}

	KubeAPIServerConnectType string

	DeployClusterStatus string
)

const (
	KubeAPIServerConnectTypeFirstMasterIP KubeAPIServerConnectType = "firstMasterIP"
	KubeAPIServerConnectTypeKeepalived    KubeAPIServerConnectType = "keepalived"
	KubeAPIServerConnectTypeLoadBalancer  KubeAPIServerConnectType = "loadbalancer"

	DeployClusterStatusNotRunning         DeployClusterStatus = "notRunning"
	DeployClusterStatusRunning            DeployClusterStatus = "running"
	DeployClusterStatusSuccessful         DeployClusterStatus = "successful"
	DeployClusterStatusFailed             DeployClusterStatus = "failed"
	DeployClusterStatusWorkedButHaveError DeployClusterStatus = "workedButHaveError"
)

var (
	wizardData *Cluster
)

func NewCluster() *Cluster {

	cluster := new(Cluster)
	cluster.init()
	return cluster
}

func (cluster *Cluster) init() {

	cluster.Info = NewClusterInfo()
	cluster.DeploymentStatus = DeployClusterStatusNotRunning
	cluster.Nodes = make([]*Node, 0, 0)
	cluster.Wizard = NewWizardData()
}

func NewClusterInfo() *ClusterInfo {

	info := new(ClusterInfo)
	info.init()
	return info
}

func (info *ClusterInfo) init() {

	info.KubeAPIServerConnection = NewKubeAPIServerConnectionData()
	info.Labels = make([]*Label, 0, 0)
	info.Annotations = make([]*Annotation, 0, 0)
}

func NewKubeAPIServerConnectionData() *KubeAPIServerConnectionData {

	data := new(KubeAPIServerConnectionData)
	data.init()
	return data
}

func (data *KubeAPIServerConnectionData) init() {

	data.KubeAPIServerConnectType = KubeAPIServerConnectTypeFirstMasterIP
}

func init() {

	ClearCurrentWizardData()
}

func GetCurrentWizard() *Cluster {

	return wizardData
}

func ClearCurrentWizardData() {

	wizardData = NewCluster()
}
