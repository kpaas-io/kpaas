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

package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
)

func TestGetWizardProgress(t *testing.T) {

	wizard.ClearCurrentWizardData()
	responseData := getWizardProgressData(t)

	assert.Equal(t, api.KubeAPIServerConnectTypeFirstMasterIP, responseData.ClusterData.KubeAPIServerConnectType)
	assert.Equal(t, "", responseData.ClusterData.ShortName)
	assert.Equal(t, "", responseData.ClusterData.Name)
	assert.Equal(t, "", responseData.ClusterData.VIP)
	assert.Equal(t, "", responseData.ClusterData.NetInterfaceName)
	assert.Equal(t, "", responseData.ClusterData.LoadbalancerIP)
	assert.Equal(t, uint16(0), responseData.ClusterData.LoadbalancerPort)
	assert.Equal(t, wizard.DefaultNodePortMinimum, responseData.ClusterData.NodePortMinimum)
	assert.Equal(t, wizard.DefaultNodePortMaximum, responseData.ClusterData.NodePortMaximum)
	assert.Empty(t, responseData.ClusterData.Labels)
	assert.Empty(t, responseData.ClusterData.Annotations)
	assert.Empty(t, responseData.NodesData)
	assert.Empty(t, responseData.CheckingData)
	assert.Empty(t, responseData.DeploymentData)
	assert.Equal(t, constant.CheckResultPending, responseData.CheckResult)
	assert.Equal(t, api.DeployClusterStatusPending, responseData.DeployClusterStatus)
}

func TestGetWizardProgress2(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	wizardData.Info.ShortName = "test-cluster"
	wizardData.Info.Name = "ClusterName"
	wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType = wizard.KubeAPIServerConnectTypeKeepalived
	wizardData.Info.KubeAPIServerConnection.VIP = "192.168.31.100"
	wizardData.Info.KubeAPIServerConnection.NetInterfaceName = "em0"
	wizardData.Info.NodePortMinimum = 20000
	wizardData.Info.NodePortMaximum = 29999
	wizardData.Info.Labels = []*wizard.Label{
		{
			Key:   "for-test",
			Value: "yes",
		},
	}
	wizardData.Info.Annotations = []*wizard.Annotation{
		{
			Key:   "comment",
			Value: "Icanspeakenglish",
		},
	}

	responseData := getWizardProgressData(t)

	assert.Equal(t, api.KubeAPIServerConnectTypeKeepalived, responseData.ClusterData.KubeAPIServerConnectType)
	assert.Equal(t, "test-cluster", responseData.ClusterData.ShortName)
	assert.Equal(t, "ClusterName", responseData.ClusterData.Name)
	assert.Equal(t, "192.168.31.100", responseData.ClusterData.VIP)
	assert.Equal(t, "em0", responseData.ClusterData.NetInterfaceName)
	assert.Equal(t, "", responseData.ClusterData.LoadbalancerIP)
	assert.Equal(t, uint16(0), responseData.ClusterData.LoadbalancerPort)
	assert.Equal(t, uint16(20000), responseData.ClusterData.NodePortMinimum)
	assert.Equal(t, uint16(29999), responseData.ClusterData.NodePortMaximum)
	assert.Equal(t, []api.Label{{Key: "for-test", Value: "yes"}}, responseData.ClusterData.Labels)
	assert.Equal(t, []api.Annotation{{Key: "comment", Value: "Icanspeakenglish"}}, responseData.ClusterData.Annotations)
}

func TestGetWizardProgress3(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	wizardData.Nodes = []*wizard.Node{
		{
			Name:         "master1",
			Description:  "desc1",
			MachineRoles: []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
			Labels: []*wizard.Label{
				{
					Key:   "kpaas.io/test",
					Value: "yes",
				},
			},
			Taints: []*wizard.Taint{
				{
					Key:    "taint1",
					Value:  "taint-value",
					Effect: wizard.TaintEffectNoExecute,
				},
			},
			DockerRootDirectory: "/mnt/docker",
			ConnectionData: wizard.ConnectionData{
				IP:                 "192.168.31.101",
				Port:               22,
				Username:           "kpaas",
				AuthenticationType: wizard.AuthenticationTypePassword,
				Password:           "123456",
			},
			CheckReport: &wizard.CheckReport{
				CheckItems:  make([]*wizard.CheckItem, 0),
				CheckResult: constant.CheckResultPending,
			},
		},
	}

	responseData := getWizardProgressData(t)
	assert.Len(t, responseData.NodesData, 1)
	node := responseData.NodesData[0]
	assert.Equal(t, "master1", node.Name)
	assert.Equal(t, "desc1", node.Description)
	assert.Equal(t, []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd}, node.MachineRoles)
	assert.Equal(t, []api.Label{{Key: "kpaas.io/test", Value: "yes"}}, node.Labels)
	assert.Equal(t, []api.Taint{{Key: "taint1", Value: "taint-value", Effect: api.TaintEffectNoExecute}}, node.Taints)
	assert.Equal(t, "/mnt/docker", node.DockerRootDirectory)
	assert.Equal(t, "192.168.31.101", node.IP)
	assert.Equal(t, uint16(22), node.Port)
	assert.Equal(t, "kpaas", node.Username)
	assert.Equal(t, api.AuthenticationTypePassword, node.AuthenticationType)
	assert.Equal(t, "", node.Password)
	assert.Equal(t, "", node.PrivateKeyName)
}

func TestGetWizardProgress4(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	wizardData.ClusterCheckResult = constant.CheckResultSuccessful
	wizardData.Nodes = []*wizard.Node{
		{
			Name: "master1",
			CheckReport: &wizard.CheckReport{
				CheckItems: []*wizard.CheckItem{
					{
						ItemName:    "check 1",
						CheckResult: constant.CheckResultSuccessful,
					},
					{
						ItemName:    "check 2",
						CheckResult: constant.CheckResultSuccessful,
					},
				},
				CheckResult: constant.CheckResultSuccessful,
			},
		},
	}

	responseData := getWizardProgressData(t)
	assert.Len(t, responseData.CheckingData, 1)
	checkData := responseData.CheckingData[0]
	assert.Equal(t, "master1", checkData.Name)
	assert.Equal(t, []api.CheckingItem{
		{
			CheckingPoint: "check 1",
			Result:        constant.CheckResultSuccessful,
		},
		{
			CheckingPoint: "check 2",
			Result:        constant.CheckResultSuccessful,
		},
	}, checkData.Items)
	assert.Equal(t, constant.CheckResultSuccessful, responseData.CheckResult)
}

func TestGetWizardProgress5(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	wizardData.Nodes = []*wizard.Node{
		{
			Name: "master1",
			DeploymentReports: map[constant.MachineRole]*wizard.DeploymentReport{
				constant.MachineRoleMaster: {
					Role:   constant.MachineRoleMaster,
					Status: wizard.DeployStatusSuccessful,
				},
				constant.MachineRoleEtcd: {
					Role:   constant.MachineRoleEtcd,
					Status: wizard.DeployStatusSuccessful,
				},
			},
			CheckReport: &wizard.CheckReport{
				CheckItems:  make([]*wizard.CheckItem, 0),
				CheckResult: constant.CheckResultSuccessful,
			},
		},
	}

	responseData := getWizardProgressData(t)
	assert.Len(t, responseData.DeploymentData, 2)
	sortRoles(responseData.DeploymentData)
	assert.Equal(t, []api.DeploymentResponseData{
		{
			Role: constant.MachineRoleMaster,
			Nodes: []api.DeploymentNode{
				{
					Name:   "master1",
					Status: api.DeployStatusSuccessful,
				},
			},
		},
		{
			Role: constant.MachineRoleEtcd,
			Nodes: []api.DeploymentNode{
				{
					Name:   "master1",
					Status: api.DeployStatusSuccessful,
				},
			},
		},
	}, responseData.DeploymentData)
}

func TestGetWizardProgress6(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	wizardData.DeployClusterStatus = wizard.DeployClusterStatusSuccessful
	responseData := getWizardProgressData(t)
	assert.Equal(t, api.DeployClusterStatusSuccessful, responseData.DeployClusterStatus)
}

func TestGetWizardProgress7(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	wizardData.Info.ShortName = "test-cluster"
	wizardData.Info.Name = "ClusterName"
	wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType = wizard.KubeAPIServerConnectTypeLoadBalancer
	wizardData.Info.KubeAPIServerConnection.LoadbalancerIP = "192.168.31.200"
	wizardData.Info.KubeAPIServerConnection.LoadbalancerPort = uint16(3434)

	responseData := getWizardProgressData(t)

	assert.Equal(t, api.KubeAPIServerConnectTypeLoadBalancer, responseData.ClusterData.KubeAPIServerConnectType)
	assert.Equal(t, "test-cluster", responseData.ClusterData.ShortName)
	assert.Equal(t, "ClusterName", responseData.ClusterData.Name)
	assert.Equal(t, "", responseData.ClusterData.VIP)
	assert.Equal(t, "", responseData.ClusterData.NetInterfaceName)
	assert.Equal(t, "192.168.31.200", responseData.ClusterData.LoadbalancerIP)
	assert.Equal(t, uint16(3434), responseData.ClusterData.LoadbalancerPort)
}

func getWizardProgressData(t *testing.T) (responseData *api.GetWizardResponse) {

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", "/api/v1/deploy/wizard/progresses", bytes.NewReader([]byte{}))
	GetWizardProgress(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData = new(api.GetWizardResponse)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)
	return responseData
}
