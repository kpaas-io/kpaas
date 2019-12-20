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

const ActionTypeTestConnection Type = "TestConnection"

// TestConnectionActionConfig represents the config for a test-connection action
type TestConnectionActionConfig struct {
	Node            *pb.Node
	LogFileBasePath string
}

type TestConnectionAction struct {
	Base
}

// NewTestConnectionAction returns a test-connection action based on the config.
// User should use this function to create a test-connection action.
func NewTestConnectionAction(cfg *TestConnectionActionConfig) (Action, error) {
	var err error
	if cfg == nil {
		err = fmt.Errorf("action config is nil")
	} else if cfg.Node == nil {
		err = fmt.Errorf("invalid test connection config: node is nil")
	}
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	actionName := GenActionName(ActionTypeTestConnection)
	return &TestConnectionAction{
		Base: Base{
			Name:              actionName,
			ActionType:        ActionTypeTestConnection,
			Status:            ActionPending,
			LogFilePath:       GenActionLogFilePath(cfg.LogFileBasePath, actionName, cfg.Node.Name),
			CreationTimestamp: time.Now(),
			Node:              cfg.Node,
		},
	}, nil
}
