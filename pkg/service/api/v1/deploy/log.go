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
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

// @ID DownloadLog
// @Summary Download the log detail
// @Description Download the deployment log details to check the cause of the error
// @Tags log
// @Produce text/plain
// @Param id path int true "Log ID"
// @Success 200 {string} string "Log File Content"
// @Failure 400 {object} h.AppErr
// @Failure 404 {object} h.AppErr
// @Router /api/v1/deploy/wizard/logs/{id} [get]
func DownloadLog(c *gin.Context) {

	logIdString := c.Param("id")
	if len(logIdString) == 0 {
		h.E(c, h.ENotFound.WithPayload("log not exist"))
		return
	}

	logId, err := strconv.ParseUint(logIdString, 10, 64)
	if err != nil {
		h.E(c, h.EParamsError.WithPayload(err))
		return
	}

	content := wizard.GetLog(logId)
	if content == nil {
		h.E(c, h.ENotFound)
		return
	}

	h.R(c, string(content))
}
