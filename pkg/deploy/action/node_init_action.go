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

const (
	nodeInitItemPending = "pending"
	nodeInitItemDoing   = "doing"
	nodeInitItemDone    = "done"
	nodeInitItemFailed  = "failed"
)

// NodeInitActionConfig represents the config for a node init action
type NodeInitActionConfig struct {
	NodeInitConfig  *pb.NodeDeployConfig
	LogFileBasePath string
}

type nodeInitAction struct {
	base
	nodeInitConfig *pb.NodeDeployConfig
	initItems      []*nodeInitItem
	sync.RWMutex
}

type NodeInitItemStatus string

type nodeInitItem struct {
	name        string
	description string
	status      NodeInitItemStatus
	err         *pb.Error
}

// NewNodeInitAction returns a node init action based on the config.
// User should use this function to create a node init action.
func NewNodeInitAction(cfg *NodeInitActionConfig) (Action, error) {
	var err error
	if cfg == nil {
		err = fmt.Errorf("action config is nil")
	} else if cfg.NodeInitConfig == nil {
		err = fmt.Errorf("invalid config: node init config is nil")
	} else if cfg.NodeInitConfig.Node == nil {
		err = fmt.Errorf("invalid node init config: node is nil")
	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	actionName := getNodeInitActionName(cfg)
	return &nodeInitAction{
		base: base{
			name:              actionName,
			actionType:        ActionTypeNodeInit,
			status:            ActionPending,
			logFilePath:       GenActionLogFilePath(cfg.LogFileBasePath, actionName),
			creationTimestamp: time.Now(),
		},
		nodeInitConfig: cfg.NodeInitConfig,
	}, nil
}

// return node name as the action name, temporarily
func getNodeInitActionName(cfg *NodeInitActionConfig) string {
	return cfg.NodeInitConfig.Node.GetName()
}
