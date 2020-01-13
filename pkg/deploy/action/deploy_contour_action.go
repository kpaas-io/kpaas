// Copyright 2020 Shanghai JingDuo Information Technology co., Ltd.
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

package action

import (
	"errors"
	"fmt"
	"time"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const ActionTypeDeployContour Type = "DeployContour"

type DeployContourActionConfig struct {
	ClusterConfig   *pb.ClusterConfig
	MasterNodes     []*pb.Node
	LogFileBasePath string
}

type DeployContourAction struct {
	Base
	config *DeployContourActionConfig
}

func NewDeployContourAction(config *DeployContourActionConfig) (Action, error) {

	if config == nil {
		return nil, fmt.Errorf("action config is nil")
	}
	if len(config.MasterNodes) <= 0 {
		return nil, errors.New("master node is empty")
	}

	actionName := GenActionName(ActionTypeDeployContour)
	return &DeployContourAction{
		Base: Base{
			Name:              actionName,
			ActionType:        ActionTypeDeployContour,
			Status:            ActionPending,
			LogFilePath:       GenActionLogFilePath(config.LogFileBasePath, actionName, config.ClusterConfig.ClusterName), // /app/deploy/logs/unknown/deploy-ingress/{clusterName}-DeployContour-{randomUint64}.log
			CreationTimestamp: time.Now(),
			Node:              config.MasterNodes[0],
		},
		config: config,
	}, nil
}
