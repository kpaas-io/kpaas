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

// Service for nodes information manage

package deploy

import (
	"github.com/gin-gonic/gin"
)

// @ID AddNode
// @Summary Add Node Information
// @Description Add deployment candidate to node list
// @Tags node
// @Accept application/json
// @Produce application/json
// @Param node body api.NodeData true "Node information"
// @Success 201 {object} api.NodeData
// @Failure 400 {object} h.AppErr
// @Failure 409 {object} h.AppErr
// @Router /api/v1/deploy/wizard/nodes [post]
func AddNode(c *gin.Context) {

}

// @ID UpdateNode
// @Summary Update Node Information
// @Description Update a node information which in deployment candidate node list
// @Tags node
// @Accept application/json
// @Produce application/json
// @Param node body api.NodeData true "Node information"
// @Param ip path int true "Node IP Address"
// @Success 200 {object} api.NodeData
// @Failure 400 {object} h.AppErr
// @Failure 404 {object} h.AppErr
// @Failure 409 {object} h.AppErr
// @Router /api/v1/deploy/wizard/nodes/{ip} [put]
func UpdateNode(c *gin.Context) {

}

// @ID DeleteNode
// @Summary Delete a node
// @Description Delete a node from deployment candidate node list
// @Tags node
// @Produce application/json
// @Param ip path int true "Node IP Address"
// @Success 204
// @Failure 400 {object} h.AppErr
// @Failure 404 {object} h.AppErr
// @Failure 409 {object} h.AppErr
// @Router /api/v1/deploy/wizard/nodes/{ip} [delete]
func DeleteNode(c *gin.Context) {

}
