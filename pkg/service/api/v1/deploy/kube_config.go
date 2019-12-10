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

	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

// @ID DownloadKubeConfig
// @Summary Download the log detail
// @Description Download kubeconfig file
// @Tags kubeconfig
// @Produce text/plain
// @Success 200 {string} string "Kube Config File Content"
// @Failure 400 {object} h.AppErr
// @Failure 404 {object} h.AppErr
// @Router /api/v1/deploy/wizard/kubeconfigs [get]
func DownloadKubeConfig(c *gin.Context) {

	wizardData := wizard.GetCurrentWizard()
	if wizardData.DeployClusterStatus != wizard.DeployClusterStatusSuccessful &&
		wizardData.DeployClusterStatus != wizard.DeployClusterStatusWorkedButHaveError {
		h.E(c, h.EStatusError.WithPayload("Current cluster has not been deployed yet"))
		return
	}

	if wizardData.KubeConfig == nil ||
		*wizardData.KubeConfig == "" {
		h.E(c, h.ENotFound.WithPayload("kubeconfig file has not ready yet, try it later"))

		fetchKubeConfigContent()
		return
	}

	h.R(c, wizardData.KubeConfig)
}
