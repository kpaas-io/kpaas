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
	"github.com/kpaas-io/kpaas/pkg/service/api/v1/helm"
	"github.com/kpaas-io/kpaas/pkg/service/config"
	clientUtils "github.com/kpaas-io/kpaas/pkg/service/grpcutils/client"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/common"
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

	if wizardData.GetCheckResult() != constant.CheckResultSuccessful {
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

	go deployNetwork()
	go listenDeploymentData()

	h.R(c, api.SuccessfulOption{Success: resp.GetAccepted()})
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
		DeployItems:         *nodeList,
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

	resp, err := client.GetDeployResult(grpcContext, &protos.GetDeployResultRequest{})
	if err != nil {
		logrus.Errorf("call deploy controller error, errorMessage: %v", err)
		return
	}

	wizardData := wizard.GetCurrentWizard()
	wizardData.SetClusterDeploymentStatus(
		computeClusterDeployStatus(resp),
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
			constant.DeployItem(item.DeployItem.GetRole()),
			convertDeployControllerDeployResultToModelDeployResult(item.GetStatus()),
			failureDetail)
	}

	switch wizardData.DeployClusterStatus {
	case wizard.DeployClusterStatusSuccessful, wizard.DeployClusterStatusWorkedButHaveError:

		fetchKubeConfigContent()
	}

}

func computeClusterDeployStatus(resp *protos.GetDeployResultReply) wizard.DeployClusterStatus {

	status := convertDeployControllerDeployClusterStatusToModelDeployClusterStatus(resp.GetStatus())
	switch status {
	case wizard.DeployClusterStatusPending, wizard.DeployClusterStatusRunning, wizard.DeployClusterStatusSuccessful, wizard.DeployClusterStatusDeployServiceUnknown:
		return status
	}

	for _, deployItem := range resp.GetItems() {

		if deployItem.GetDeployItem() == nil {
			continue
		}

		if constant.MachineRole(deployItem.GetDeployItem().GetRole()) != constant.MachineRoleMaster {
			continue
		}

		if convertDeployControllerDeployResultToModelDeployResult(deployItem.GetStatus()) == wizard.DeployStatusSuccessful {
			return wizard.DeployClusterStatusWorkedButHaveError
		}
	}

	return wizard.DeployClusterStatusFailed
}

func fetchKubeConfigContent() {

	wizardData := wizard.GetCurrentWizard()
	client := clientUtils.GetDeployController()
	ctx := context.Background()

	var node *wizard.Node
	for _, iterateNode := range wizardData.Nodes {
		if iterateNode.IsMatchMachineRole(constant.MachineRoleMaster) {
			node = iterateNode
			break
		}
	}

	if node == nil {
		logrus.Errorf("There is no master finish, it is not possible.")
		return
	}

	fetchResponse, err := client.FetchKubeConfig(ctx, &protos.FetchKubeConfigRequest{Node: &protos.Node{
		Name: node.Name,
		Ip:   node.IP,
		Ssh:  convertModelConnectionDataToDeployControllerSSHData(&node.ConnectionData),
	}})

	if err != nil {
		logrus.Errorf("Call gRPC deploy controller error, errorMessage: %v", err)
		return
	}
	kubeConfig := string(fetchResponse.GetKubeConfig())
	wizardData.KubeConfig = &kubeConfig
}

func deployNetwork() {
	wizardData := wizard.GetCurrentWizard()
	networkOptions := wizardData.GetNetworkOptions()

	for _, node := range wizardData.Nodes {
		node.SetDeployResult(
			constant.DeployItemNetwork, wizard.DeployStatusPending, nil)
	}
	logrus.Debugf("waiting for kubernetes to be ready and kubeconfig")
	for {
		// abort deploying of network components if deploying of cluster failed.
		if wizardData.GetDeployClusterStatus() == wizard.DeployClusterStatusFailed {
			logrus.Errorf("deploy cluster failed, abort deploying of network components")
			for _, node := range wizardData.Nodes {
				node.SetDeployResult(
					constant.DeployItemNetwork, wizard.DeployStatusAborted, nil,
				)
			}
			return
		}
		fetchKubeConfigContent()
		if wizardData.KubeConfig != nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	// install calico by helm.
	if networkOptions.NetworkType == api.NetworkTypeCalico {

		err := installCalicoNetwork(networkOptions.CalicoOptions)
		if err != nil {
			for _, node := range wizardData.Nodes {
				node.SetDeployResult(
					constant.DeployItemNetwork, wizard.DeployStatusFailed,
					&common.FailureDetail{
						Reason:     "failed to install calico by helm",
						Detail:     fmt.Sprintf("error %v happened in installing calico by helm", err),
						FixMethods: "check help docs of helm about the error message",
					})
			}
			return
		}
		for _, node := range wizardData.Nodes {
			// TODO: check pod status of network components.
			node.SetDeployResult(
				constant.DeployItemNetwork, wizard.DeployStatusSuccessful, nil)
		}
	}
}

func installCalicoNetwork(options *api.CalicoOptions) error {
	calicoValues := api.HelmValues{}
	if options != nil {
		calicoValues["encap_mode"] = string(options.EncapsulationMode)
		calicoValues["vxlan_port"] = string(options.VxlanPort)
		calicoValues["ipv4_pool"] = options.InitialPodIPs
		calicoValues["ip_detection.method"] = string(options.IPDetectionMethod)
		if options.IPDetectionMethod == api.IPDetectionMethodInterface {
			calicoValues["ip_detection.interface"] = options.IPDetectionInterface
		}
		calicoValues["veth_mtu"] = options.VethMtu
	}
	_, err := helm.RunInstallReleaseAction(nil, &api.HelmRelease{
		Name:      "calico",
		Namespace: "kube-system",
		// TODO: chart path here can be modified and put charts into docker image
		Chart:  "charts/calico",
		Values: calicoValues,
	})
	if err != nil {
		logrus.WithError(err).Errorf("error happened in running helm release action")
		return err
	}
	return nil
}
