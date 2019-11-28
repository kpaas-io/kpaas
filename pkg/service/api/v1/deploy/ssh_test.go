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

	"github.com/kpaas-io/kpaas/pkg/service/model/api"
)

func TestTestConnectNode(t *testing.T) {

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	body := api.ConnectionData{
		IP:   "192.168.31.101",
		Port: uint16(22),
		SSHLoginData: api.SSHLoginData{
			Username:           "root",
			AuthenticationType: api.AuthenticationTypePassword,
			Password:           "123456",
		},
	}
	bodyContent, err := json.Marshal(body)
	assert.Nil(t, err)
	bodyReader := bytes.NewReader(bodyContent)

	ctx.Request = httptest.NewRequest("POST", "/api/v1/ssh/tests", bodyReader)

	TestConnectNode(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.SuccessfulOption)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)
	assert.True(t, responseData.Success)
}
