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

package task

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
)

func init() {
	RegisterProcessor(TaskTypeInitMaster, new(initMasterProcessor))
}

type initMasterProcessor struct {
}

func (p *initMasterProcessor) SplitTask(t Task) error {
	task, err := p.verifyTask(t)
	if err != nil {
		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: t.GetName(),
	})

	logger.Debug("Start to split init master task into action")

	//task := t.(*InitMasterTask)
	var actions []action.Action
	actionCfg := &action.InitMasterActionConfig{
		CertKey:         task.CertKey,
		Node:            task.Node,
		EtcdNodes:       task.EtcdNodes,
		MasterNodes:     task.MasterNodes,
		ClusterConfig:   task.ClusterConfig,
		LogFileBasePath: task.LogFileDir,
	}
	act, err := action.NewInitMasterAction(actionCfg)
	if err != nil {
		return err
	}
	actions = append(actions, act)
	task.Actions = actions

	logger.Debugf("Finish to split init master task: %d actions", len(actions))

	return nil
}

// Verify if the task is valid.
func (p *initMasterProcessor) verifyTask(t Task) (*InitMasterTask, error) {
	if t == nil {
		return nil, consts.ErrEmptyTask
	}

	task, ok := t.(*InitMasterTask)
	if !ok {
		return nil, fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, t)
	}

	return task, nil
}
