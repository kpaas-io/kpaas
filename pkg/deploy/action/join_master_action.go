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

const ActionTypeJoinMaster Type = "JoinMaster"

type JoinMasterActionConfig struct {
	Node            *pb.Node
	MasterNodes     []*pb.Node
	ClusterConfig   *pb.ClusterConfig
	LogFileBasePath string
}

type JoinMasterAction struct {
	Base
	MasterNodes   []*pb.Node
	ClusterConfig *pb.ClusterConfig
}

func NewJoinMasterAction(cfg *JoinMasterActionConfig) (Action, error) {
	if cfg == nil {
		return nil, fmt.Errorf("action config is nil")
	}
	if cfg.Node == nil {
		return nil, fmt.Errorf("invalid action config: node is nil")
	}

	actionName := GenActionName(ActionTypeJoinMaster)
	return &JoinMasterAction{
		Base: Base{
			Name:              actionName,
			Node:              cfg.Node,
			ActionType:        ActionTypeJoinMaster,
			Status:            ActionPending,
			LogFilePath:       GenActionLogFilePath(cfg.LogFileBasePath, actionName, cfg.Node.Name),
			CreationTimestamp: time.Now(),
		},
		MasterNodes:   cfg.MasterNodes,
		ClusterConfig: cfg.ClusterConfig,
	}, nil
}
