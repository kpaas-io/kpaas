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

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
)

func init() {
	RegisterProcessor(TaskTypeDeployMaster, new(deployMasterProcessor))
}

// deployMasterProcessor implements the specific logic to deploy master
type deployMasterProcessor struct {
}

func (p *deployMasterProcessor) SplitTask(t Task) error {
	deployMasterTask, err := p.verifyTask(t)
	if err != nil {
		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: t.GetName(),
	})

	logger.Debug("Start to split deploy task")

	// split task into subtask: init first master, join remain masters
	var subTasks []Task

	for i := range deployMasterTask.Nodes {
		task, err := p.createDeployMasterSubTask(i, deployMasterTask)
		if err != nil {
			err = fmt.Errorf("failed to create init master sub tasks: %s", err)
			logger.Error(err)
			return err
		}
		subTasks = append(subTasks, task)
	}

	deployMasterTask.SubTasks = subTasks
	logger.Debugf("Finish to split deploy master task: %d sub tasks", len(subTasks))

	return nil
}

// Verify if the task is valid.
func (p *deployMasterProcessor) verifyTask(t Task) (*deployMasterTask, error) {
	if t == nil {
		return nil, consts.ErrEmptyTask
	}

	masterTask, ok := t.(*deployMasterTask)
	if !ok {
		return nil, fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, t)
	}

	return masterTask, nil
}

func (p *deployMasterProcessor) createDeployMasterSubTask(index int, parent *deployMasterTask) (Task, error) {
	switch index {
	case 0:
		config := &InitMasterTaskConfig{
			node:            parent.Nodes[index],
			etcdNodes:       parent.EtcdNodes,
			MasterNodes:     parent.Nodes,
			clusterConfig:   parent.ClusterConfig,
			logFileBasePath: parent.GetLogFileDir(),
			Priority:        int(InitMasterPriority),
			parent:          parent.GetName(),
		}

		taskName := "initMaster"
		return NewInitMasterTask(taskName, config)
	default:
		config := &JoinMasterTaskConfig{
			node:            parent.Nodes[index],
			masterNodes:     parent.Nodes,
			clusterConfig:   parent.ClusterConfig,
			logFileBasePath: parent.GetLogFileDir(),
			priority:        int(JoinMasterPriority),
			parent:          parent.GetName(),
		}
		// Use the role name as the task name for now.
		taskName := fmt.Sprintf("%v-%v", "joinMaster", config.node.GetName())
		return NewJoinMasterTask(taskName, config)
	}
}
