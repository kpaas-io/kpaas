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

	"github.com/kpaas-io/kpaas/pkg/constant"
	grpcClient "github.com/kpaas-io/kpaas/pkg/service/grpcutils/client"
	"github.com/kpaas-io/kpaas/pkg/service/grpcutils/mock"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

func TestCheckNodeList(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	mockNode := wizard.NewNode()
	mockNode.Name = "master1"
	wizardData.Nodes = []*wizard.Node{mockNode}

	grpcClient.SetDeployController(mock.NewDeployController())

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/checks", nil)

	CheckNodeList(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.SuccessfulOption)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.True(t, responseData.Success)
}

func TestCheckNodeList2(t *testing.T) {

	wizard.ClearCurrentWizardData()

	grpcClient.SetDeployController(mock.NewDeployController())

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/checks", nil)

	CheckNodeList(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(h.AppErr)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestGetCheckingNodeListResult(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	wizardData.ClusterCheckResult = constant.CheckResultPassed
	wizardData.Nodes = []*wizard.Node{
		{
			Name: "master1",
			CheckReport: &wizard.CheckReport{
				CheckItems: []*wizard.CheckItem{
					{
						ItemName:    "check 1",
						CheckResult: constant.CheckResultPassed,
					},
					{
						ItemName:    "check 2",
						CheckResult: constant.CheckResultPassed,
					},
				},
				CheckResult: constant.CheckResultPassed,
			},
		},
	}

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", "/api/v1/deploy/wizard/checks", nil)

	GetCheckingNodeListResult(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.GetCheckingResultResponse)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Len(t, responseData.Nodes, 1)
	checkData := responseData.Nodes[0]
	assert.Equal(t, "master1", checkData.Name)
	assert.Equal(t, []api.CheckingItem{
		{
			CheckingPoint: "check 1",
			Result:        constant.CheckResultPassed,
		},
		{
			CheckingPoint: "check 2",
			Result:        constant.CheckResultPassed,
		},
	}, checkData.Items)

	assert.Equal(t, constant.CheckResultPassed, responseData.Result)
}
