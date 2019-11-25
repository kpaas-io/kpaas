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

package application

import (
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	deploy2 "github.com/kpaas-io/kpaas/pkg/service/api/v1/deploy"
	_ "github.com/kpaas-io/kpaas/pkg/service/swaggerdocs"
)

func (a *app) setRoutes() {
	a.httpHandler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := a.httpHandler.Group("/api/v1")
	wizardGroup := v1.Group("/deploy/wizard")

	wizardGroup.GET("/progresses", deploy2.GetWizardProgress)
	wizardGroup.DELETE("/progresses", deploy2.ClearWizard)

	wizardGroup.POST("/clusters", deploy2.SetCluster)

	wizardGroup.POST("/nodes", deploy2.AddNode)
	wizardGroup.PUT("/nodes/{ip}", deploy2.UpdateNode)
	wizardGroup.DELETE("/nodes/{ip}", deploy2.DeleteNode)

	wizardGroup.POST("/checks", deploy2.CheckNodeList)
	wizardGroup.GET("/checks", deploy2.GetCheckingNodeListResult)

	wizardGroup.POST("/deploys", deploy2.Deploy)
	wizardGroup.GET("/deploys", deploy2.GetDeployReport)

	wizardGroup.GET("/logs/{id}", deploy2.DownloadLog)

	v1.POST("/ssh/tests", deploy2.TestConnectNode)

	v1.POST("/ssh_certificates", deploy2.AddSSHCertificate)
	v1.GET("/ssh_certificates", deploy2.GetCertificateList)
}
