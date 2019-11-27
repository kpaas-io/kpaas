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
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
)

func TestSetCluster(t *testing.T) {

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	body := api.Cluster{
		Name:                     "cluster-name",
		ShortName:                "short-name",
		KubeAPIServerConnectType: api.KubeAPIServerConnectTypeFirstMasterIP,
		NodePortMinimum:          uint16(15000),
		NodePortMaximum:          uint16(15999),
		Labels: []api.Label{
			{
				Key:   "label-key",
				Value: "value",
			},
		},
		Annotations: []api.Annotation{
			{
				Key:   "annotation-key",
				Value: "value",
			},
		},
	}
	bodyContent, err := json.Marshal(body)
	assert.Nil(t, err)
	bodyReader := bytes.NewReader(bodyContent)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/clusters", bodyReader)

	SetCluster(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.SuccessfulOption)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.True(t, responseData.Success)

	wizardData := wizard.GetCurrentWizard()
	assert.Equal(t, "cluster-name", wizardData.Info.Name)
	assert.Equal(t, "short-name", wizardData.Info.ShortName)
	assert.Equal(t, wizard.KubeAPIServerConnectTypeFirstMasterIP, wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType)
	assert.Equal(t, "", wizardData.Info.KubeAPIServerConnection.VIP)
	assert.Equal(t, "", wizardData.Info.KubeAPIServerConnection.NetInterfaceName)
	assert.Equal(t, "", wizardData.Info.KubeAPIServerConnection.LoadbalancerIP)
	assert.Equal(t, uint16(0), wizardData.Info.KubeAPIServerConnection.LoadbalancerPort)
	assert.Equal(t, uint16(15000), wizardData.Info.NodePortMinimum)
	assert.Equal(t, uint16(15999), wizardData.Info.NodePortMaximum)
	assert.Equal(t, []*wizard.Label{{Key: "label-key", Value: "value"}}, wizardData.Info.Labels)
	assert.Equal(t, []*wizard.Annotation{{Key: "annotation-key", Value: "value"}}, wizardData.Info.Annotations)
}

func TestSetCluster2(t *testing.T) {

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	body := api.Cluster{
		Name:                     "cluster-name",
		ShortName:                "short-name",
		KubeAPIServerConnectType: api.KubeAPIServerConnectTypeKeepalived,
		VIP:                      "192.168.31.200",
		NetInterfaceName:         "eth0",
		NodePortMinimum:          uint16(15000),
		NodePortMaximum:          uint16(15999),
		Labels: []api.Label{
			{
				Key:   "label-key",
				Value: "value",
			},
		},
		Annotations: []api.Annotation{
			{
				Key:   "annotation-key",
				Value: "value",
			},
		},
	}
	bodyContent, err := json.Marshal(body)
	assert.Nil(t, err)
	bodyReader := bytes.NewReader(bodyContent)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/clusters", bodyReader)

	SetCluster(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.SuccessfulOption)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.True(t, responseData.Success)

	wizardData := wizard.GetCurrentWizard()
	assert.Equal(t, "cluster-name", wizardData.Info.Name)
	assert.Equal(t, "short-name", wizardData.Info.ShortName)
	assert.Equal(t, wizard.KubeAPIServerConnectTypeKeepalived, wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType)
	assert.Equal(t, "192.168.31.200", wizardData.Info.KubeAPIServerConnection.VIP)
	assert.Equal(t, "eth0", wizardData.Info.KubeAPIServerConnection.NetInterfaceName)
	assert.Equal(t, "", wizardData.Info.KubeAPIServerConnection.LoadbalancerIP)
	assert.Equal(t, uint16(0), wizardData.Info.KubeAPIServerConnection.LoadbalancerPort)
	assert.Equal(t, uint16(15000), wizardData.Info.NodePortMinimum)
	assert.Equal(t, uint16(15999), wizardData.Info.NodePortMaximum)
	assert.Equal(t, []*wizard.Label{{Key: "label-key", Value: "value"}}, wizardData.Info.Labels)
	assert.Equal(t, []*wizard.Annotation{{Key: "annotation-key", Value: "value"}}, wizardData.Info.Annotations)
}

func TestSetCluster3(t *testing.T) {

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	body := api.Cluster{
		Name:                     "cluster-name",
		ShortName:                "short-name",
		KubeAPIServerConnectType: api.KubeAPIServerConnectTypeLoadBalancer,
		LoadbalancerIP:           "192.168.31.100",
		LoadbalancerPort:         uint16(4332),
		NodePortMinimum:          uint16(16000),
		NodePortMaximum:          uint16(16999),
	}
	bodyContent, err := json.Marshal(body)
	assert.Nil(t, err)
	bodyReader := bytes.NewReader(bodyContent)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/clusters", bodyReader)

	SetCluster(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.SuccessfulOption)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.True(t, responseData.Success)

	wizardData := wizard.GetCurrentWizard()
	assert.Equal(t, "cluster-name", wizardData.Info.Name)
	assert.Equal(t, "short-name", wizardData.Info.ShortName)
	assert.Equal(t, wizard.KubeAPIServerConnectTypeLoadBalancer, wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType)
	assert.Equal(t, "", wizardData.Info.KubeAPIServerConnection.VIP)
	assert.Equal(t, "", wizardData.Info.KubeAPIServerConnection.NetInterfaceName)
	assert.Equal(t, "192.168.31.100", wizardData.Info.KubeAPIServerConnection.LoadbalancerIP)
	assert.Equal(t, uint16(4332), wizardData.Info.KubeAPIServerConnection.LoadbalancerPort)
	assert.Equal(t, uint16(16000), wizardData.Info.NodePortMinimum)
	assert.Equal(t, uint16(16999), wizardData.Info.NodePortMaximum)
}
