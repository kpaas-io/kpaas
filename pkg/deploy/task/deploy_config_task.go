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

const TaskTypeDeployConfig Type = "DeployConfig"

// DeployConfigTaskConfig represents the config for a deploy config task.
type DeployConfigTaskConfig struct {
	NodeConfigs     []*pb.NodeDeployConfig
	MasterNodes     []*pb.Node
	ClusterConfig   *pb.ClusterConfig
	LogFileBasePath string
	Priority        int
	Parent          string
}

type DeployConfigTask struct {
	Base
	NodeConfigs   []*pb.NodeDeployConfig
	MasterNodes   []*pb.Node
	ClusterConfig *pb.ClusterConfig
}

// NewDeployConfigTask returns a deploy config task based on the config.
// User should use this function to create a deploy config task.
func NewDeployConfigTask(taskName string, taskConfig *DeployConfigTaskConfig) (Task, error) {
	var err error
	if taskConfig == nil {
		err = fmt.Errorf("invalid task config: nil")
	} else if len(taskConfig.NodeConfigs) == 0 {
		err = fmt.Errorf("invalid task config: NodeConfigs is empty")
	} else if len(taskConfig.MasterNodes) == 0 {
		err = fmt.Errorf("invalid task config: MasterNodes is empty")
	} else if taskConfig.ClusterConfig == nil {
		err = fmt.Errorf("invalid task config: ClusterConfig is nil")
	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	task := &DeployConfigTask{
		Base: Base{
			Name:                taskName,
			TaskType:            TaskTypeDeployConfig,
			Status:              TaskPending,
			LogFileDir:          GenTaskLogFileDir(taskConfig.LogFileBasePath, taskName),
			CreationTimestamp:   time.Now(),
			Priority:            taskConfig.Priority,
			Parent:              taskConfig.Parent,
			FailureCanBeIgnored: true,
		},
		NodeConfigs:   taskConfig.NodeConfigs,
		MasterNodes:   taskConfig.MasterNodes,
		ClusterConfig: taskConfig.ClusterConfig,
	}

	return task, nil
}
