// Copyright 2020 Shanghai JingDuo Information Technology co., Ltd.
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

	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
)

// SetNetwork set and store network options
// @ID SetNetwork
// @Summary set network options
// @Description set network options
// @Tags network
// @Accept application/json
// @Produce application/json
// @Param networkOptions body protos.NetworkOptions true "options of network components in the cluster"
// @Success 201 {object} api.SuccessfulOption
// @Failure 400 {object} h.AppErr
// @Router /v1/deploy/wizard/networks [post]
func SetNetwork(c *gin.Context) {
	logger := log.ReqEntry(c)
	networkOptions := &protos.NetworkOptions{}
	err := c.BindJSON(networkOptions)
	if err != nil {
		logger.WithError(err).Info("failed to parse network options in request body")
		h.EBindBodyError.WithPayload("failed to parse network options in request body")
		return
	}

	// TODO: validate network options.

	wizardData := wizard.GetCurrentWizard()
	wizardData.SetNetworkOptions(networkOptions)
	h.R(c, &api.SuccessfulOption{Success: true})
}

// GetNetwork get currently stored network options
// @ID GetNetwork
// @Summary get current network options
// @Description get currently stored network options, returns default options if nothing stored.
// @Tags network
// @Produce application/json
// @Success 200 {object} protos.NetworkOptions
// @Router /v1/deploy/wizard/networks [get]
func GetNetwork(c *gin.Context) {
	logger := log.ReqEntry(c)
	wizardData := wizard.GetCurrentWizard()
	logger.WithField("cluster", wizardData.Info.ShortName).Debug("get network options of cluster")
	networkOptions := wizardData.GetNetworkOptions()
	h.R(c, networkOptions)
}
