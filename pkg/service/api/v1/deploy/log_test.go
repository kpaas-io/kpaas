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
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

func TestDownloadLog(t *testing.T) {

	var err error
	content := "abcdefghijklmn"
	logId, err := wizard.SetLogByString(content)
	assert.Nil(t, err)
	assert.Greater(t, logId, uint64(0))

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/deploy/wizard/logs/%d", logId), nil)
	ctx.Params = gin.Params{
		{
			Key:   "id",
			Value: fmt.Sprintf("%d", logId),
		},
	}

	DownloadLog(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	var bodyContent = new(string)
	err = json.Unmarshal(resp.Body.Bytes(), bodyContent)
	assert.Nil(t, err)
	fmt.Printf("result: %s\n", resp.Body.String())

	assert.Equal(t, content, *bodyContent)
}

func TestDownloadLog_NotFound(t *testing.T) {

	var err error
	rand.Seed(time.Now().UnixNano())

	logId := rand.Uint64()
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/deploy/wizard/logs/%d", logId), nil)
	ctx.Params = gin.Params{
		{
			Key:   "id",
			Value: fmt.Sprintf("%d", logId),
		},
	}

	DownloadLog(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	responseData := new(h.AppErr)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)
	fmt.Printf("result: %s\n", resp.Body.String())

	assert.Equal(t, h.ENotFound.Status, resp.Code)
}

func TestDownloadLog_IDInvalid(t *testing.T) {

	var err error
	logId := "invalid"
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/deploy/wizard/logs/%s", logId), nil)
	ctx.Params = gin.Params{
		{
			Key:   "id",
			Value: logId,
		},
	}

	DownloadLog(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	responseData := new(h.AppErr)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)
	fmt.Printf("result: %s\n", resp.Body.String())

	assert.Equal(t, h.EParamsError.Status, resp.Code)
}

func TestDownloadLog_IDEcmpty(t *testing.T) {

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", "/api/v1/deploy/wizard/logs/", nil)

	DownloadLog(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	responseData := new(h.AppErr)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)
	fmt.Printf("result: %s\n", resp.Body.String())

	assert.Equal(t, h.ENotFound.Status, resp.Code)
}
