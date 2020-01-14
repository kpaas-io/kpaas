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
	RegisterProcessor(TaskTypeDeployWorker, new(DeployWorkerProcessor))
}

type DeployWorkerProcessor struct {
}

// Spilt the task into one or more node deploy worker actions
func (processor *DeployWorkerProcessor) SplitTask(task Task) error {
	if err := processor.verifyTask(task); err != nil {

		// No need to do something when nodes empty
		if err == consts.ErrEmptyNodes {
			return nil
		}

		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: task.GetName(),
	})

	logger.Debug("Start to split deploy node task")

	deployTask := task.(*deployWorkerTask)

	// split task into actions: will create a action for every node, the action type
	// is ActionTypeDeployWorker

	actions := make([]action.Action, 0, len(deployTask.Config.Nodes))
	for _, node := range deployTask.Config.Nodes {
		actionCfg := &action.DeployNodeActionConfig{
			NodeCfg:         node,
			ClusterConfig:   deployTask.Config.ClusterConfig,
			LogFileBasePath: deployTask.LogFileDir, // /app/deploy/logs/unknown/deploy-worker
			MasterNodes:     deployTask.Config.MasterNodes,
		}
		act, err := action.NewDeployWorkerAction(actionCfg)
		if err != nil {
			return err
		}
		actions = append(actions, act)
	}

	deployTask.Actions = actions

	logger.Debugf("Finish to split deploy node task: %d actions", len(actions))

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

	if len(deployTask.Config.Nodes) == 0 {
		return consts.ErrEmptyNodes
	}

	return nil
}
