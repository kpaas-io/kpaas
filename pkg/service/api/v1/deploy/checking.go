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
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/service/config"
	clientUtils "github.com/kpaas-io/kpaas/pkg/service/grpcutils/client"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
)

// @ID CheckNodeList
// @Summary check node list
// @Description Check if the node meets the pre-deployment requirements
// @Tags checking
// @Produce application/json
// @Success 201 {object} api.SuccessfulOption
// @Router /api/v1/deploy/wizard/checks [post]
func CheckNodeList(c *gin.Context) {

	wizardData := wizard.GetCurrentWizard()
	if len(wizardData.Nodes) == 0 {
		h.E(c, h.ENotFound.WithPayload("No node information, node list is empty, please add node information"))
		return
	}

	if wizardData.GetCheckResult() == constant.CheckResultRunning {
		h.E(c, h.EStatusError.WithPayload("It was checking"))
		return
	}

	if !checkClusterConfiguration() {

		// Cluster Configuration check failed, no need to check the nodes
		// Return true because this is a go check trigger API
		h.R(c, api.SuccessfulOption{Success: true})
		return
	}

	wizardData.ClearClusterCheckingData()

	if err := wizardData.MarkNodeChecking(); err != nil {
		h.E(c, h.EStatusError.WithPayload(err))
		return
	}

	client := clientUtils.GetDeployController()

	grpcContext, cancel := context.WithTimeout(context.Background(), config.Config.DeployController.GetTimeout())
	defer cancel()

	resp, err := client.CheckNodes(grpcContext, getCallCheckNodesData())
	if err != nil {
		h.E(c, h.EDeployControllerError.WithPayload(err))
		log.ReqEntry(c).Errorf("call deploy controller error, errorMessage: %v", err)
		wizardData.ClearClusterDeployData()
		return
	}

	if resp.GetErr() != nil {

		log.ReqEntry(c).Errorf("call deploy controller result error, error: %#v", resp.GetErr())
	}

	go listenCheckNodesData()

	h.R(c, api.SuccessfulOption{Success: resp.GetAccepted()})
}

// @ID GetCheckNodeListResult
// @Summary Get the result of check node
// @Description Get the result of the check node
// @Tags checking
// @Produce application/json
// @Success 200 {object} api.GetCheckingResultResponse
// @Router /api/v1/deploy/wizard/checks [get]
func GetCheckingNodeListResult(c *gin.Context) {

	responseData := new(api.GetCheckingResultResponse)
	wizardData := wizard.GetCurrentWizard()
	checkResults := getWizardCheckingData()
	responseData.Nodes = *checkResults
	responseData.Result = wizardData.GetCheckResult()
	responseData.Cluster = getCheckedClusterConfiguration()

	h.R(c, responseData)
}

func getCallCheckNodesData() *protos.CheckNodesRequest {

	requestData := &protos.CheckNodesRequest{}

	wizardData := wizard.GetCurrentWizard()
	for _, node := range wizardData.Nodes {

		nodeConfig := new(protos.NodeCheckConfig)
		for _, role := range node.MachineRoles {
			nodeConfig.Roles = append(nodeConfig.Roles, string(role))
		}

		nodeConfig.Node = &protos.Node{
			Name: node.Name,
			Ip:   node.IP,
			Ssh:  convertModelConnectionDataToDeployControllerSSHData(&node.ConnectionData),
		}

		requestData.Configs = append(requestData.Configs, nodeConfig)
	}

	return requestData
}

func listenCheckNodesData() {

	wizardData := wizard.GetCurrentWizard()
	for {
		if wizardData.GetCheckResult() != constant.CheckResultRunning {
			break
		}

		refreshCheckResultOneTime()
		time.Sleep(time.Second)
	}
}

func refreshCheckResultOneTime() {

	client := clientUtils.GetDeployController()

	grpcContext, cancel := context.WithTimeout(context.Background(), config.Config.DeployController.GetTimeout())
	defer cancel()

	resp, err := client.GetCheckNodesResult(grpcContext, &protos.GetCheckNodesResultRequest{})
	if err != nil {
		logrus.Errorf("call deploy controller error, errorMessage: %v", err)
		return
	}

	wizardData := wizard.GetCurrentWizard()
	wizardData.SetClusterCheckResult(
		convertDeployControllerCheckResultToModelCheckResult(resp.GetStatus()),
		convertDeployControllerErrorToFailureDetail(resp.GetErr()))

	for _, node := range resp.Nodes {

		wizardNode := wizardData.GetNodeByName(node.NodeName)
		if wizardNode == nil {

			logrus.Errorf("iterate check result, can not find node(%s) from cluster data", node.NodeName)
			continue
		}

		wizardNode.SetCheckResult(
			convertDeployControllerCheckResultToModelCheckResult(node.GetStatus()),
			convertDeployControllerErrorToFailureDetail(node.GetErr()))

		for _, item := range node.Items {
			itemName := getItemNameFromDeployControllerCheckItem(item.Item)
			failureDetail := convertDeployControllerErrorToFailureDetail(item.Err)
			if failureDetail != nil && item.Logs != "" {
				var setLogContentError error
				failureDetail.LogId, setLogContentError = wizard.SetLogByString(item.Logs)
				if setLogContentError != nil {
					logrus.Errorf("Store checking error log error, %s", setLogContentError)
				}
			}
			wizardNode.SetCheckItem(itemName, convertDeployControllerCheckResultToModelCheckResult(item.Status), failureDetail)
		}
	}
}

func getItemNameFromDeployControllerCheckItem(item *protos.CheckItem) string {

	if item == nil {
		return ""
	}

	var itemName string
	if item.Name != "" {
		itemName = item.Name
	}

	if item.Description == "" {
		return itemName
	}

	if itemName != "" {
		return fmt.Sprintf("%s（%s）", itemName, item.Description)
	}

	return item.Description
}

func checkClusterConfiguration() bool {

	return len(checkWrongClusterConfiguration()) == 0
}

func checkWrongClusterConfiguration() (errs []*api.CheckingItem) {

	errs = make([]*api.CheckingItem, 0)
	wizardData := wizard.GetCurrentWizard()
	if len(wizardData.Nodes) == 0 {

		errs = append(errs, &api.CheckingItem{
			CheckingPoint: "Checking node information", // 检查节点信息
			Result:        constant.CheckResultFailed,
			Error: &api.Error{
				Reason:     "No node information",         // 无节点信息
				Detail:     "node list is empty",          // 节点列表为空
				FixMethods: "please add node information", // 请添加节点信息
			},
		})
		return
	}

	counters := map[constant.MachineRole]uint{
		constant.MachineRoleEtcd:    0,
		constant.MachineRoleMaster:  0,
		constant.MachineRoleWorker:  0,
		constant.MachineRoleIngress: 0,
	}
	for _, node := range wizardData.Nodes {

		for _, role := range node.MachineRoles {
			counters[role]++
		}
	}

	for role, counter := range counters {

		if counter > 0 {
			continue
		}

		errs = append(errs, &api.CheckingItem{
			CheckingPoint: fmt.Sprintf("Checking nodes for %s count", role),
			Result:        constant.CheckResultFailed,
			Error: &api.Error{
				Reason:     fmt.Sprintf("nodes for %s are not enough", role),                                           // %s 角色的节点数不足
				Detail:     fmt.Sprintf("%s needs at least one node", role),                                            // %s 角色节点数至少一个
				FixMethods: fmt.Sprintf("Add new node for role: %s, or edit existing node to include this role", role), // 添加一个新节点包含 %s 角色，或者编辑已有节点使他包含这个角色
			},
		})
	}

	return errs
}

func getCheckedClusterConfiguration() api.CheckClusterResponseData {
	return api.CheckClusterResponseData{
		Items: checkWrongClusterConfiguration(),
	}
}
