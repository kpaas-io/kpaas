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

type (
	DeployWorkerActionConfig struct {
		Node            *pb.NodeDeployConfig
		ClusterConfig   *pb.ClusterConfig
		LogFileBasePath string
	}

	DeployWorkerAction struct {
		Base
		config *DeployWorkerActionConfig
	}
)

func NewDeployWorkerAction(config *DeployWorkerActionConfig) (Action, error) {

	if config == nil {
		return nil, fmt.Errorf("action config is nil")
	} else if config.Node == nil {
		return nil, fmt.Errorf("invalid node check config: node is nil")
	}

	actionName := getDeployWorkerActionName(config)
	return &DeployWorkerAction{
		Base: Base{
			Name:              actionName,
			ActionType:        ActionTypeDeployEtcd,
			Status:            ActionPending,
			LogFilePath:       GenActionLogFilePath(config.LogFileBasePath, actionName),
			CreationTimestamp: time.Now(),
		},
		config: config,
	}, nil
}

func getDeployWorkerActionName(config *DeployWorkerActionConfig) string {

	return fmt.Sprintf("worker-%s", config.Node.GetNode().GetName())
}
