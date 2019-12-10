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
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/service/model/common"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

func TestNewCluster(t *testing.T) {

	assert.IsType(t, &Cluster{}, NewCluster())
}

func TestNewClusterInfo(t *testing.T) {

	assert.IsType(t, &ClusterInfo{}, NewClusterInfo())
}

func TestGetCurrentWizard(t *testing.T) {

	cluster := GetCurrentWizard()
	assert.Equal(t, DeployClusterStatusNotRunning, cluster.DeployClusterStatus)
	assert.NotNil(t, cluster.Wizard)
	assert.NotNil(t, cluster.Nodes)
	assert.NotNil(t, cluster.Info)
	assert.Equal(t, ProgressSettingClusterInformation, cluster.Wizard.Progress)
	assert.Equal(t, WizardModeNormal, cluster.Wizard.WizardMode)
}

func TestClearCurrentWizardData(t *testing.T) {

	var cluster *Cluster
	cluster = GetCurrentWizard()
	cluster.Wizard.Progress = ProgressSettingNodesInformation
	cluster = GetCurrentWizard()
	assert.Equal(t, ProgressSettingNodesInformation, cluster.Wizard.Progress)

	ClearCurrentWizardData()
	cluster = GetCurrentWizard()
	assert.Equal(t, ProgressSettingClusterInformation, cluster.Wizard.Progress)
}

func TestNewKubeAPIServerConnectionData(t *testing.T) {

	data := NewKubeAPIServerConnectionData()
	assert.IsType(t, &KubeAPIServerConnectionData{}, data)
	assert.Equal(t, KubeAPIServerConnectTypeFirstMasterIP, data.KubeAPIServerConnectType)
}

func TestCluster_GetCheckResult(t *testing.T) {

	tests := []struct {
		Input Cluster
		Want  constant.CheckResult
	}{
		{
			Input: Cluster{
				ClusterCheckResult: constant.CheckResultPassed,
				lock:               new(sync.RWMutex),
			},
			Want: constant.CheckResultPassed,
		},
		{
			Input: Cluster{
				ClusterCheckResult: constant.CheckResultFailed,
				lock:               new(sync.RWMutex),
			},
			Want: constant.CheckResultFailed,
		},
	}

	for _, item := range tests {

		assert.Equal(t, item.Want, item.Input.GetCheckResult())
	}
}

func TestCluster_GetDeployClusterStatus(t *testing.T) {

	tests := []struct {
		Input Cluster
		Want  DeployClusterStatus
	}{
		{
			Input: Cluster{
				DeployClusterStatus: DeployClusterStatusSuccessful,
				lock:                new(sync.RWMutex),
			},
			Want: DeployClusterStatusSuccessful,
		},
		{
			Input: Cluster{
				DeployClusterStatus: DeployClusterStatusNotRunning,
				lock:                new(sync.RWMutex),
			},
			Want: DeployClusterStatusNotRunning,
		},
	}

	for _, item := range tests {

		assert.Equal(t, item.Want, item.Input.GetDeployClusterStatus())
	}
}

func TestCluster_AddNode(t *testing.T) {

	tests := []struct {
		Input struct {
			Cluster Cluster
			Node    *Node
		}
		Want struct {
			Cluster     Cluster
			ReturnValue error
		}
	}{
		{
			Input: struct {
				Cluster Cluster
				Node    *Node
			}{
				Cluster: Cluster{
					Nodes: []*Node{},
					lock:  new(sync.RWMutex),
				},
				Node: &Node{
					Name: "node1",
				},
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
						},
					},
					lock: new(sync.RWMutex),
				},
				ReturnValue: nil,
			},
		},
		{
			Input: struct {
				Cluster Cluster
				Node    *Node
			}{
				Cluster: Cluster{
					Nodes: []*Node{},
					lock:  new(sync.RWMutex),
				},
				Node: &Node{
					Name: "node2",
					ConnectionData: ConnectionData{
						IP:                 "192.168.31.1",
						Port:               22,
						Username:           "root",
						AuthenticationType: AuthenticationTypePassword,
						Password:           "123456",
					},
				},
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP:                 "192.168.31.1",
								Port:               22,
								Username:           "root",
								AuthenticationType: AuthenticationTypePassword,
								Password:           "123456",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				ReturnValue: nil,
			},
		},
		{
			Input: struct {
				Cluster Cluster
				Node    *Node
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				Node: &Node{
					Name: "node3",
					ConnectionData: ConnectionData{
						IP: "192.168.31.1",
					},
				},
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				ReturnValue: h.EExists.WithPayload("node ip was exist"),
			},
		},
		{
			Input: struct {
				Cluster Cluster
				Node    *Node
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				Node: &Node{
					Name: "node1",
					ConnectionData: ConnectionData{
						IP: "192.168.31.30",
					},
				},
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				ReturnValue: h.EExists.WithPayload("node name was exist"),
			},
		},
	}

	for _, item := range tests {

		assert.Equal(t, item.Want.ReturnValue, item.Input.Cluster.AddNode(item.Input.Node))
		assert.Equal(t, item.Want.Cluster, item.Input.Cluster)
	}
}

func TestCluster_UpdateNode(t *testing.T) {

	tests := []struct {
		Input struct {
			Cluster Cluster
			Node    *Node
		}
		Want struct {
			Cluster     Cluster
			ReturnValue error
		}
	}{
		{
			Input: struct {
				Cluster Cluster
				Node    *Node
			}{
				Cluster: Cluster{
					Nodes: []*Node{},
					lock:  new(sync.RWMutex),
				},
				Node: &Node{
					Name: "node1",
				},
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{},
					lock:  new(sync.RWMutex),
				},
				ReturnValue: h.ENotFound.WithPayload("node ip not exist"),
			},
		},
		{
			Input: struct {
				Cluster Cluster
				Node    *Node
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP:                 "192.168.31.1",
								Port:               22,
								Username:           "root",
								AuthenticationType: AuthenticationTypePassword,
								Password:           "123456",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				Node: &Node{
					Name: "node1",
					ConnectionData: ConnectionData{
						IP:                 "192.168.31.1",
						Port:               22,
						Username:           "root",
						AuthenticationType: AuthenticationTypePassword,
						Password:           "45678",
					},
				},
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
							ConnectionData: ConnectionData{
								IP:                 "192.168.31.1",
								Port:               22,
								Username:           "root",
								AuthenticationType: AuthenticationTypePassword,
								Password:           "45678",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				ReturnValue: nil,
			},
		},
		{
			Input: struct {
				Cluster Cluster
				Node    *Node
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
						{
							Name: "node3",
							ConnectionData: ConnectionData{
								IP: "192.168.31.2",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				Node: &Node{
					Name: "node3",
					ConnectionData: ConnectionData{
						IP: "192.168.31.1",
					},
				},
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
						{
							Name: "node3",
							ConnectionData: ConnectionData{
								IP: "192.168.31.2",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				ReturnValue: h.EExists.WithPayload("node name was exist"),
			},
		},
	}

	for _, item := range tests {

		assert.Equal(t, item.Want.ReturnValue, item.Input.Cluster.UpdateNode(item.Input.Node))
		assert.Equal(t, item.Want.Cluster, item.Input.Cluster)
	}
}

func TestCluster_DeleteNode(t *testing.T) {

	tests := []struct {
		Input struct {
			Cluster Cluster
			IP      string
		}
		Want struct {
			Cluster     Cluster
			ReturnValue error
		}
	}{
		{
			Input: struct {
				Cluster Cluster
				IP      string
			}{
				Cluster: Cluster{
					Nodes: []*Node{},
					lock:  new(sync.RWMutex),
				},
				IP: "192.168.31.1",
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{},
					lock:  new(sync.RWMutex),
				},
				ReturnValue: h.EExists.WithPayload("node not exist"),
			},
		},
		{
			Input: struct {
				Cluster Cluster
				IP      string
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP: "192.168.31.2",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				IP: "192.168.31.1",
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP: "192.168.31.2",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				ReturnValue: h.EExists.WithPayload("node not exist"),
			},
		},
		{
			Input: struct {
				Cluster Cluster
				IP      string
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				IP: "192.168.31.1",
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{},
					lock:  new(sync.RWMutex),
				},
				ReturnValue: nil,
			},
		},
		{
			Input: struct {
				Cluster Cluster
				IP      string
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP: "192.168.31.2",
							},
						},
						{
							Name: "node3",
							ConnectionData: ConnectionData{
								IP: "192.168.31.3",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				IP: "192.168.31.2",
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
						{
							Name: "node3",
							ConnectionData: ConnectionData{
								IP: "192.168.31.3",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				ReturnValue: nil,
			},
		},
	}

	for _, item := range tests {

		assert.Equal(t, item.Want.ReturnValue, item.Input.Cluster.DeleteNode(item.Input.IP))
		assert.Equal(t, item.Want.Cluster, item.Input.Cluster)
	}
}

func TestCluster_GetNode(t *testing.T) {

	tests := []struct {
		Input struct {
			Cluster Cluster
			IP      string
		}
		Want *Node
	}{
		{
			Input: struct {
				Cluster Cluster
				IP      string
			}{
				Cluster: Cluster{
					Nodes: []*Node{},
					lock:  new(sync.RWMutex),
				},
				IP: "192.168.31.1",
			},
			Want: nil,
		},
		{
			Input: struct {
				Cluster Cluster
				IP      string
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				IP: "192.168.31.1",
			},
			Want: &Node{
				Name: "node2",
				ConnectionData: ConnectionData{
					IP: "192.168.31.1",
				},
			},
		},
		{
			Input: struct {
				Cluster Cluster
				IP      string
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				IP: "192.168.31.3",
			},
			Want: nil,
		},
	}

	for _, item := range tests {

		assert.Equal(t, item.Want, item.Input.Cluster.GetNode(item.Input.IP))
	}
}

func TestCluster_GetNodeByName(t *testing.T) {

	tests := []struct {
		Input struct {
			Cluster Cluster
			Name    string
		}
		Want *Node
	}{
		{
			Input: struct {
				Cluster Cluster
				Name    string
			}{
				Cluster: Cluster{
					Nodes: []*Node{},
					lock:  new(sync.RWMutex),
				},
				Name: "node1",
			},
			Want: nil,
		},
		{
			Input: struct {
				Cluster Cluster
				Name    string
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				Name: "node2",
			},
			Want: &Node{
				Name: "node2",
				ConnectionData: ConnectionData{
					IP: "192.168.31.1",
				},
			},
		},
		{
			Input: struct {
				Cluster Cluster
				Name    string
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					lock: new(sync.RWMutex),
				},
				Name: "node1",
			},
			Want: nil,
		},
	}

	for _, item := range tests {

		assert.Equal(t, item.Want, item.Input.Cluster.GetNodeByName(item.Input.Name))
	}
}

func TestCluster_MarkNodeChecking(t *testing.T) {

	tests := []struct {
		Input Cluster
		Want  struct {
			Cluster     Cluster
			ReturnValue error
		}
	}{
		{
			Input: Cluster{
				Nodes:              []*Node{},
				ClusterCheckResult: constant.CheckResultNotRunning,
				lock:               new(sync.RWMutex),
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes:              []*Node{},
					ClusterCheckResult: constant.CheckResultNotRunning,
					lock:               new(sync.RWMutex),
				},
				ReturnValue: nil,
			},
		},
		{
			Input: Cluster{
				Nodes: []*Node{
					{
						Name: "node2",
						ConnectionData: ConnectionData{
							IP: "192.168.31.1",
						},
					},
				},
				ClusterCheckResult: constant.CheckResultNotRunning,
				lock:               new(sync.RWMutex),
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					ClusterCheckResult: constant.CheckResultChecking,
					lock:               new(sync.RWMutex),
				},
				ReturnValue: nil,
			},
		},
		{
			Input: Cluster{
				Nodes: []*Node{
					{
						Name: "node1",
						ConnectionData: ConnectionData{
							IP: "192.168.31.1",
						},
					},
				},
				ClusterCheckResult: constant.CheckResultChecking,
				lock:               new(sync.RWMutex),
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node1",
							ConnectionData: ConnectionData{
								IP: "192.168.31.1",
							},
						},
					},
					ClusterCheckResult: constant.CheckResultChecking,
					lock:               new(sync.RWMutex),
				},
				ReturnValue: errors.New("was checking"),
			},
		},
	}

	for _, item := range tests {

		assert.Equal(t, item.Want.ReturnValue, item.Input.MarkNodeChecking())
		assert.Equal(t, item.Want.Cluster, item.Input)
	}
}

func TestCluster_ClearClusterCheckingData(t *testing.T) {

	tests := []struct {
		Input Cluster
		Want  Cluster
	}{
		{
			Input: Cluster{
				Nodes:              []*Node{},
				ClusterCheckResult: constant.CheckResultNotRunning,
				lock:               new(sync.RWMutex),
			},
			Want: Cluster{
				Nodes:              []*Node{},
				ClusterCheckResult: constant.CheckResultNotRunning,
				lock:               new(sync.RWMutex),
			},
		},
		{
			Input: Cluster{
				Nodes: []*Node{
					{
						Name:        "node2",
						CheckReport: &CheckReport{},
					},
				},
				ClusterCheckResult: constant.CheckResultChecking,
				lock:               new(sync.RWMutex),
			},
			Want: Cluster{
				Nodes: []*Node{
					{
						Name: "node2",
						CheckReport: &CheckReport{
							CheckItems:  make([]*CheckItem, 0, 0),
							CheckResult: constant.CheckResultNotRunning,
						},
					},
				},
				ClusterCheckResult: constant.CheckResultNotRunning,
				lock:               new(sync.RWMutex),
			},
		},
	}

	for _, item := range tests {

		item.Input.ClearClusterCheckingData()
		assert.Equal(t, item.Want, item.Input)
	}
}

func TestCluster_SetClusterCheckResult(t *testing.T) {

	tests := []struct {
		Input struct {
			Cluster       Cluster
			CheckResult   constant.CheckResult
			FailureDetail *common.FailureDetail
		}
		Want Cluster
	}{
		{
			Input: struct {
				Cluster       Cluster
				CheckResult   constant.CheckResult
				FailureDetail *common.FailureDetail
			}{
				Cluster: Cluster{
					ClusterCheckResult: constant.CheckResultChecking,
					lock:               new(sync.RWMutex),
				},
				CheckResult:   constant.CheckResultPassed,
				FailureDetail: nil,
			},
			Want: Cluster{
				ClusterCheckResult: constant.CheckResultPassed,
				lock:               new(sync.RWMutex),
			},
		},
		{
			Input: struct {
				Cluster       Cluster
				CheckResult   constant.CheckResult
				FailureDetail *common.FailureDetail
			}{
				Cluster: Cluster{
					ClusterCheckResult: constant.CheckResultChecking,
					lock:               new(sync.RWMutex),
				},
				CheckResult: constant.CheckResultFailed,
				FailureDetail: &common.FailureDetail{
					Reason:     "reason",
					Detail:     "detail",
					FixMethods: "fix",
					LogId:      1,
				},
			},
			Want: Cluster{
				ClusterCheckResult: constant.CheckResultFailed,
				ClusterCheckError: &common.FailureDetail{
					Reason:     "reason",
					Detail:     "detail",
					FixMethods: "fix",
					LogId:      1,
				},
				lock: new(sync.RWMutex),
			},
		},
	}

	for _, item := range tests {

		item.Input.Cluster.SetClusterCheckResult(item.Input.CheckResult, item.Input.FailureDetail)
		assert.Equal(t, item.Want, item.Input.Cluster)
	}
}

func TestCluster_ClearClusterDeployData(t *testing.T) {

	tests := []struct {
		Input Cluster
		Want  Cluster
	}{
		{
			Input: Cluster{
				Nodes:               []*Node{},
				DeployClusterStatus: DeployClusterStatusSuccessful,
				lock:                new(sync.RWMutex),
			},
			Want: Cluster{
				Nodes:               []*Node{},
				DeployClusterStatus: DeployClusterStatusSuccessful,
				lock:                new(sync.RWMutex),
			},
		},
		{
			Input: Cluster{
				Nodes: []*Node{
					{
						Name: "node2",
						DeploymentReports: map[constant.MachineRole]*DeploymentReport{
							constant.MachineRoleMaster: {
								Role: constant.MachineRoleMaster,
							},
						},
					},
				},
				DeployClusterStatus: DeployClusterStatusSuccessful,
				lock:                new(sync.RWMutex),
			},
			Want: Cluster{
				Nodes: []*Node{
					{
						Name:              "node2",
						DeploymentReports: map[constant.MachineRole]*DeploymentReport{},
					},
				},
				DeployClusterStatus: DeployClusterStatusNotRunning,
				lock:                new(sync.RWMutex),
			},
		},
		{
			Input: Cluster{
				Nodes: []*Node{
					{
						Name: "node2",
						DeploymentReports: map[constant.MachineRole]*DeploymentReport{
							constant.MachineRoleMaster: {
								Role: constant.MachineRoleMaster,
							},
						},
					},
				},
				DeployClusterStatus: DeployClusterStatusFailed,
				DeployClusterError:  common.NewFailureDetail(),
				lock:                new(sync.RWMutex),
			},
			Want: Cluster{
				Nodes: []*Node{
					{
						Name:              "node2",
						DeploymentReports: map[constant.MachineRole]*DeploymentReport{},
					},
				},
				DeployClusterStatus: DeployClusterStatusNotRunning,
				lock:                new(sync.RWMutex),
			},
		},
	}

	for _, item := range tests {

		item.Input.ClearClusterDeployData()
		assert.Equal(t, item.Want, item.Input)
	}
}

func TestCluster_MarkNodeDeploying(t *testing.T) {

	tests := []struct {
		Input Cluster
		Want  struct {
			Cluster     Cluster
			ReturnValue error
		}
	}{
		{
			Input: Cluster{
				Nodes:               []*Node{},
				DeployClusterStatus: DeployClusterStatusSuccessful,
				lock:                new(sync.RWMutex),
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes:               []*Node{},
					DeployClusterStatus: DeployClusterStatusSuccessful,
					lock:                new(sync.RWMutex),
				},
				ReturnValue: nil,
			},
		},
		{
			Input: Cluster{
				Nodes: []*Node{
					{
						Name: "node2",
					},
				},
				DeployClusterStatus: DeployClusterStatusNotRunning,
				lock:                new(sync.RWMutex),
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
						},
					},
					DeployClusterStatus: DeployClusterStatusRunning,
					lock:                new(sync.RWMutex),
				},
				ReturnValue: nil},
		},
		{
			Input: Cluster{
				Nodes: []*Node{
					{
						Name: "node2",
					},
				},
				DeployClusterStatus: DeployClusterStatusRunning,
				lock:                new(sync.RWMutex),
			},
			Want: struct {
				Cluster     Cluster
				ReturnValue error
			}{
				Cluster: Cluster{
					Nodes: []*Node{
						{
							Name: "node2",
						},
					},
					DeployClusterStatus: DeployClusterStatusRunning,
					lock:                new(sync.RWMutex),
				},
				ReturnValue: errors.New("was running"),
			},
		},
	}

	for _, item := range tests {

		assert.Equal(t, item.Want.ReturnValue, item.Input.MarkNodeDeploying())
		assert.Equal(t, item.Want.Cluster, item.Input)
	}
}

func TestCluster_SetClusterDeploymentStatus(t *testing.T) {

	tests := []struct {
		Input struct {
			Cluster             Cluster
			DeployClusterStatus DeployClusterStatus
			FailureDetail       *common.FailureDetail
		}
		Want Cluster
	}{
		{
			Input: struct {
				Cluster             Cluster
				DeployClusterStatus DeployClusterStatus
				FailureDetail       *common.FailureDetail
			}{
				Cluster: Cluster{
					DeployClusterStatus: DeployClusterStatusRunning,
					lock:                new(sync.RWMutex),
				},
				DeployClusterStatus: DeployClusterStatusSuccessful,
				FailureDetail:       nil,
			},
			Want: Cluster{
				DeployClusterStatus: DeployClusterStatusSuccessful,
				lock:                new(sync.RWMutex),
			},
		},
		{
			Input: struct {
				Cluster             Cluster
				DeployClusterStatus DeployClusterStatus
				FailureDetail       *common.FailureDetail
			}{
				Cluster: Cluster{
					DeployClusterStatus: DeployClusterStatusRunning,
					lock:                new(sync.RWMutex),
				},
				DeployClusterStatus: DeployClusterStatusFailed,
				FailureDetail: &common.FailureDetail{
					Reason:     "reason",
					Detail:     "detail",
					FixMethods: "fix",
					LogId:      1,
				},
			},
			Want: Cluster{
				DeployClusterStatus: DeployClusterStatusFailed,
				DeployClusterError: &common.FailureDetail{
					Reason:     "reason",
					Detail:     "detail",
					FixMethods: "fix",
					LogId:      1,
				},
				lock: new(sync.RWMutex),
			},
		},
	}

	for _, item := range tests {

		item.Input.Cluster.SetClusterDeploymentStatus(item.Input.DeployClusterStatus, item.Input.FailureDetail)
		assert.Equal(t, item.Want, item.Input.Cluster)
	}
}

func TestCluster_AddNodeList(t *testing.T) {

	tests := []struct {
		BaseNodeList []*Node
		Input        []*Node
		Want         error
		WantNodeList []*Node
	}{
		{
			BaseNodeList: []*Node{},
			Input: []*Node{
				{
					ConnectionData:      ConnectionData{},
					Name:                "k8s-master1",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
			Want: nil,
			WantNodeList: []*Node{
				{
					ConnectionData:      ConnectionData{},
					Name:                "k8s-master1",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
		},
		{
			BaseNodeList: []*Node{
				{
					ConnectionData: ConnectionData{
						IP:   "192.168.31.1",
						Port: 22,
					},
					Name:                "k8s-master1",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
			Input: []*Node{
				{
					ConnectionData: ConnectionData{
						IP:   "192.168.31.2",
						Port: 22,
					},
					Name:                "k8s-master2",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
			Want: nil,
			WantNodeList: []*Node{
				{
					ConnectionData: ConnectionData{
						IP:   "192.168.31.1",
						Port: 22,
					},
					Name:                "k8s-master1",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
				{
					ConnectionData: ConnectionData{
						IP:   "192.168.31.2",
						Port: 22,
					},
					Name:                "k8s-master2",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
		},
		{
			BaseNodeList: []*Node{
				{
					ConnectionData: ConnectionData{
						IP:   "192.168.31.1",
						Port: 22,
					},
					Name:                "k8s-master1",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
			Input: []*Node{
				{
					ConnectionData: ConnectionData{
						IP:   "192.168.31.1",
						Port: 22,
					},
					Name:                "k8s-master2",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
			Want: h.EExists.WithPayload(fmt.Sprintf("node ip 192.168.31.1 was exist")),
			WantNodeList: []*Node{
				{
					ConnectionData: ConnectionData{
						IP:   "192.168.31.1",
						Port: 22,
					},
					Name:                "k8s-master1",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
		},
		{
			BaseNodeList: []*Node{
				{
					ConnectionData: ConnectionData{
						IP:   "192.168.31.1",
						Port: 22,
					},
					Name:                "k8s-master1",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
			Input: []*Node{
				{
					ConnectionData: ConnectionData{
						IP:   "192.168.31.2",
						Port: 22,
					},
					Name:                "k8s-master1",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
			Want: h.EExists.WithPayload(fmt.Sprintf("node name k8s-master1 was exist")),
			WantNodeList: []*Node{
				{
					ConnectionData: ConnectionData{
						IP:   "192.168.31.1",
						Port: 22,
					},
					Name:                "k8s-master1",
					Description:         "k8s-description",
					MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleWorker},
					Labels:              []*Label{},
					Taints:              []*Taint{},
					CheckReport:         &CheckReport{},
					DeploymentReports:   make(map[constant.MachineRole]*DeploymentReport),
					DockerRootDirectory: "/var/lib/docker",
					rwLock:              sync.RWMutex{},
				},
			},
		},
	}

	for _, test := range tests {
		cluster := NewCluster()
		cluster.Nodes = test.BaseNodeList
		assert.Equal(t, test.Want, cluster.AddNodeList(test.Input))
		assert.Equal(t, test.WantNodeList, cluster.Nodes)
	}
}
