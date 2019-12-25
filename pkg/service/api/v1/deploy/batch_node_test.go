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
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

func TestUploadBatchNodes(t *testing.T) {

	logrus.SetLevel(logrus.TraceLevel)
	tests := []struct {
		BaseNodeList  []*wizard.Node
		Input         string
		Want          *api.GetNodeListResponse
		ResponseError *h.AppErr
	}{
		{
			BaseNodeList: []*wizard.Node{},
			Input: `
#<hostname> <user>  <role,role,role>         <IP>             <ssh port>  <password>          <login key name>        <docker path>
k8s-master1   root	    master,etcd          192.168.3.223    22          111111111111	      -		                  /var/lib/docker
k8s-master2   root	    master,etcd          192.168.3.224    22          111111111111	      -   	                  /var/lib/docker
`,
			Want: &api.GetNodeListResponse{
				Nodes: []api.NodeData{
					{
						NodeBaseData: api.NodeBaseData{
							Name:                "k8s-master1",
							MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
							DockerRootDirectory: "/var/lib/docker",
							Labels:              []api.Label{},
							Taints:              []api.Taint{},
						},
						ConnectionData: api.ConnectionData{
							SSHLoginData: api.SSHLoginData{
								Username:           "root",
								AuthenticationType: api.AuthenticationTypePassword,
								Password:           "",
							},
							IP:   "192.168.3.223",
							Port: 22,
						},
					},
					{
						NodeBaseData: api.NodeBaseData{
							Name:                "k8s-master2",
							MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
							DockerRootDirectory: "/var/lib/docker",
							Labels:              []api.Label{},
							Taints:              []api.Taint{},
						},
						ConnectionData: api.ConnectionData{
							SSHLoginData: api.SSHLoginData{
								Username:           "root",
								AuthenticationType: api.AuthenticationTypePassword,
								Password:           "",
							},
							IP:   "192.168.3.224",
							Port: 22,
						},
					},
				},
			},
		},
		{
			BaseNodeList: []*wizard.Node{},
			Input: `
#<hostname> <user>  <role,role,role>         <IP>             <ssh port>  <password>          <login key name>        <docker path>
k8s-master1   root	    master,etcd          192.168.3.223    22          -	      TestKey		                  /var/lib/docker
k8s-master2   root	    master,etcd          192.168.3.224    22          -	      TestKey   	                  /var/lib/docker
`,
			Want: &api.GetNodeListResponse{
				Nodes: []api.NodeData{
					{
						NodeBaseData: api.NodeBaseData{
							Name:                "k8s-master1",
							MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
							DockerRootDirectory: "/var/lib/docker",
							Labels:              []api.Label{},
							Taints:              []api.Taint{},
						},
						ConnectionData: api.ConnectionData{
							SSHLoginData: api.SSHLoginData{
								Username:           "root",
								AuthenticationType: api.AuthenticationTypePrivateKey,
								PrivateKeyName:     "TestKey",
							},
							IP:   "192.168.3.223",
							Port: 22,
						},
					},
					{
						NodeBaseData: api.NodeBaseData{
							Name:                "k8s-master2",
							MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
							DockerRootDirectory: "/var/lib/docker",
							Labels:              []api.Label{},
							Taints:              []api.Taint{},
						},
						ConnectionData: api.ConnectionData{
							SSHLoginData: api.SSHLoginData{
								Username:           "root",
								AuthenticationType: api.AuthenticationTypePrivateKey,
								PrivateKeyName:     "TestKey",
							},
							IP:   "192.168.3.224",
							Port: 22,
						},
					},
				},
			},
		},
		{
			BaseNodeList: []*wizard.Node{},
			Input: `
#<hostname>    <user>    <role,role,role>    <IP>             <ssh port>    <password>       <login key name>    <docker path>
k8s-master1    root      master,etcd         192.168.3.223    22             111111111111    -                   /var/lib/docker
k8s-master2    root      master,etcd         192.168.3.224    22             111111111111    -                   /var/lib/docker
k8s-master3    root      master,etcd         192.168.3.227    22             111111111111    -                   /var/lib/docker
k8s-worker4    root      worker,etcd         192.168.3.226    22             -               worker_key          /var/lib/docker
k8s-worker5    root      worker,etcd         192.168.3.229    22             -               worker_key          /var/lib/docker
k8s-worker6    root      worker              192.168.3.230    22             -               worker_key          /var/lib/docker
`,
			Want: &api.GetNodeListResponse{
				Nodes: []api.NodeData{
					{
						NodeBaseData: api.NodeBaseData{
							Name:                "k8s-master1",
							MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
							DockerRootDirectory: "/var/lib/docker",
							Labels:              []api.Label{},
							Taints:              []api.Taint{},
						},
						ConnectionData: api.ConnectionData{
							SSHLoginData: api.SSHLoginData{
								Username:           "root",
								AuthenticationType: api.AuthenticationTypePassword,
							},
							IP:   "192.168.3.223",
							Port: 22,
						},
					},
					{
						NodeBaseData: api.NodeBaseData{
							Name:                "k8s-master2",
							MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
							DockerRootDirectory: "/var/lib/docker",
							Labels:              []api.Label{},
							Taints:              []api.Taint{},
						},
						ConnectionData: api.ConnectionData{
							SSHLoginData: api.SSHLoginData{
								Username:           "root",
								AuthenticationType: api.AuthenticationTypePassword,
							},
							IP:   "192.168.3.224",
							Port: 22,
						},
					},
					{
						NodeBaseData: api.NodeBaseData{
							Name:                "k8s-master3",
							MachineRoles:        []constant.MachineRole{constant.MachineRoleMaster, constant.MachineRoleEtcd},
							DockerRootDirectory: "/var/lib/docker",
							Labels:              []api.Label{},
							Taints:              []api.Taint{},
						},
						ConnectionData: api.ConnectionData{
							SSHLoginData: api.SSHLoginData{
								Username:           "root",
								AuthenticationType: api.AuthenticationTypePassword,
							},
							IP:   "192.168.3.227",
							Port: 22,
						},
					},
					{
						NodeBaseData: api.NodeBaseData{
							Name:                "k8s-worker4",
							MachineRoles:        []constant.MachineRole{constant.MachineRoleWorker, constant.MachineRoleEtcd},
							DockerRootDirectory: "/var/lib/docker",
							Labels:              []api.Label{},
							Taints:              []api.Taint{},
						},
						ConnectionData: api.ConnectionData{
							SSHLoginData: api.SSHLoginData{
								Username:           "root",
								AuthenticationType: api.AuthenticationTypePrivateKey,
								PrivateKeyName:     "worker_key",
							},
							IP:   "192.168.3.226",
							Port: 22,
						},
					},
					{
						NodeBaseData: api.NodeBaseData{
							Name:                "k8s-worker5",
							MachineRoles:        []constant.MachineRole{constant.MachineRoleWorker, constant.MachineRoleEtcd},
							DockerRootDirectory: "/var/lib/docker",
							Labels:              []api.Label{},
							Taints:              []api.Taint{},
						},
						ConnectionData: api.ConnectionData{
							SSHLoginData: api.SSHLoginData{
								Username:           "root",
								AuthenticationType: api.AuthenticationTypePrivateKey,
								PrivateKeyName:     "worker_key",
							},
							IP:   "192.168.3.229",
							Port: 22,
						},
					},
					{
						NodeBaseData: api.NodeBaseData{
							Name:                "k8s-worker6",
							MachineRoles:        []constant.MachineRole{constant.MachineRoleWorker},
							DockerRootDirectory: "/var/lib/docker",
							Labels:              []api.Label{},
							Taints:              []api.Taint{},
						},
						ConnectionData: api.ConnectionData{
							SSHLoginData: api.SSHLoginData{
								Username:           "root",
								AuthenticationType: api.AuthenticationTypePrivateKey,
								PrivateKeyName:     "worker_key",
							},
							IP:   "192.168.3.230",
							Port: 22,
						},
					},
				},
			},
		},
		{
			BaseNodeList: []*wizard.Node{
				{
					Name: "k8s-master1",
				},
			},
			Input: `
#<hostname> <user>  <role,role,role>         <IP>             <ssh port>  <password>          <login key name>        <docker path>
k8s-master1   root	    master,etcd          192.168.3.223    22          111111111111	      -		                  /var/lib/docker
`,
			Want:          nil,
			ResponseError: h.EExists.WithPayload(fmt.Sprintf("node name k8s-master1 was exist")),
		},
		{
			BaseNodeList: []*wizard.Node{},
			Input: `
#<hostname> <user>  <role,role,role>         <IP>             <ssh port>  <password>          <login key name>        <docker path>
k8s-master1   -	    master,etcd          192.168.3.223    22          111111111111	      -		                  /var/lib/docker
`,
			Want:          nil,
			ResponseError: h.EParamsError.WithPayload(fmt.Sprintf(`"username" illegal`)),
		},
		{
			BaseNodeList:  []*wizard.Node{},
			Input:         ``,
			Want:          nil,
			ResponseError: h.EParamsError.WithPayload(fmt.Sprintf(`node list empty`)),
		},
		{
			BaseNodeList: []*wizard.Node{},
			Input: `
k8s-master1   root	    master,etcd          192.168.3.223    22          111111111111	      -		                  /var/lib/docker
k8s-master1   root	    master,etcd          192.168.3.224    22          111111111111	      -		                  /var/lib/docker
`,
			Want:          nil,
			ResponseError: h.EParamsError.WithPayload(fmt.Sprintf(`node name k8s-master1 was duplicated`)),
		},
		{
			BaseNodeList: []*wizard.Node{},
			Input: `
k8s-master1   root	    master,etcd          192.168.3.223    22          111111111111	      -		                  /var/lib/docker
k8s-master2   root	    master,etcd          192.168.3.223    22          111111111111	      -		                  /var/lib/docker
`,
			Want:          nil,
			ResponseError: h.EParamsError.WithPayload(fmt.Sprintf(`node ip 192.168.3.223 was duplicated`)),
		},
		{
			BaseNodeList: []*wizard.Node{},
			Input: `
k8s-master1   root	    master,etcd          192.168.3.223    99999999999999999999999999999          111111111111	      -		                  /var/lib/docker
`,
			Want:          nil,
			ResponseError: h.EParamsError.WithPayload(fmt.Sprintf(`strconv.Atoi: parsing "99999999999999999999999999999": value out of range`)),
		},
	}

	for _, test := range tests {

		wizard.ClearCurrentWizardData()
		wizardData := wizard.GetCurrentWizard()
		wizardData.Nodes = test.BaseNodeList
		var err error
		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		bodyReader := strings.NewReader(test.Input)
		ctx.Request = httptest.NewRequest("POST", "/api/v1/deploy/wizard/batchnodes", bodyReader)

		UploadBatchNodes(ctx)

		resp.Flush()
		assert.True(t, resp.Body.Len() > 0)
		fmt.Printf("result: %s\n", resp.Body.String())
		if test.Want != nil {
			responseData := new(api.GetNodeListResponse)
			err = json.Unmarshal(resp.Body.Bytes(), responseData)
			assert.Nil(t, err)

			assert.Equal(t, test.Want, responseData)
		} else {
			responseData := new(h.AppErr)
			err = json.Unmarshal(resp.Body.Bytes(), responseData)
			assert.Nil(t, err)

			assert.Equal(t, test.ResponseError.Error(), responseData.Error())
			assert.Equal(t, test.ResponseError.Status, resp.Code)
		}
	}
}

func TestTryToMatchBatchNodes(t *testing.T) {

	tests := []struct {
		Input          []byte
		WantMatches    [][]string
		WantGroupNames []string
	}{
		{
			Input: []byte(`
#<hostname> <user>  <role,role,role>         <IP>             <ssh port>  <password>          <login key name>        <docker path>
k8s-master1   root	    master,etcd          192.168.3.223    22          111111111111	      -		                  /var/lib/docker
`),
			WantMatches: [][]string{
				{
					"k8s-master1   root	    master,etcd          192.168.3.223    22          111111111111	      -		                  /var/lib/docker",
					"k8s-master1",
					"root",
					"master,etcd",
					"192.168.3.223",
					"22",
					"111111111111",
					"-",
					"/var/lib/docker",
				},
			},
			WantGroupNames: []string{
				"",
				"nodeName",
				"username",
				"roles",
				"ip",
				"port",
				"password",
				"privateKeyName",
				"dockerPath",
			},
		},
	}

	for _, test := range tests {

		matches, groupNames := tryToMatchBatchNodes(test.Input)
		assert.Equal(t, test.WantMatches, matches)
		assert.Equal(t, test.WantGroupNames, groupNames)
	}
}

func TestSplitInputRoles(t *testing.T) {

	tests := []struct {
		Input map[string]string
		Want  []constant.MachineRole
	}{
		{
			Input: map[string]string{
				"roles": "master,etcd",
			},
			Want: []constant.MachineRole{
				constant.MachineRoleMaster,
				constant.MachineRoleEtcd,
			},
		},
		{
			Input: map[string]string{
				"roles": "worker,etcd",
			},
			Want: []constant.MachineRole{
				constant.MachineRoleWorker,
				constant.MachineRoleEtcd,
			},
		},
		{
			Input: map[string]string{
				"roles": "",
			},
			Want: []constant.MachineRole{},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Want, splitInputRoles(test.Input))
	}
}
