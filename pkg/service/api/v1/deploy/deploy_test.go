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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

func TestDeploy(t *testing.T) {

	wizard.ClearCurrentWizardData()
	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/deploys", nil)

	Deploy(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(h.AppErr)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestDeploy2(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	node := wizard.NewNode()
	node.Name = "master1"
	node.CheckReport = &wizard.CheckReport{
		CheckResult: constant.CheckResultFailed,
	}
	wizardData.Nodes = []*wizard.Node{
		node,
	}

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/deploys", nil)

	Deploy(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(h.AppErr)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, h.EStatusError.Msg, responseData.Msg)
}

func TestDeploy3(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	wizardData.ClusterCheckResult = constant.CheckResultSuccessful
	node := wizard.NewNode()
	node.Name = "master1"
	node.CheckReport = &wizard.CheckReport{
		CheckResult: constant.CheckResultSuccessful,
	}
	wizardData.Nodes = []*wizard.Node{
		node,
	}

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/deploys", nil)

	Deploy(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.SuccessfulOption)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.True(t, responseData.Success)
}

func TestGetDeployReport(t *testing.T) {

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
		},
	}
	wizardData.DeployClusterStatus = wizard.DeployClusterStatusSuccessful

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", "/api/v1/deploy/wizard/deploys", nil)

	GetDeployReport(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.GetDeploymentReportResponse)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, api.DeployClusterStatusSuccessful, responseData.DeployClusterStatus)
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
	}, sortRoles(responseData.Roles))
	assert.Nil(t, responseData.DeployClusterError)
}

func TestFetchKubeConfigContent(t *testing.T) {

	tests := []struct {
		OriginNodeList []*wizard.Node
		WantKubeConfig string
	}{
		{
			OriginNodeList: []*wizard.Node{
				{
					ConnectionData: wizard.ConnectionData{
						IP:                 "192.168.1.1",
						Port:               22,
						Username:           "root",
						AuthenticationType: wizard.AuthenticationTypePassword,
						Password:           "123456",
					},
					Name:         "k8s-master1",
					MachineRoles: []constant.MachineRole{constant.MachineRoleMaster},
				},
			},
			WantKubeConfig: "kube config content",
		},
	}

	for _, test := range tests {

		wizard.ClearCurrentWizardData()
		wizardData := wizard.GetCurrentWizard()
		wizardData.Nodes = test.OriginNodeList
		fetchKubeConfigContent()
		assert.Equal(t, test.WantKubeConfig, *wizardData.KubeConfig)
	}
}

func TestComputeClusterDeployStatus(t *testing.T) {

	tests := []struct {
		Input *protos.GetDeployResultReply
		Want  wizard.DeployClusterStatus
	}{
		{
			Input: &protos.GetDeployResultReply{
				Status: string(constant.OperationStatusSuccessful),
				Err:    nil,
				Items:  nil,
			},
			Want: wizard.DeployClusterStatusSuccessful,
		},
		{
			Input: &protos.GetDeployResultReply{
				Status: string(constant.OperationStatusPending),
				Err:    nil,
				Items:  nil,
			},
			Want: wizard.DeployClusterStatusPending,
		},
		{
			Input: &protos.GetDeployResultReply{
				Status: string(constant.OperationStatusRunning),
				Err:    nil,
				Items:  nil,
			},
			Want: wizard.DeployClusterStatusRunning,
		},
		{
			Input: &protos.GetDeployResultReply{
				Status: string(constant.OperationStatusUnknown),
				Err:    nil,
				Items:  nil,
			},
			Want: wizard.DeployClusterStatusDeployServiceUnknown,
		},
		{
			Input: &protos.GetDeployResultReply{
				Status: string(constant.OperationStatusFailed),
				Err:    nil,
				Items:  nil,
			},
			Want: wizard.DeployClusterStatusFailed,
		},
		{
			Input: &protos.GetDeployResultReply{
				Status: string(constant.OperationStatusFailed),
				Err:    nil,
				Items: []*protos.DeployItemResult{
					{
						DeployItem: &protos.DeployItem{
							Role: string(constant.MachineRoleEtcd),
						},
						Status: string(constant.OperationStatusFailed),
					},
				},
			},
			Want: wizard.DeployClusterStatusFailed,
		},
		{
			Input: &protos.GetDeployResultReply{
				Status: string(constant.OperationStatusFailed),
				Err:    nil,
				Items: []*protos.DeployItemResult{
					{
						DeployItem: nil,
						Status:     string(constant.OperationStatusFailed),
					},
				},
			},
			Want: wizard.DeployClusterStatusFailed,
		},
		{
			Input: &protos.GetDeployResultReply{
				Status: string(constant.OperationStatusFailed),
				Err:    nil,
				Items: []*protos.DeployItemResult{
					{
						DeployItem: &protos.DeployItem{
							Role: string(constant.MachineRoleMaster),
						},
						Status: string(constant.OperationStatusFailed),
					},
				},
			},
			Want: wizard.DeployClusterStatusFailed,
		},
		{
			Input: &protos.GetDeployResultReply{
				Status: string(constant.OperationStatusFailed),
				Err:    nil,
				Items: []*protos.DeployItemResult{
					{
						DeployItem: &protos.DeployItem{
							Role: string(constant.MachineRoleMaster),
						},
						Status: string(constant.OperationStatusSuccessful),
					},
					{
						DeployItem: &protos.DeployItem{
							Role: string(constant.MachineRoleMaster),
						},
						Status: string(constant.OperationStatusFailed),
					},
				},
			},
			Want: wizard.DeployClusterStatusWorkedButHaveError,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.Want, computeClusterDeployStatus(test.Input))
	}
}

func sortRoles(roles []api.DeploymentResponseData) []api.DeploymentResponseData {

	sort.SliceStable(roles, func(i, j int) bool {
		if roles[i].Role == roles[j].Role {
			return false
		}
		if roles[i].Role == constant.MachineRoleMaster {
			return true
		}
		if roles[j].Role == constant.MachineRoleMaster {
			return false
		}
		if roles[i].Role == constant.MachineRoleWorker {
			return true
		}
		if roles[j].Role == constant.MachineRoleWorker {
			return false
		}
		return false
	})
	return roles
}
