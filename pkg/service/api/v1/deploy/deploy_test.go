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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

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
	wizardData.Nodes = []*wizard.Node{
		{
			Name: "master1",
			CheckItems: []*wizard.CheckItem{
				{
					ItemName:    "check 1",
					CheckResult: wizard.CheckResultFailed,
				},
			},
		},
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
	wizardData.Nodes = []*wizard.Node{
		{
			Name: "master1",
			CheckItems: []*wizard.CheckItem{
				{
					ItemName:    "check 1",
					CheckResult: wizard.CheckResultPassed,
				},
				{
					ItemName:    "check 2",
					CheckResult: wizard.CheckResultPassed,
				},
			},
		},
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
			DeploymentReports: []*wizard.DeploymentReport{
				{
					Role:   wizard.MachineRoleMaster,
					Status: wizard.DeployStatusCompleted,
				},
				{
					Role:   wizard.MachineRoleEtcd,
					Status: wizard.DeployStatusCompleted,
				},
			},
		},
	}
	wizardData.DeploymentStatus = wizard.DeployClusterStatusSuccessful

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
			Role: api.MachineRoleMaster,
			Nodes: []api.DeploymentNode{
				{
					Name:   "master1",
					Status: api.DeployStatusCompleted,
				},
			},
		},
		{
			Role: api.MachineRoleEtcd,
			Nodes: []api.DeploymentNode{
				{
					Name:   "master1",
					Status: api.DeployStatusCompleted,
				},
			},
		},
	}, responseData.Nodes)
}
