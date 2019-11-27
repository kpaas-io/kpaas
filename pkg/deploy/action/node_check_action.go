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

	"github.com/sirupsen/logrus"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// NodeCheckActionConfig represents the config for a node check action
type NodeCheckActionConfig struct {
	NodeCheckConfig *pb.NodeCheckConfig
	LogFileBasePath string
}

type nodeCheckAction struct {
	base
	nodeCheckConfig *pb.NodeCheckConfig
	checkItems      []*nodeCheckItem
}

type nodeCheckItemStatus string

const (
	nodeCheckItemFailed    nodeCheckItemStatus = "failed"
	nodeCheckItemSucessful nodeCheckItemStatus = "sucessful"
)

type nodeCheckItem struct {
	name        string
	description string
	status      nodeCheckItemStatus
	err         *pb.Error
}

// NewNodeCheckAction returns a node check action based on the config.
// User should use this function to create a node check action.
func NewNodeCheckAction(cfg *NodeCheckActionConfig) (Action, error) {
	var err error
	if cfg == nil {
		err = fmt.Errorf("action config is nil")
	} else if cfg.NodeCheckConfig == nil {
		err = fmt.Errorf("Invalid config: node check config is nil")
	} else if cfg.NodeCheckConfig.Node == nil {
		err = fmt.Errorf("Invalid node check config: node is nil")
	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	actionName := getActionName(cfg)
	return &nodeCheckAction{
		base: base{
			name:              actionName,
			actionType:        ActionTypeNodeCheck,
			status:            ActionPending,
			logFilePath:       GenActionLogFilePath(cfg.LogFileBasePath, actionName),
			creationTimestamp: time.Now(),
		},
		nodeCheckConfig: cfg.NodeCheckConfig,
	}, nil
}

func getActionName(cfg *NodeCheckActionConfig) string {
	// now we used the node name as the the action name, this may be changed in the future.
	return cfg.NodeCheckConfig.Node.GetName()
}
