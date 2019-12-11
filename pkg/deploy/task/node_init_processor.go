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

type nodeInitProcessor struct{}

func (p *nodeInitProcessor) SplitTask(t Task) error {
	if err := p.verifyTask(t); err != nil {
		logrus.Errorf("invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: t.GetName(),
	})

	logger.Debug("Start to split node init task")

	initTask := t.(*nodeInitTask)

	// split task into actions: will create a action for every node, the action type
	// is NodeInitAction
	actions := make([]action.Action, 0, len(initTask.nodeConfigs))
	for _, subConfig := range initTask.nodeConfigs {
		actionCfg := &action.NodeInitActionConfig{
			NodeInitConfig:  subConfig,
			LogFileBasePath: initTask.logFilePath,
		}
		act, err := action.NewNodeInitAction(actionCfg)
		if err != nil {
			return err
		}
		actions = append(actions, act)
	}
	initTask.actions = actions

	logrus.Debugf("Finish to split node init task: %d actions", len(actions))
	return nil
}

// Verify if the task is valid
func (p *nodeInitProcessor) verifyTask(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	nodeInitTask, ok := t.(*nodeInitTask)
	if !ok {
		return fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, t)
	}

	if len(nodeInitTask.nodeConfigs) == 0 {
		return fmt.Errorf("nodeConfig is empty")
	}

	return nil
}
