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

type DeployWorkerProcessor struct {
}

// Spilt the task into one or more node deploy worker actions
func (processor *DeployWorkerProcessor) SplitTask(task Task) error {
	if err := processor.verifyTask(task); err != nil {
		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: task.GetName(),
	})

	logger.Debug("Start to split deploy worker task")

	deployTask := task.(*deployWorkerTask)

	// split task into actions: will create a action for every node, the action type
	// is ActionTypeDeployWorker

	actions := make([]action.Action, 0, len(deployTask.Nodes))
	for _, node := range deployTask.Nodes {
		actionCfg := &action.DeployWorkerActionConfig{
			Node:            node,
			ClusterConfig:   deployTask.Cluster,
			LogFileBasePath: deployTask.LogFilePath,
		}
		act, err := action.NewDeployWorkerAction(actionCfg)
		if err != nil {
			return err
		}
		actions = append(actions, act)
	}
	deployTask.Actions = actions

	logger.Debugf("Finish to split deploy worker task: %d actions", len(actions))

	return nil
}

// Verify if the task is valid.
func (processor *DeployWorkerProcessor) verifyTask(task Task) error {
	if task == nil {
		return consts.ErrEmptyTask
	}

	deployTask, ok := task.(*deployWorkerTask)
	if !ok {
		return fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, task)
	}

	if len(deployTask.Nodes) == 0 {
		return fmt.Errorf("nodes is empty")
	}

	return nil
}
