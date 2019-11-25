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

// @ID CheckNodeList
// @Summary check node list
// @Description Check if the node meets the pre-deployment requirements
// @Tags checking
// @Produce application/json
// @Success 201 {object} api.SuccessfulOption
// @Router /api/v1/deploy/wizard/checks [post]
func CheckNodeList(c *gin.Context) {

}

// @ID GetCheckNodeListResult
// @Summary Get the result of check node
// @Description Get the result of the check node
// @Tags checking
// @Produce application/json
// @Success 200 {object} api.GetCheckingResultResponse
// @Router /api/v1/deploy/wizard/checks [get]
func GetCheckingNodeListResult(c *gin.Context) {

}