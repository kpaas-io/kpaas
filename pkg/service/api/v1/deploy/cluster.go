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

// Service for cluster information manage

package deploy

import (
	"github.com/gin-gonic/gin"

	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
	"github.com/kpaas-io/kpaas/pkg/utils/validator"
)

// @ID SetCluster
// @Summary Set Cluster Information
// @Description Store new cluster information
// @Tags cluster
// @Accept application/json
// @Produce application/json
// @Param cluster body api.Cluster true "RequiredFields: shortName, name, kubeAPIServerConnectType"
// @Success 201 {object} api.SuccessfulOption
// @Failure 400 {object} h.AppErr
// @Failure 409 {object} h.AppErr
// @Router /api/v1/deploy/wizard/clusters [post]
func SetCluster(c *gin.Context) {

	requestData, hasError := getClusterRequestData(c)
	if hasError {
		return
	}

	wizardData := wizard.GetCurrentWizard()
	wizardData.Info.Name = requestData.Name
	wizardData.Info.ShortName = requestData.ShortName
	wizardData.Info.KubeAPIServerConnection = wizard.NewKubeAPIServerConnectionData()

	switch requestData.KubeAPIServerConnectType {
	case api.KubeAPIServerConnectTypeFirstMasterIP:
		wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType = wizard.KubeAPIServerConnectTypeFirstMasterIP
	case api.KubeAPIServerConnectTypeKeepalived:
		wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType = wizard.KubeAPIServerConnectTypeKeepalived
		wizardData.Info.KubeAPIServerConnection.VIP = requestData.VIP
		wizardData.Info.KubeAPIServerConnection.NetInterfaceName = requestData.NetInterfaceName
	case api.KubeAPIServerConnectTypeLoadBalancer:
		wizardData.Info.KubeAPIServerConnection.KubeAPIServerConnectType = wizard.KubeAPIServerConnectTypeLoadBalancer
		wizardData.Info.KubeAPIServerConnection.LoadbalancerIP = requestData.LoadbalancerIP
		wizardData.Info.KubeAPIServerConnection.LoadbalancerPort = requestData.LoadbalancerPort
	}
	wizardData.Info.NodePortMinimum = requestData.NodePortMinimum
	wizardData.Info.NodePortMaximum = requestData.NodePortMaximum
	wizardData.Info.Labels = make([]*wizard.Label, 0, len(requestData.Labels))
	for _, label := range requestData.Labels {
		wizardData.Info.Labels = append(wizardData.Info.Labels, &wizard.Label{
			Key:   label.Key,
			Value: label.Value,
		})
	}
	wizardData.Info.Annotations = make([]*wizard.Annotation, 0, len(requestData.Annotations))
	for _, annotation := range requestData.Annotations {
		wizardData.Info.Annotations = append(wizardData.Info.Annotations, &wizard.Annotation{
			Key:   annotation.Key,
			Value: annotation.Value,
		})
	}

	h.R(c, api.SuccessfulOption{Success: true})
}

// @ID GetCluster
// @Summary Get Cluster Information
// @Description Describe cluster information
// @Tags cluster
// @Produce application/json
// @Success 200 {object} api.Cluster
// @Router /api/v1/deploy/wizard/clusters [get]
func GetCluster(c *gin.Context) {

	clusterInfo := getWizardClusterInfo()

	h.R(c, clusterInfo)
}

func getClusterRequestData(c *gin.Context) (requestData *api.Cluster, hasError bool) {

	requestData = new(api.Cluster)
	logger := log.ReqEntry(c)

	if err := validator.Params(c, requestData); err != nil {
		logger.Info(err)
		h.E(c, err)
		return nil, true
	}

	logger.WithField("data", requestData)
	return requestData, false
}
