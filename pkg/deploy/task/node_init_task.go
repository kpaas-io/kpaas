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

	"github.com/sirupsen/logrus"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// NodeInitTaskConfig represents the config for a node init task
type NodeInitTaskConfig struct {
	Nodes           []*pb.Node
	LogFileBasePath string
	Priority        int
	Parent          string
}

type nodeInitTask struct {
	base
	nodes []*pb.Node
}

// NewNodeInitTask returns a node init task based on the config.
// User should use this function to create a node init task.
func NewNodeInitTask(taskName string, taskConfig *NodeInitTaskConfig) (Task, error) {
	var err error
	if taskConfig == nil {
		err = fmt.Errorf("invalid task config: nil")
	} else if len(taskConfig.Nodes) == 0 {
		err = fmt.Errorf("invalid task config: nodes is empty")
	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	task := &nodeInitTask{
		base: base{
			name:              taskName,
			taskType:          TaskTypeNodeInit,
			status:            TaskPending,
			logFilePath:       GenTaskLogFilePath(taskConfig.LogFileBasePath, taskName),
			creationTimestamp: time.Now(),
			priority:          taskConfig.Priority,
			parent:            taskConfig.Parent,
		},
		nodes: taskConfig.Nodes,
	}

	return task, nil
}
