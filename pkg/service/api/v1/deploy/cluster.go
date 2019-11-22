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

}
