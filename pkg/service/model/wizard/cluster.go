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
	"errors"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/service/model/common"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/idcreator"
)

type (
	Cluster struct {
		ClusterId          uint64
		Info               *ClusterInfo
		Nodes              []*Node
		DeploymentStatus   DeployClusterStatus
		DeployClusterError *common.FailureDetail
		ClusterCheckResult constant.CheckResult
		CheckClusterError  *common.FailureDetail
		Wizard             *WizardData
		KubeConfig         *string
		lock               *sync.RWMutex
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

	DefaultNodePortMinimum uint16 = 30000
	DefaultNodePortMaximum uint16 = 32767
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
	cluster.ClusterCheckResult = constant.CheckResultNotRunning
	cluster.Nodes = make([]*Node, 0, 0)
	cluster.Wizard = NewWizardData()
	cluster.lock = &sync.RWMutex{}
	var err error
	cluster.ClusterId, err = idcreator.NextID()
	if err != nil {
		logrus.Error(err)
	}
	cluster.KubeConfig = new(string)
}

func (cluster *Cluster) GetCheckResult() constant.CheckResult {

	cluster.lock.RLock()
	defer cluster.lock.RUnlock()

	return cluster.ClusterCheckResult
}

func (cluster *Cluster) AddNode(node *Node) error {

	cluster.lock.Lock()
	defer cluster.lock.Unlock()

	for _, iterateNode := range cluster.Nodes {

		if iterateNode.IP == node.IP {

			return h.EExists.WithPayload("node ip was exist")
		}

		if iterateNode.Name == node.Name {

			return h.EExists.WithPayload("node name was exist")
		}
	}

	wizardData.Nodes = append(wizardData.Nodes, node)
	return nil
}

func (cluster *Cluster) UpdateNode(node *Node) error {

	cluster.lock.Lock()
	defer cluster.lock.Unlock()

	var targetNode *Node

	for _, iterateNode := range cluster.Nodes {

		if iterateNode.IP == node.IP {

			targetNode = iterateNode
			break
		}
	}

	if targetNode == nil {
		return h.ENotFound.WithPayload("node ip not exist")
	}

	for _, iterateNode := range cluster.Nodes {

		if iterateNode.Name == node.Name && iterateNode != targetNode {

			return h.EExists.WithPayload("node name was exist")
		}
	}

	targetNode.Name = node.Name
	targetNode.Description = node.Description
	targetNode.DockerRootDirectory = node.DockerRootDirectory
	targetNode.MachineRoles = node.MachineRoles
	targetNode.Labels = node.Labels
	targetNode.Taints = node.Taints
	targetNode.ConnectionData.IP = node.ConnectionData.IP
	targetNode.ConnectionData.Port = node.ConnectionData.Port
	targetNode.ConnectionData.Username = node.ConnectionData.Username
	targetNode.ConnectionData.AuthenticationType = node.ConnectionData.AuthenticationType
	targetNode.ConnectionData.PrivateKeyName = node.ConnectionData.PrivateKeyName
	if len(node.ConnectionData.Password) != 0 {
		targetNode.ConnectionData.Password = node.ConnectionData.Password
	}

	return nil
}

func (cluster *Cluster) DeleteNode(ip string) error {

	cluster.lock.Lock()
	defer cluster.lock.Unlock()
	var found bool

	newList := make([]*Node, 0, len(wizardData.Nodes)-1)

	for index, iterateNode := range wizardData.Nodes {

		if iterateNode.IP != ip {
			continue
		}

		if index > 0 {
			newList = append(newList, wizardData.Nodes[:index]...)
		}
		if index < len(wizardData.Nodes)-1 {
			newList = append(newList, wizardData.Nodes[index+1:]...)
		}

		found = true
		break
	}

	if !found {
		return h.EExists.WithPayload("node not exist")
	}

	cluster.Nodes = newList

	return nil
}

func (cluster *Cluster) GetNode(ip string) *Node {

	for _, node := range cluster.Nodes {
		if node.IP == ip {
			return node
		}
	}

	return nil
}

func (cluster *Cluster) GetNodeByName(name string) *Node {

	for _, node := range cluster.Nodes {
		if node.Name == name {
			return node
		}
	}

	return nil
}

func (cluster *Cluster) MarkNodeChecking() error {

	cluster.lock.Lock()
	defer cluster.lock.Unlock()

	if len(cluster.Nodes) <= 0 {
		return nil
	}

	if cluster.ClusterCheckResult == constant.CheckResultChecking {
		return errors.New("was checking")
	}

	cluster.ClusterCheckResult = constant.CheckResultChecking

	return nil
}

func (cluster *Cluster) ClearClusterCheckingData() error {

	cluster.lock.Lock()
	defer cluster.lock.Unlock()

	if len(cluster.Nodes) <= 0 {
		return nil
	}

	cluster.ClusterCheckResult = constant.CheckResultNotRunning
	cluster.CheckClusterError = nil

	for _, node := range cluster.Nodes {

		node.CheckReport.init()
	}

	return nil
}

func (cluster *Cluster) SetClusterCheckResult(result constant.CheckResult, failureDetail *common.FailureDetail) {

	cluster.lock.Lock()
	defer cluster.lock.Unlock()

	cluster.ClusterCheckResult = result
	if failureDetail != nil {
		cluster.CheckClusterError = failureDetail.Clone()
	}
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
	info.NodePortMinimum = DefaultNodePortMinimum
	info.NodePortMaximum = DefaultNodePortMaximum
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
