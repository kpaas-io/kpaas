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
	"context"

	"github.com/gin-gonic/gin"

	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/service/config"
	clientUtils "github.com/kpaas-io/kpaas/pkg/service/grpcutils/client"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/sshcertificate"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
	"github.com/kpaas-io/kpaas/pkg/utils/validator"
)

const (
	deployControllerAuthCredentialPassword   = "password"
	deployControllerAuthCredentialPrivateKey = "privatekey"
)

// @ID TestSSH
// @Summary Test a node connection
// @Description Try to connection a node using ssh
// @Tags ssh
// @Accept application/json
// @Produce application/json
// @Param node body api.ConnectionData true "Node information"
// @Success 201 {object} api.SuccessfulOption
// @Failure 400 {object} h.AppErr
// @Failure 409 {object} h.AppErr
// @Failure 500 {object} h.AppErr
// @Router /api/v1/ssh/tests [post]
func TestConnectNode(c *gin.Context) {

	requestData, hasError := getConnectionData(c)
	if hasError {
		return
	}

	client := clientUtils.GetDeployController()

	grpcContext, cancel := context.WithTimeout(context.Background(), config.Config.DeployController.GetTimeout())
	defer cancel()

	resp, err := client.TestConnection(grpcContext, getCallTestConnectionData(requestData))
	if err != nil {
		h.E(c, h.EDeployControllerError.WithPayload(err))
		log.ReqEntry(c).Errorf("call deploy controller error, errorMessage: %v", err)
		return
	}

	h.R(c, api.TestConnectionResponse{
		SuccessfulOption: api.SuccessfulOption{Success: resp.GetPassed()},
		Error:            convertDeployControllerErrorToAPIError(resp.GetErr()),
	})
}

func getCallTestConnectionData(requestData *api.ConnectionData) *protos.TestConnectionRequest {

	auth := &protos.Auth{
		UserName: requestData.Username,
	}

	switch requestData.AuthenticationType {
	case api.AuthenticationTypePassword:
		auth.Type = deployControllerAuthCredentialPassword
		auth.Credential = requestData.Password
	case api.AuthenticationTypePrivateKey:
		auth.Type = deployControllerAuthCredentialPrivateKey
		auth.Credential = sshcertificate.GetPrivateKey(requestData.PrivateKeyName)
	}

	return &protos.TestConnectionRequest{Node: &protos.Node{
		Name: requestData.IP,
		Ip:   requestData.IP,
		Ssh: &protos.SSH{
			Port: uint32(requestData.Port),
			Auth: auth,
		},
	}}
}

func getConnectionData(c *gin.Context) (requestData *api.ConnectionData, hasError bool) {

	requestData = new(api.ConnectionData)
	logger := log.ReqEntry(c)

	if err := validator.Params(c, requestData); err != nil {
		logger.Info(err)
		h.E(c, err)
		return nil, true
	}

	logger.WithField("data", requestData)
	return requestData, false
}
