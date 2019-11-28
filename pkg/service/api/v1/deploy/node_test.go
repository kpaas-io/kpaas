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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

func TestAddNode(t *testing.T) {

	wizard.ClearCurrentWizardData()
	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	body := api.NodeData{
		NodeBaseData: api.NodeBaseData{
			Name:         "name",
			Description:  "description",
			MachineRoles: []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
			Labels: []api.Label{
				{
					Key:   "label-key",
					Value: "value",
				},
			},
			Taints: []api.Taint{
				{
					Key:    "taint-key",
					Value:  "value",
					Effect: api.TaintEffectNoExecute,
				},
			},
			DockerRootDirectory: "/var/lib/docker",
		},
		ConnectionData: api.ConnectionData{
			IP:   "192.168.30.140",
			Port: uint16(22),
			SSHLoginData: api.SSHLoginData{
				Username:           "root",
				AuthenticationType: api.AuthenticationTypePassword,
				Password:           "123456",
			},
		},
	}
	bodyContent, err := json.Marshal(body)
	assert.Nil(t, err)
	bodyReader := bytes.NewReader(bodyContent)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/nodes", bodyReader)

	AddNode(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.NodeData)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, body, *responseData)
}

func TestUpdateNode(t *testing.T) {

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
				IP:                 "192.168.31.140",
				Port:               22,
				Username:           "kpaas",
				AuthenticationType: wizard.AuthenticationTypePassword,
				Password:           "123456",
			},
		},
	}

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	body := api.NodeData{
		NodeBaseData: api.NodeBaseData{
			Name:         "name",
			Description:  "description",
			MachineRoles: []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
			Labels: []api.Label{
				{
					Key:   "label-key",
					Value: "value",
				},
			},
			Taints: []api.Taint{
				{
					Key:    "taint-key",
					Value:  "value",
					Effect: api.TaintEffectNoExecute,
				},
			},
			DockerRootDirectory: "/var/lib/docker",
		},
		ConnectionData: api.ConnectionData{
			IP:   "192.168.31.140",
			Port: uint16(22),
			SSHLoginData: api.SSHLoginData{
				Username:           "root",
				AuthenticationType: api.AuthenticationTypePassword,
				Password:           "123456",
			},
		},
	}
	bodyContent, err := json.Marshal(body)
	assert.Nil(t, err)
	bodyReader := bytes.NewReader(bodyContent)
	ctx.Request = httptest.NewRequest("PUT", "/api/v1/deploy/wizard/nodes/192.168.31.140", bodyReader)
	ctx.Params = gin.Params{
		{
			Key:   "ip",
			Value: "192.168.31.140",
		},
	}

	UpdateNode(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.NodeData)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, body, *responseData)
}

func TestUpdateNode_NotExist(t *testing.T) {

	wizard.ClearCurrentWizardData()

	var err error
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	body := api.NodeData{
		NodeBaseData: api.NodeBaseData{
			Name:         "name",
			Description:  "description",
			MachineRoles: []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
			Labels: []api.Label{
				{
					Key:   "label-key",
					Value: "value",
				},
			},
			Taints: []api.Taint{
				{
					Key:    "taint-key",
					Value:  "value",
					Effect: api.TaintEffectNoExecute,
				},
			},
			DockerRootDirectory: "/var/lib/docker",
		},
		ConnectionData: api.ConnectionData{
			IP:   "192.168.31.140",
			Port: uint16(22),
			SSHLoginData: api.SSHLoginData{
				Username:           "root",
				AuthenticationType: api.AuthenticationTypePassword,
				Password:           "123456",
			},
		},
	}
	bodyContent, err := json.Marshal(body)
	assert.Nil(t, err)
	bodyReader := bytes.NewReader(bodyContent)
	ctx.Request = httptest.NewRequest("PUT", "/api/v1/deploy/wizard/nodes/192.168.31.140", bodyReader)
	ctx.Params = gin.Params{
		{
			Key:   "ip",
			Value: "192.168.31.140",
		},
	}

	UpdateNode(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(h.AppErr)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, h.ENotFound.Status, resp.Code)
}

func TestDeleteNode(t *testing.T) {

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
				IP:                 "192.168.31.140",
				Port:               22,
				Username:           "kpaas",
				AuthenticationType: wizard.AuthenticationTypePassword,
				Password:           "123456",
			},
		},
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("DELETE", "/api/v1/deploy/wizard/nodes/192.168.31.140", nil)
	ctx.Params = gin.Params{
		{
			Key:   "ip",
			Value: "192.168.31.140",
		},
	}

	DeleteNode(ctx)
	resp.Flush()

	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestGetNodeList(t *testing.T) {

	var err error
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
				IP:                 "192.168.31.140",
				Port:               22,
				Username:           "kpaas",
				AuthenticationType: wizard.AuthenticationTypePassword,
				Password:           "123456",
			},
		},
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", "/api/v1/deploy/wizard/nodes", nil)

	GetNodeList(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.GetNodeListResponse)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, []api.NodeData{
		{
			NodeBaseData: api.NodeBaseData{
				Name:         "master1",
				Description:  "desc1",
				MachineRoles: []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
				Labels: []api.Label{
					{
						Key:   "kpaas.io/test",
						Value: "yes",
					},
				},
				Taints: []api.Taint{
					{
						Key:    "taint1",
						Value:  "taint-value",
						Effect: api.TaintEffectNoExecute,
					},
				},
				DockerRootDirectory: "/mnt/docker",
			},
			ConnectionData: api.ConnectionData{
				IP:   "192.168.31.140",
				Port: 22,
				SSHLoginData: api.SSHLoginData{
					Username:           "kpaas",
					AuthenticationType: api.AuthenticationTypePassword,
					Password:           "",
				},
			},
		},
	}, responseData.Nodes)
}

func TestGetNode(t *testing.T) {

	var err error
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
				IP:                 "192.168.31.140",
				Port:               22,
				Username:           "kpaas",
				AuthenticationType: wizard.AuthenticationTypePassword,
				Password:           "123456",
			},
		},
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest("GET", "/api/v1/deploy/wizard/nodes/192.168.31.140", nil)
	ctx.Params = gin.Params{
		{
			Key:   "ip",
			Value: "192.168.31.140",
		},
	}

	GetNode(ctx)
	resp.Flush()
	assert.True(t, resp.Body.Len() > 0)
	fmt.Printf("result: %s\n", resp.Body.String())
	responseData := new(api.NodeData)
	err = json.Unmarshal(resp.Body.Bytes(), responseData)
	assert.Nil(t, err)

	assert.Equal(t, api.NodeData{
		NodeBaseData: api.NodeBaseData{
			Name:         "master1",
			Description:  "desc1",
			MachineRoles: []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
			Labels: []api.Label{
				{
					Key:   "kpaas.io/test",
					Value: "yes",
				},
			},
			Taints: []api.Taint{
				{
					Key:    "taint1",
					Value:  "taint-value",
					Effect: api.TaintEffectNoExecute,
				},
			},
			DockerRootDirectory: "/mnt/docker",
		},
		ConnectionData: api.ConnectionData{
			IP:   "192.168.31.140",
			Port: 22,
			SSHLoginData: api.SSHLoginData{
				Username:           "kpaas",
				AuthenticationType: api.AuthenticationTypePassword,
				Password:           "",
			},
		},
	}, *responseData)
}
