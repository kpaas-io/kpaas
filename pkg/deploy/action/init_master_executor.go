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

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/master"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type initMasterExecutor struct {
}

func (a *initMasterExecutor) Execute(act Action) error {
	action, ok := act.(*InitMasterAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be init master action, but is %T", act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	logger.Debug("Start to init first master action")

	config := &master.InitMasterOperationConfig{
		Logger:        logger,
		Node:          action.Node,
		MasterNodes:   action.MasterNodes,
		EtcdNodes:     action.EtcdNodes,
		ClusterConfig: action.ClusterConfig,
	}

	op, err := master.NewInitMasterOperation(config)
	if err != nil {
		return fmt.Errorf("failed to get init master operation, error: %v", err)
	}

	logger.Debugf("Start to init master on nodes: %s", action.Node.Name)

	action.Status = ActionDone
	if err := op.Do(); err != nil {
		action.Status = ActionFailed
		action.Err = &pb.Error{
			Reason:     err.Error(),
			Detail:     err.Error(),
			FixMethods: "",
		}
	}

	logger.Debug("Finish to execute init master action")
	return nil
}
