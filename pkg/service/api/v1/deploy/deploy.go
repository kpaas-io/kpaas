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

// @ID LaunchDeployment
// @Summary Launch deployment
// @Description Launch deployment
// @Tags deploy
// @Produce application/json
// @Success 201 {object} api.SuccessfulOption
// @Failure 404 {object} h.AppErr
// @Router /api/v1/deploy/wizard/deploys [post]
func Deploy(c *gin.Context) {

	wizardData := wizard.GetCurrentWizard()
	if len(wizardData.Nodes) <= 0 {
		h.E(c, h.ENotFound.WithPayload("No node information, node list is empty, please add node information"))
		return
	}

	if wizardData.GetCheckResult() != constant.CheckResultPassed {
		h.E(c, h.EStatusError.WithPayload("current check result status is not passed"))
		return
	}

	switch wizardData.GetDeployClusterStatus() {
	case wizard.DeployClusterStatusSuccessful:
		h.E(c, h.EStatusError.WithPayload("It was deploy succeed"))
		return
	case wizard.DeployClusterStatusRunning:
		h.E(c, h.EStatusError.WithPayload("It was deploying"))
		return
	}

	wizardData.ClearClusterDeployData()

	if err := wizardData.MarkNodeDeploying(); err != nil {
		h.E(c, h.EStatusError.WithPayload(err))
		return
	}

	client := clientUtils.GetDeployController()

	grpcContext, cancel := context.WithTimeout(context.Background(), config.Config.DeployController.GetTimeout())
	defer cancel()

	resp, err := client.Deploy(grpcContext, getCallDeployData())
	if err != nil {
		h.E(c, h.EDeployControllerError.WithPayload(err))
		log.ReqEntry(c).Errorf("call deploy controller error, errorMessage: %v", err)
		wizardData.ClearClusterDeployData()
		return
	}

	if resp.GetErr() != nil {

		log.ReqEntry(c).Errorf("call deploy controller result error, error: %#v", resp.GetErr())
	}

	go listenDeploymentData()

	h.R(c, api.SuccessfulOption{Success: resp.GetAcceptd()})
}

// @ID GetDeploymentReport
// @Summary Get the result of deployment
// @Description Get the result of the deployment
// @Tags deploy
// @Produce application/json
// @Success 200 {object} api.GetDeploymentReportResponse
// @Router /api/v1/deploy/wizard/deploys [get]
func GetDeployReport(c *gin.Context) {

	wizardData := wizard.GetCurrentWizard()
	nodeList := getWizardDeploymentData()
	responseData := api.GetDeploymentReportResponse{
		Roles:               *nodeList,
		DeployClusterStatus: convertModelDeployClusterStatusToAPIDeployClusterStatus(wizardData.DeployClusterStatus),
		DeployClusterError:  convertModelErrorToAPIError(wizardData.DeployClusterError),
	}

	h.R(c, responseData)
}

func getCallDeployData() *protos.DeployRequest {

	return &protos.DeployRequest{
		NodeConfigs:   buildCallDeployDataNodesPart(),
		ClusterConfig: buildCallDeployDataClusterPart(),
	}
}

func buildCallDeployDataNodesPart() (nodeConfigs []*protos.NodeDeployConfig) {

	wizardData := wizard.GetCurrentWizard()
	nodeConfigs = make([]*protos.NodeDeployConfig, 0, len(wizardData.Nodes))
	for _, node := range wizardData.Nodes {

		nodeConfig := new(protos.NodeDeployConfig)
		for _, role := range node.MachineRoles {
			nodeConfig.Roles = append(nodeConfig.Roles, string(role))
		}

		nodeConfig.Node = &protos.Node{
			Name: node.Name,
			Ip:   node.IP,
			Ssh:  convertModelConnectionDataToDeployControllerSSHData(&node.ConnectionData),
		}

		nodeConfig.Labels = make(map[string]string)
		for _, label := range node.Labels {
			nodeConfig.Labels[label.Key] = label.Value
		}

		nodeConfig.Taints = make([]*protos.Taint, 0, len(node.Taints))
		for _, taint := range node.Taints {
			nodeConfig.Taints = append(nodeConfig.Taints, &protos.Taint{
				Key:    taint.Key,
				Value:  taint.Value,
				Effect: string(taint.Effect),
			})
		}

		nodeConfigs = append(nodeConfigs, nodeConfig)
	}

	return
}

func buildCallDeployDataClusterPart() (clusterConfig *protos.ClusterConfig) {

	wizardData := wizard.GetCurrentWizard()
	clusterConfig = &protos.ClusterConfig{
		ClusterName: wizardData.Info.ShortName,
		KubeAPIServerConnect: &protos.KubeAPIServerConnect{
			Type:         string(wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType),
			Keepalived:   nil,
			Loadbalancer: nil,
		},
		NodePortRange: &protos.NodePortRange{
			From: uint32(wizardData.Info.NodePortMinimum),
			To:   uint32(wizardData.Info.NodePortMaximum),
		},
		NodeLabels:      make(map[string]string),
		NodeAnnotations: make(map[string]string),
	}

	switch wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType {
	case wizard.KubeAPIServerConnectTypeKeepalived:
		clusterConfig.KubeAPIServerConnect.Keepalived = &protos.Keepalived{
			Vip:              wizardData.Info.KubeAPIServerConnection.VIP,
			NetInterfaceName: wizardData.Info.KubeAPIServerConnection.NetInterfaceName,
		}
	case wizard.KubeAPIServerConnectTypeLoadBalancer:
		clusterConfig.KubeAPIServerConnect.Loadbalancer = &protos.Loadbalancer{
			Ip:   wizardData.Info.KubeAPIServerConnection.LoadbalancerIP,
			Port: uint32(wizardData.Info.KubeAPIServerConnection.LoadbalancerPort),
		}
	}

	for _, label := range wizardData.Info.Labels {
		clusterConfig.NodeLabels[label.Key] = label.Value
	}

	for _, annotation := range wizardData.Info.Annotations {
		clusterConfig.NodeAnnotations[annotation.Key] = annotation.Value
	}

	return
}

func listenDeploymentData() {

	wizardData := wizard.GetCurrentWizard()
	for {
		if wizardData.GetDeployClusterStatus() != wizard.DeployClusterStatusRunning {
			break
		}

		refreshDeployResultOneTime()
		time.Sleep(time.Second)
	}
}

func refreshDeployResultOneTime() {

	client := clientUtils.GetDeployController()

	grpcContext, cancel := context.WithTimeout(context.Background(), config.Config.DeployController.GetTimeout())
	defer cancel()

	resp, err := client.GetDeployResult(grpcContext, &protos.GetDeployResultRequest{WithLogs: true})
	if err != nil {
		logrus.Errorf("call deploy controller error, errorMessage: %v", err)
		return
	}

	wizardData := wizard.GetCurrentWizard()
	wizardData.SetClusterDeploymentStatus(
		convertDeployControllerDeployClusterStatusToModelDeployClusterStatus(resp.GetStatus()),
		convertDeployControllerErrorToFailureDetail(resp.GetErr()))

	for _, item := range resp.Items {

		wizardNode := wizardData.GetNodeByName(item.DeployItem.NodeName)
		if wizardNode == nil {

			logrus.Errorf("iterate deployment result, can not find node(%s) from cluster data", item.DeployItem.NodeName)
			continue
		}

		failureDetail := convertDeployControllerErrorToFailureDetail(item.GetErr())

		if failureDetail != nil && item.Logs != "" {
			var setLogContentError error
			failureDetail.LogId, setLogContentError = wizard.SetLogByString(item.Logs)
			if setLogContentError != nil {
				logrus.Errorf("Store deploy error log error, %s", setLogContentError)
			}
		}

		wizardNode.SetDeployResult(
			constant.MachineRole(item.DeployItem.Role),
			convertDeployControllerDeployResultToModelDeployResult(item.GetStatus()),
			failureDetail)
	}

	switch wizardData.DeployClusterStatus {
	case wizard.DeployClusterStatusSuccessful, wizard.DeployClusterStatusWorkedButHaveError:
		// TODO Lucky Update KubeConfig Data
	}
}
