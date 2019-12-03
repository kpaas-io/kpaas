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
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

func TestDownloadKubeConfig(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	wizardData.DeployClusterStatus = wizard.DeployClusterStatusSuccessful
	kubeConfig := "123456"
	wizardData.KubeConfig = &kubeConfig

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", "/api/v1/deploy/wizard/kubeconfigs", nil)

	DownloadKubeConfig(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(string)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, "123456", *responseData)
}

func TestDownloadKubeConfig2(t *testing.T) {

	wizard.ClearCurrentWizardData()
	wizardData := wizard.GetCurrentWizard()
	wizardData.DeployClusterStatus = wizard.DeployClusterStatusSuccessful

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", "/api/v1/deploy/wizard/kubeconfigs", nil)

	DownloadKubeConfig(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(h.AppErr)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, h.ENotFound.Status, resp.Code)
}

func TestDownloadKubeConfig3(t *testing.T) {

	wizard.ClearCurrentWizardData()

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", "/api/v1/deploy/wizard/kubeconfigs", nil)

	DownloadKubeConfig(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(h.AppErr)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, h.EStatusError.Status, resp.Code)
}
