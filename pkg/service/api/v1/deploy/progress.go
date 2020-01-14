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
	"github.com/gin-gonic/gin"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
)

// Service for wizard progress

// @ID GetWizardProgress
// @Summary Get all of current deploy wizard data
// @Description Get all data, include current progress, cluster and node data. deploying progress or error.
// @Tags wizard
// @Produce application/json
// @Success 200 {object} api.GetWizardResponse
// @Router /api/v1/deploy/wizard/progresses [get]
func GetWizardProgress(c *gin.Context) {

	clusterInfo := getWizardClusterInfo()
	nodes := getWizardNodes()
	networkOptions := getWizardNetworkOptions()
	checkingData := getWizardCheckingData()
	deploymentData := getWizardDeploymentData()

	responseData := api.GetWizardResponse{
		ClusterData:         *clusterInfo,
		NetworkOptions:      *networkOptions,
		NodesData:           *nodes,
		CheckingData:        *checkingData,
		DeploymentData:      *deploymentData,
		CheckResult:         wizard.GetCurrentWizard().GetCheckResult(),
		DeployClusterStatus: convertModelDeployClusterStatusToAPIDeployClusterStatus(wizard.GetCurrentWizard().DeployClusterStatus),
	}

	h.R(c, responseData)
}

// @ID ClearWizard
// @Summary Clear all of current deploy wizard data
// @Description Clear all data, include current progress, cluster and node data. deploying progress or error.
// @Tags wizard
// @Success 204
// @Router /api/v1/deploy/wizard/progresses [delete]
func ClearWizard(c *gin.Context) {

	wizard.ClearCurrentWizardData()
	log.ReqEntry(c).Warn("clear wizard data")

	h.R(c, nil)
}

func getWizardClusterInfo() *api.Cluster {

	wizardData := wizard.GetCurrentWizard()
	clusterInfo := &api.Cluster{
		ShortName:       wizardData.Info.ShortName,
		Name:            wizardData.Info.Name,
		NodePortMinimum: wizardData.Info.NodePortMinimum,
		NodePortMaximum: wizardData.Info.NodePortMaximum,
	}

	switch wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType {
	case wizard.KubeAPIServerConnectTypeFirstMasterIP:
		clusterInfo.KubeAPIServerConnectType = api.KubeAPIServerConnectTypeFirstMasterIP
	case wizard.KubeAPIServerConnectTypeKeepalived:
		clusterInfo.KubeAPIServerConnectType = api.KubeAPIServerConnectTypeKeepalived
		clusterInfo.VIP = wizardData.Info.KubeAPIServerConnection.VIP
		clusterInfo.NetInterfaceName = wizardData.Info.KubeAPIServerConnection.NetInterfaceName
	case wizard.KubeAPIServerConnectTypeLoadBalancer:
		clusterInfo.KubeAPIServerConnectType = api.KubeAPIServerConnectTypeLoadBalancer
		clusterInfo.LoadbalancerIP = wizardData.Info.KubeAPIServerConnection.LoadbalancerIP
		clusterInfo.LoadbalancerPort = wizardData.Info.KubeAPIServerConnection.LoadbalancerPort
	}

	clusterInfo.Labels = make([]api.Label, 0, len(wizardData.Info.Labels))
	for _, label := range wizardData.Info.Labels {
		clusterInfo.Labels = append(clusterInfo.Labels, convertModelLabelToAPILabel(label))
	}

	clusterInfo.Annotations = make([]api.Annotation, 0, len(wizardData.Info.Annotations))
	for _, annotation := range wizardData.Info.Annotations {
		clusterInfo.Annotations = append(clusterInfo.Annotations, convertModelAnnotationToAPIAnnotation(annotation))
	}

	return clusterInfo
}

func getWizardNodes() *[]api.NodeData {

	wizardData := wizard.GetCurrentWizard()
	nodes := new([]api.NodeData)
	*nodes = make([]api.NodeData, 0, len(wizardData.Nodes))

	for _, node := range wizardData.Nodes {

		apiNode := convertModelNodeToAPINode(node)

		*nodes = append(*nodes, *apiNode)
	}

	return nodes
}

func getWizardNetworkOptions() *api.NetworkOptions {
	return wizard.GetCurrentWizard().GetNetworkOptions()
}

func getWizardCheckingData() *[]api.CheckingResultResponseData {

	wizardData := wizard.GetCurrentWizard()
	responseData := new([]api.CheckingResultResponseData)

	*responseData = make([]api.CheckingResultResponseData, 0, len(wizardData.Nodes))

	for _, node := range wizardData.Nodes {

		checkingResult := api.CheckingResultResponseData{
			Name:  node.Name,
			Items: []api.CheckingItem{},
		}

		for _, checkItem := range node.CheckReport.CheckItems {

			checkingResult.Items = append(checkingResult.Items, api.CheckingItem{
				CheckingPoint: checkItem.ItemName,
				Result:        checkItem.CheckResult,
				Error:         convertModelErrorToAPIError(checkItem.Error),
			})
		}

		*responseData = append(*responseData, checkingResult)
	}

	return responseData
}

func getWizardDeploymentData() *[]api.DeploymentResponseData {

	wizardData := wizard.GetCurrentWizard()
	responseData := new([]api.DeploymentResponseData)
	*responseData = make([]api.DeploymentResponseData, 0, 0)

	deployList := make(map[constant.DeployItem][]*api.DeploymentNode)

	for _, node := range wizardData.Nodes {

		for _, report := range node.DeploymentReports {

			deployItem := report.DeployItem
			if _, machineRoleExist := deployList[deployItem]; !machineRoleExist {
				deployList[deployItem] = make([]*api.DeploymentNode, 0, 0)
			}

			deployList[deployItem] = append(deployList[deployItem], &api.DeploymentNode{
				Name:   node.Name,
				Status: convertModelDeployStatusToAPIDeployStatus(report.Status),
				Error:  convertModelErrorToAPIError(report.Error),
			})
		}
	}

	for role, nodes := range deployList {

		nodeList := api.DeploymentResponseData{
			DeployItem: role,
			Nodes:      make([]api.DeploymentNode, 0, len(nodes)),
		}
		for _, node := range nodes {

			nodeList.Nodes = append(nodeList.Nodes, *node)
		}

		*responseData = append(*responseData, nodeList)
	}

	return responseData
}
