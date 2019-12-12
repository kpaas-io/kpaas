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

// Service for nodes information manage

package deploy

import (
	"github.com/gin-gonic/gin"

	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/sshcertificate"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
	"github.com/kpaas-io/kpaas/pkg/utils/validator"
)

// @ID GetNodeList
// @Summary Get nodes information
// @Description Get nodes information
// @Tags node
// @Produce application/json
// @Success 200 {object} api.GetNodeListResponse
// @Router /api/v1/deploy/wizard/nodes [get]
func GetNodeList(c *gin.Context) {

	responseData := new(api.GetNodeListResponse)
	nodes := getWizardNodes()
	responseData.Nodes = *nodes

	h.R(c, responseData)
}

// @ID GetNode
// @Summary Get a node information
// @Description Get a node information
// @Tags node
// @Produce application/json
// @Param ip path int true "Node IP Address"
// @Success 200 {object} api.NodeData
// @Failure 400 {object} h.AppErr
// @Failure 404 {object} h.AppErr
// @Router /api/v1/deploy/wizard/nodes/{ip} [get]
func GetNode(c *gin.Context) {

	ip := c.Param("ip")
	if len(ip) <= 0 {

		h.E(c, h.EParamsError.WithPayload("path parameter \"ip\" required"))
		return
	}

	node := wizard.GetCurrentWizard().GetNode(ip)
	if node == nil {

		h.E(c, h.ENotFound.WithPayload("node ip not exist"))
		return
	}

	h.R(c, convertModelNodeToAPINode(node))
}

// @ID AddNode
// @Summary Add Node Information
// @Description Add deployment candidate to node list
// @Tags node
// @Accept application/json
// @Produce application/json
// @Param node body api.NodeData true "Node information"
// @Success 201 {object} api.NodeData
// @Failure 400 {object} h.AppErr
// @Failure 409 {object} h.AppErr
// @Router /api/v1/deploy/wizard/nodes [post]
func AddNode(c *gin.Context) {

	requestData, hasError := getNodeRequestData(c)
	if hasError {
		return
	}

	node := wizard.NewNode()
	node.Name = requestData.Name
	node.Description = requestData.Description
	node.DockerRootDirectory = requestData.DockerRootDirectory
	node.MachineRoles = requestData.MachineRoles

	node.Labels = make([]*wizard.Label, 0, len(requestData.Labels))
	for _, label := range requestData.Labels {

		node.Labels = append(node.Labels, &wizard.Label{
			Key:   label.Key,
			Value: label.Value,
		})
	}

	node.Taints = make([]*wizard.Taint, 0, len(requestData.Taints))
	for _, taint := range requestData.Taints {

		node.Taints = append(node.Taints, &wizard.Taint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: convertAPITaintEffectToModelTaintEffect(taint.Effect),
		})
	}

	node.IP = requestData.IP
	node.Port = requestData.Port
	node.Username = requestData.Username
	node.AuthenticationType = convertAPIAuthenticationTypeToModelAuthenticationType(requestData.AuthenticationType)
	switch requestData.AuthenticationType {
	case api.AuthenticationTypePassword:
		node.Password = requestData.Password
	case api.AuthenticationTypePrivateKey:
		node.PrivateKeyName = requestData.PrivateKeyName
	}

	err := wizard.GetCurrentWizard().AddNode(node)
	if err != nil {
		h.E(c, err)
		log.ReqEntry(c).Info(err)
		return
	}

	h.R(c, requestData)
}

// @ID UpdateNode
// @Summary Update Node Information
// @Description Update a node information which in deployment candidate node list
// @Tags node
// @Accept application/json
// @Produce application/json
// @Param node body api.UpdateNodeData true "Node information"
// @Param ip path int true "Node IP Address"
// @Success 200 {object} api.NodeData
// @Failure 400 {object} h.AppErr
// @Failure 404 {object} h.AppErr
// @Failure 409 {object} h.AppErr
// @Router /api/v1/deploy/wizard/nodes/{ip} [put]
func UpdateNode(c *gin.Context) {

	requestData, ip, hasError := getUpdateNodeRequestData(c)
	if hasError {
		return
	}

	var node = wizard.NewNode()
	node.Name = requestData.Name
	node.Description = requestData.Description
	node.DockerRootDirectory = requestData.DockerRootDirectory
	node.MachineRoles = requestData.MachineRoles

	node.Labels = make([]*wizard.Label, 0, len(requestData.Labels))
	for _, label := range requestData.Labels {

		node.Labels = append(node.Labels, &wizard.Label{
			Key:   label.Key,
			Value: label.Value,
		})
	}

	node.Taints = make([]*wizard.Taint, 0, len(requestData.Taints))
	for _, taint := range requestData.Taints {

		node.Taints = append(node.Taints, &wizard.Taint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: convertAPITaintEffectToModelTaintEffect(taint.Effect),
		})
	}

	node.IP = ip
	node.Port = requestData.Port
	node.Username = requestData.Username
	node.AuthenticationType = convertAPIAuthenticationTypeToModelAuthenticationType(requestData.AuthenticationType)
	switch requestData.AuthenticationType {
	case api.AuthenticationTypePassword:
		node.Password = requestData.Password
	case api.AuthenticationTypePrivateKey:
		node.PrivateKeyName = requestData.PrivateKeyName
	}

	err := wizard.GetCurrentWizard().UpdateNode(node)
	if err != nil {
		h.E(c, err)
		log.ReqEntry(c).Info(err)
		return
	}

	h.R(c, api.NodeData{
		NodeBaseData: requestData.NodeBaseData,
		ConnectionData: api.ConnectionData{
			IP:           ip,
			Port:         requestData.Port,
			SSHLoginData: requestData.SSHLoginData,
		},
	})
}

// @ID DeleteNode
// @Summary Delete a node
// @Description Delete a node from deployment candidate node list
// @Tags node
// @Produce application/json
// @Param ip path int true "Node IP Address"
// @Success 204
// @Failure 400 {object} h.AppErr
// @Failure 404 {object} h.AppErr
// @Failure 409 {object} h.AppErr
// @Router /api/v1/deploy/wizard/nodes/{ip} [delete]
func DeleteNode(c *gin.Context) {

	ip := c.Param("ip")
	if len(ip) <= 0 {
		h.E(c, h.EParamsError.WithPayload("path parameter \"ip\" required"))
		return
	}

	err := wizard.GetCurrentWizard().DeleteNode(ip)
	if err != nil {
		h.E(c, err)
		log.ReqEntry(c).Info(err)
		return
	}

	h.R(c, nil)
}

func getNodeRequestData(c *gin.Context) (*api.NodeData, bool) {

	requestData := new(api.NodeData)
	logger := log.ReqEntry(c)

	if err := validator.Params(c, requestData); err != nil {
		logger.Info(err)
		h.E(c, err)
		return nil, true
	}

	if requestData.AuthenticationType == api.AuthenticationTypePrivateKey {
		validateFunction := validator.ValidateStringOptions(requestData.PrivateKeyName, "privateKeyName", sshcertificate.GetNameList())
		if err := validateFunction(); err != nil {
			h.E(c, h.EParamsError.WithPayload(err))
			return nil, true
		}
	}

	logger = logger.WithField("data", requestData)
	logger.Debug("Request data")

	if requestData.DockerRootDirectory == "" {
		requestData.DockerRootDirectory = wizard.DefaultDockerRootDirectory
		logger.Debug("Use default docker root directory")
	}

	return requestData, false
}

func getUpdateNodeRequestData(c *gin.Context) (*api.UpdateNodeData, string, bool) {

	requestData := new(api.UpdateNodeData)
	logger := log.ReqEntry(c)

	ip := c.Param("ip")
	if len(ip) <= 0 {

		h.E(c, h.EParamsError.WithPayload("path parameter \"ip\" required"))
		return nil, "", true
	}

	if err := validator.Params(c, requestData); err != nil {
		logger.Info(err)
		h.E(c, err)
		return nil, "", true
	}

	if requestData.AuthenticationType == api.AuthenticationTypePrivateKey {
		validateFunction := validator.ValidateStringOptions(requestData.PrivateKeyName, "privateKeyName", sshcertificate.GetNameList())
		if err := validateFunction(); err != nil {
			h.E(c, h.EParamsError.WithPayload(err))
			return nil, "", true
		}
	}

	logger.WithField("data", requestData)
	return requestData, ip, false
}
