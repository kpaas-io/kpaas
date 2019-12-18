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

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// Executor represents the interface of an action executor.
// Concrete executors implements the logic of actions.
type Executor interface {
	Execute(act Action) *pb.Error
}

var _executorRegistry map[Type]Executor

// RegisterExecutor is to register an Executor for an action type
func RegisterExecutor(actionType Type, exec Executor) error {
	if _executorRegistry == nil {
		_executorRegistry = make(map[Type]Executor)
	}
	if exec == nil {
		err := fmt.Errorf("the Executor to be registered is nil")
		logrus.Error(err)
		return err
	}
	if _, ok := _executorRegistry[actionType]; ok {
		err := fmt.Errorf("the Executor for type %v has already been registered", actionType)
		logrus.Error(err)
		return err
	}
	_executorRegistry[actionType] = exec
	return nil
}

// NewExecutor is a simple factory method to return an action executor based on action type.
func NewExecutor(actionType Type) (Executor, error) {
	exec, ok := _executorRegistry[actionType]
	if !ok {
		return nil, fmt.Errorf("%s: %s", consts.MsgActionTypeUnsupported, actionType)
	}

	return exec, nil
}

// ExecuteAction creates and run the executor for an action,
// a *sync.WaitGroup should be passed in.
func ExecuteAction(act Action, wg *sync.WaitGroup) {
	defer wg.Done()

	if act == nil {
		return
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction:     act.GetName(),
		consts.LogFieldActionType: act.GetType(),
	})

	logger.Debug("Start to execute action")

	executor, err := NewExecutor(act.GetType())
	if err != nil {
		act.SetStatus(ActionFailed)
		act.SetErr(&pb.Error{
			Reason: consts.MsgActionExecutorCreationFailed,
			Detail: err.Error(),
		})
		deploy.PBErrLogger(act.GetErr(), logger).Error()
		return
	}

	act.SetStatus(ActionDoing)

	if exeErr := executor.Execute(act); exeErr != nil {
		act.SetStatus(ActionFailed)
		act.SetErr(exeErr)
		deploy.PBErrLogger(act.GetErr(), logger).Error()
		return
	}

	act.SetStatus(ActionDone)
	logger.Debug("Finish to execute action")
}

func errOfTypeMismatched(expected, actual interface{}) *pb.Error {
	return &pb.Error{
		Reason: consts.MsgActionTypeMismatched,
		Detail: fmt.Sprintf(consts.MsgActionTypeMismatchedDetail, expected, actual),
	}
}
