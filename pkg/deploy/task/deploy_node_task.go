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
	"time"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const TaskTypeDeployNode Type = "DeployNode"

type DeployNodeTaskConfig struct {
	BaseTaskConfig
	MasterNodes   []*protos.Node
	Nodes         []*protos.NodeDeployConfig
	ClusterConfig *protos.ClusterConfig
	PostActions   []action.Action
}

type deployNodeTask struct {
	Base
	Config *DeployNodeTaskConfig
}

// NewDeployNodeTask returns a deploy k8s worker task based on the config.
// User should use this function to create a deploy worker task.
func NewDeployNodeTask(taskName string, taskConfig *DeployNodeTaskConfig) (Task, error) {

	if taskConfig == nil {

		return nil, fmt.Errorf("invalid task config: nil")
	}

	if len(taskConfig.Nodes) == 0 {

		return nil, fmt.Errorf("invalid task config: nodes is empty")
	}

	task := &deployNodeTask{
		Base: Base{
			Name:                taskName,
			TaskType:            TaskTypeDeployNode,
			Status:              TaskPending,
			LogFileDir:          GenTaskLogFileDir(taskConfig.LogFileBasePath, taskName), // /app/deploy/logs/unknown/deploy-{role}
			CreationTimestamp:   time.Now(),
			Priority:            taskConfig.Priority,
			Parent:              taskConfig.Parent,
			FailureCanBeIgnored: true,
		},
		Config: taskConfig,
	}

	return task, nil
}
