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

}

// @ID GetSSHCertificate
// @Summary Get SSH login keys list
// @Description Get SSH login certificate keys list
// @Tags ssh_certificate
// @Produce application/json
// @Success 200 {object} api.GetSSHCertificateListResponse
// @Router /api/v1/ssh_certificates [get]
func GetCertificateList(c *gin.Context) {

}
