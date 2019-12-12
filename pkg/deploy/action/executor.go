package action

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

import (
	"fmt"
	"sync"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// Executor represents the interface of an action executor.
// Concrete executors implements the logic of actions.
type Executor interface {
	Execute(act Action) error
}

// NewExecutor is a simple factory method to return an action executor based on action type.
func NewExecutor(actionType Type) (Executor, error) {
	var executor Executor
	switch actionType {
	case ActionTypeNodeCheck:
		executor = &nodeCheckExecutor{}
	case ActionTypeNodeInit:
		executor = &nodeInitExecutor{}
	case ActionTypeDeployEtcd:
		executor = &deployEtcdExecutor{}
	case ActionTypeConnectivityCheck:
		executor = &connectivityCheckExecutor{}
	default:
		return nil, fmt.Errorf("%s: %s", consts.MsgActionTypeUnsupported, actionType)
	}

	return executor, nil
}

// ExecuteAction creates and run the executor for an action,
// a *sync.WaitGroup should be passed in.
func ExecuteAction(act Action, wg *sync.WaitGroup) {
	defer wg.Done()

	if act == nil {
		return
	}

	executor, err := NewExecutor(act.GetType())
	if err != nil {
		act.SetStatus(ActionFailed)
		act.SetErr(&pb.Error{
			Reason: consts.MsgActionExecutorCreationFailed,
			Detail: err.Error(),
		})
		return
	}

	act.SetStatus(ActionDoing)

	err = executor.Execute(act)
	if err != nil {
		act.SetStatus(ActionFailed)
		act.SetErr(&pb.Error{
			Reason: consts.MsgActionExecutionFailed,
			Detail: err.Error(),
		})
		return
	}

	act.SetStatus(ActionDone)
}
