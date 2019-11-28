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

	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/sshcertificate"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
	"github.com/kpaas-io/kpaas/pkg/utils/validator"
)

// @ID AddSSHCertificate
// @Summary Add SSH login private key
// @Description Add SSH login private key
// @Tags ssh_certificate
// @Accept application/json
// @Produce application/json
// @Param certificate body api.SSHCertificate true "Certificate information"
// @Success 201 {object} api.SuccessfulOption
// @Router /api/v1/ssh_certificates [post]
func AddSSHCertificate(c *gin.Context) {

	requestData, hasError := getSSHCertificateRequestData(c)
	if hasError {
		return
	}

	sshcertificate.AddCertificate(requestData.Name, requestData.Content)

	h.R(c, api.SuccessfulOption{Success: true})
}

// @ID GetSSHCertificate
// @Summary Get SSH login keys list
// @Description Get SSH login certificate keys list
// @Tags ssh_certificate
// @Produce application/json
// @Success 200 {object} api.GetSSHCertificateListResponse
// @Router /api/v1/ssh_certificates [get]
func GetCertificateList(c *gin.Context) {

	h.R(c, api.GetSSHCertificateListResponse{
		Names: sshcertificate.GetNameList(),
	})
}

func getSSHCertificateRequestData(c *gin.Context) (requestData *api.SSHCertificate, hasError bool) {

	requestData = new(api.SSHCertificate)
	logger := log.ReqEntry(c)

	if err := validator.Params(c, requestData); err != nil {
		logger.Info(err)
		h.E(c, err)
		return nil, true
	}

	logger.WithField("data", requestData)
	return requestData, false
}
