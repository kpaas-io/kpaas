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
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// NodeInitActionConfig represents the config for a master node init action
type NodeMasterInitActionConfig struct {
	NodeInitConfig  *pb.NodeDeployConfig
	ClusterConfig   *pb.ClusterConfig
	LogFileBasePath string
}

type NodeMasterInitAction struct {
	Base
	sync.RWMutex

	NodeInitConfig *pb.NodeDeployConfig
	ClusterConfig  *pb.ClusterConfig
	InitItems      []*NodeMasterInitItem
}

type NodeMasterInitItemStatus string

type NodeMasterInitItem struct {
	Name        string
	Description string
	Status      NodeInitItemStatus
	Err         *pb.Error
}

// NewNodeInitAction returns a master node init action based on the config.
// User should use this function to create a master node init action.
func NewMasterNodeInitAction(cfg *NodeMasterInitActionConfig) (Action, error) {
	var err error
	if cfg == nil {
		err = fmt.Errorf("action config is nil")
	} else if cfg.NodeInitConfig == nil {
		err = fmt.Errorf("Invalid config: node init config is nil")
	} else if cfg.NodeInitConfig.Node == nil {
		err = fmt.Errorf("Invalid node init config: node is nil")
	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	actionName := getNodMasterInitActionName(cfg)
	return &NodeMasterInitAction{
		Base: Base{
			Name:              actionName,
			ActionType:        ActionTypeNodeInit,
			Status:            ActionPending,
			LogFilePath:       GenActionLogFilePath(cfg.LogFileBasePath, actionName),
			CreationTimestamp: time.Now(),
		},
		NodeInitConfig: cfg.NodeInitConfig,
		ClusterConfig:  cfg.ClusterConfig,
	}, nil
}

// return master node name as the action name, temporarily
func getNodMasterInitActionName(cfg *NodeMasterInitActionConfig) string {
	return cfg.NodeInitConfig.Node.GetName()
}
