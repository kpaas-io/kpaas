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

// DeployTaskConfig represents the config for a deploy task.
type DeployTaskConfig struct {
	NodeConfigs     []*pb.NodeDeployConfig
	ClusterConfig   *pb.ClusterConfig
	LogFileBasePath string
	Priority        int
}

type deployTask struct {
	base
	nodeConfigs   []*pb.NodeDeployConfig
	clusterConfig *pb.ClusterConfig
}

// NewDeployTask returns a deploy task based on the config.
// User should use this function to create a deploy task.
func NewDeployTask(taskName string, taskConfig *DeployTaskConfig) (Task, error) {
	var err error
	if taskConfig == nil {
		err = fmt.Errorf("Invalid task config: nil")

	} else if len(taskConfig.NodeConfigs) == 0 {
		err = fmt.Errorf("Invalid task config: node deploy configs is empty")

	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	task := &deployTask{
		base: base{
			name:              taskName,
			taskType:          TaskTypeNodeCheck,
			status:            TaskPending,
			logFilePath:       GenTaskLogFilePath(taskConfig.LogFileBasePath, taskName),
			creationTimestamp: time.Now(),
			priority:          taskConfig.Priority,
		},
		nodeConfigs:   taskConfig.NodeConfigs,
		clusterConfig: taskConfig.ClusterConfig,
	}

	return task, nil
}
