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

const TaskTypeDeployEtcd Type = "DeployEtcd"

// DeployEtcdTaskConfig represents the config for a deploy etcd task.
type DeployEtcdTaskConfig struct {
	Nodes           []*pb.Node
	LogFileBasePath string
	Priority        int
	Parent          string
}

type DeployEtcdTask struct {
	Base

	Nodes []*pb.Node
}

// NewDeployEtcdTask returns a deploy etcd task based on the config.
// User should use this function to create a deploy etcd task.
func NewDeployEtcdTask(taskName string, taskConfig *DeployEtcdTaskConfig) (Task, error) {
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

	task := &DeployEtcdTask{
		Base: Base{
			Name:              taskName,
			TaskType:          TaskTypeDeployEtcd,
			Status:            TaskPending,
			LogFileDir:        GenTaskLogFileDir(taskConfig.LogFileBasePath, taskName),
			CreationTimestamp: time.Now(),
			Priority:          taskConfig.Priority,
			Parent:            taskConfig.Parent,
		},
		Nodes: taskConfig.Nodes,
	}

	return task, nil
}
