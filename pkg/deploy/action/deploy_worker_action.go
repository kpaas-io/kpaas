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

package action

import (
	"fmt"
	"time"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const ActionTypeDeployWorker Type = "DeployWorker"

type DeployWorkerActionConfig struct {
	NodeCfg         *pb.NodeDeployConfig
	ClusterConfig   *pb.ClusterConfig
	MasterNodes     []*pb.Node
	LogFileBasePath string
}

type DeployWorkerAction struct {
	Base
	config *DeployNodeActionConfig
}

func NewDeployWorkerAction(config *DeployNodeActionConfig) (Action, error) {

	if config == nil {
		return nil, fmt.Errorf("action config is nil")
	}
	if config.NodeCfg == nil {
		return nil, fmt.Errorf("invalid action config: NodeCfg is nil")
	}
	if config.NodeCfg.Node == nil {
		return nil, fmt.Errorf("invalid action config: NodeCfg.Node is nil")
	}

	actionName := GenActionName(ActionTypeDeployWorker)
	return &DeployWorkerAction{
		Base: Base{
			Name:              actionName,
			ActionType:        ActionTypeDeployWorker,
			Status:            ActionPending,
			LogFilePath:       GenActionLogFilePath(config.LogFileBasePath, actionName, config.NodeCfg.Node.Name), // /app/deploy/logs/unknown/deploy-worker/{node}-DeployWorker-{randomUint64}.log
			CreationTimestamp: time.Now(),
			Node:              config.NodeCfg.Node,
		},
		config: config,
	}, nil
}
