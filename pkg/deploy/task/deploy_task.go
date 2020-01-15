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

	"github.com/kpaas-io/kpaas/pkg/constant"

	"github.com/sirupsen/logrus"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type Operation string
type Priority int

const TaskTypeDeploy Type = "Deploy"

const (
	initOperation   Operation = "initialization"
	deployOperation Operation = "deployment"

	initPriority          Priority = 10
	DeployEtcdPriority    Priority = 20
	DeployMasterPriority  Priority = 30
	DeployWorkerPriority  Priority = 40
	DeployIngressPriority Priority = 50
	ConfigPriority        Priority = 60
)

var (
	Priorities = map[constant.MachineRole]Priority{
		constant.MachineRoleEtcd:    DeployEtcdPriority,
		constant.MachineRoleMaster:  DeployMasterPriority,
		constant.MachineRoleWorker:  DeployWorkerPriority,
		constant.MachineRoleIngress: DeployIngressPriority,
	}
)

// DeployTaskConfig represents the config for a deploy task.
type DeployTaskConfig struct {
	NodeConfigs     []*pb.NodeDeployConfig
	ClusterConfig   *pb.ClusterConfig
	LogFileBasePath string
	Priority        int
}

type DeployTask struct {
	Base
	NodeConfigs   []*pb.NodeDeployConfig
	ClusterConfig *pb.ClusterConfig
}

// NewDeployTask returns a deploy task based on the config.
// User should use this function to create a deploy task.
func NewDeployTask(taskName string, taskConfig *DeployTaskConfig) (Task, error) {
	var err error
	if taskConfig == nil {
		err = fmt.Errorf("invalid task config: nil")

	} else if len(taskConfig.NodeConfigs) == 0 {
		err = fmt.Errorf("invalid task config: node deploy configs is empty")
	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	task := &DeployTask{
		Base: Base{
			Name:              taskName,
			TaskType:          TaskTypeDeploy,
			Status:            TaskPending,
			LogFileDir:        GenTaskLogFileDir(taskConfig.LogFileBasePath, taskName), // /app/deploy/logs/unknown
			CreationTimestamp: time.Now(),
			Priority:          taskConfig.Priority,
		},
		NodeConfigs:   taskConfig.NodeConfigs,
		ClusterConfig: taskConfig.ClusterConfig,
	}

	return task, nil
}
