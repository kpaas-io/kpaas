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

const TaskTypeDeployMaster Type = "DeployMaster"

// DeploymasterTaskConfig represents the config for a deploy master task.
type DeployMasterTaskConfig struct {
	etcdNodes       []*pb.Node
	Nodes           []*pb.Node
	ClusterConfig   *pb.ClusterConfig
	LogFileBasePath string
	Priority        int
	Parent          string
}

type deployMasterTask struct {
	Base
	Nodes         []*pb.Node
	EtcdNodes     []*pb.Node
	ClusterConfig *pb.ClusterConfig
}

// NewDeploymasterTask returns a deploy master task based on the config.
// User should use this function to create a deploy master task.
func NewDeployMasterTask(taskName string, taskConfig *DeployMasterTaskConfig) (Task, error) {
	var err error
	if taskConfig == nil {
		err = fmt.Errorf("invalid task config: nil")
	} else if len(taskConfig.Nodes) == 0 {
		err = fmt.Errorf("invalid task config: nodes is empty")
	} else if len(taskConfig.etcdNodes) == 0 {
		err = fmt.Errorf("invalid task config: etcd nodes is empty")
	} else if taskConfig.ClusterConfig == nil {
		err = fmt.Errorf("nil cluster config")
	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	task := &deployMasterTask{
		Base: Base{
			Name:              taskName,
			TaskType:          TaskTypeDeployMaster,
			Status:            TaskPending,
			LogFilePath:       GenTaskLogFilePath(taskConfig.LogFileBasePath, taskName),
			CreationTimestamp: time.Now(),
			Priority:          taskConfig.Priority,
			Parent:            taskConfig.Parent,
		},
		Nodes:         taskConfig.Nodes,
		EtcdNodes:     taskConfig.etcdNodes,
		ClusterConfig: taskConfig.ClusterConfig,
	}

	return task, nil
}
