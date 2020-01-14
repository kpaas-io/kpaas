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

const TaskTypeInitMaster Type = "InitMaster"

const (
	InitMasterOperation Operation = "init"
	InitMasterPriority  Priority  = 10
)

type InitMasterTaskConfig struct {
	certKey         string
	operation       Operation
	etcdNodes       []*pb.Node
	MasterNodes     []*pb.Node
	node            *pb.Node
	roles           []string
	clusterConfig   *pb.ClusterConfig
	logFileBasePath string
	Priority        int
	parent          string
}

type InitMasterTask struct {
	Base
	CertKey       string
	Operation     Operation
	EtcdNodes     []*pb.Node
	MasterNodes   []*pb.Node
	Roles         []string
	ClusterConfig *pb.ClusterConfig
	Node          *pb.Node
}

func NewInitMasterTask(taskName string, taskConfig *InitMasterTaskConfig) (Task, error) {
	var err error
	if taskConfig == nil {
		err = fmt.Errorf("invalid task config: nil")
	} else if len(taskConfig.etcdNodes) == 0 {
		err = fmt.Errorf("invalid task config: etcd nodes is empty")
	} else if taskConfig.node == nil {
		err = fmt.Errorf("invalid task config: node is empty")
	} else if taskConfig.clusterConfig == nil {
		err = fmt.Errorf("invalid task config: cluster config is empty")
	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	task := &InitMasterTask{
		Base: Base{
			Name:              taskName,
			TaskType:          TaskTypeInitMaster,
			Status:            TaskPending,
			LogFileDir:        GenTaskLogFileDir(taskConfig.logFileBasePath, taskName),
			CreationTimestamp: time.Now(),
			Priority:          taskConfig.Priority,
			Parent:            taskConfig.parent,
		},
		CertKey:       taskConfig.certKey,
		Node:          taskConfig.node,
		Roles:         taskConfig.roles,
		EtcdNodes:     taskConfig.etcdNodes,
		MasterNodes:   taskConfig.MasterNodes,
		ClusterConfig: taskConfig.clusterConfig,
		Operation:     InitMasterOperation,
	}

	return task, nil
}
