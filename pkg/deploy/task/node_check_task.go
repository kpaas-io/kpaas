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

const TaskTypeNodeCheck Type = "NodeCheck"

// NodeCheckTaskConfig represents the config for a node check task.
type NodeCheckTaskConfig struct {
	NodeConfigs     []*pb.NodeCheckConfig
	NetworkOptions  *pb.NetworkOptions
	LogFileBasePath string
	Priority        int
}

type NodeCheckTask struct {
	Base
	NodeConfigs    []*pb.NodeCheckConfig
	NetworkOptions *pb.NetworkOptions
}

// NewNodeCheckTask returns a node check task based on the config.
// User should use this function to create a node check task.
func NewNodeCheckTask(taskName string, taskConfig *NodeCheckTaskConfig) (Task, error) {
	var err error
	if taskConfig == nil {
		err = fmt.Errorf("invalid task config: nil")

	} else if len(taskConfig.NodeConfigs) == 0 {
		err = fmt.Errorf("invalid task config: node configs is empty")

	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	task := &NodeCheckTask{
		Base: Base{
			Name:              taskName,
			TaskType:          TaskTypeNodeCheck,
			Status:            TaskPending,
			LogFileDir:        GenTaskLogFileDir(taskConfig.LogFileBasePath, taskName),
			CreationTimestamp: time.Now(),
			Priority:          taskConfig.Priority,
		},
		NodeConfigs: taskConfig.NodeConfigs,
	}

	return task, nil
}
