// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"github.com/gin-gonic/gin"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/service/api/v1/helm"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
)

const (
	defaultCalicoChartPath = "charts/calico"
)

// @ID InstallNetwork
// @Summary install network components
// @Description install network components
// @Tags deploy
// @Success 201 {object} api.SuccessfulOption
// @Failure 404 {object} h.AppErr
// @Router /api/v1/deploy/wizard/networks [post]
// InstallNetwork installs network components in installed kubernetes cluster
func InstallNetwork(c *gin.Context) {
	logger := log.ReqEntry(c)
	var networkOptions protos.NetworkOptions
	err := c.BindJSON(&networkOptions)
	if err != nil {
		logger.WithField("error", err).Warning("failed to parse network options")
		h.E(c, h.EParamsError.WithPayload("failed to get network options"))
		return
	}

	switch networkOptions.NetworkType {
	case "calico":
		appErr := installNetworkCalico(c, networkOptions.CalicoOptions, logger)
		if appErr != nil {
			h.E(c, appErr)
			return
		}
	}

	h.RJsonP(c, api.SuccessfulOption{Success: true})
}

func installNetworkCalico(
	c *gin.Context, options *protos.CalicoOptions, logger *logrus.Entry) error {
	if logger == nil {
		logrus.Warning("empty logger, create one...")
		logger = log.ReqEntry(c)
	}
	logger = logger.WithField("network-type", "calico")

	// get current cluster name.
	wizardData := wizard.GetCurrentWizard()
	// TODO: returns error if cluster is not deployed.
	calicoValues := map[string]interface{}{}

	// fill in calicoValues from options
	if options != nil {
		calicoValues["veth_mtu"] = options.VethMtu

		// TODO: get initial pod range from wizardData
		// calicoValues["ipv4pool_cidr"] = ...

		calicoValues["encap_mode"] = "vxlan"
		calicoValues["vxlan_port"] = options.VxlanPort
		calicoValues["ip_detection.method"] = options.NodeIPDetectionMethod
		if options.NodeIPDetectionMethod == "interface" {
			calicoValues["ip_detection.interface"] = options.NodeIPDetectionInterface
		}
	}

	r := &api.HelmRelease{
		Cluster:   wizardData.Info.Name,
		Namespace: "kube-system",
		Name:      "calico",
		Chart:     defaultCalicoChartPath,
		Values:    calicoValues,
	}
	_, err := helm.RunInstallReleaseAction(c, r)
	return err
}
