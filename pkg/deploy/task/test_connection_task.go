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

const TaskTypeTestConnection Type = "TestConnection"

// TestConnectionTaskConfig represents the config for a test-connection task.
type TestConnectionTaskConfig struct {
	Node            *pb.Node
	LogFileBasePath string
	Priority        int
}

type TestConnectionTask struct {
	Base

	Node *pb.Node
}

// NewTestConnectionTask returns a test-connection task based on the config.
// User should use this function to create a test-connection task.
func NewTestConnectionTask(taskName string, taskConfig *TestConnectionTaskConfig) (Task, error) {
	var err error
	if taskConfig == nil {
		err = fmt.Errorf("invalid task config: nil")

	} else if taskConfig.Node == nil {
		err = fmt.Errorf("invalid task config: node is nil")

	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	task := &TestConnectionTask{
		Base: Base{
			Name:              taskName,
			TaskType:          TaskTypeTestConnection,
			Status:            TaskPending,
			LogFilePath:       GenTaskLogFilePath(taskConfig.LogFileBasePath, taskName),
			CreationTimestamp: time.Now(),
			Priority:          taskConfig.Priority,
		},
		Node: taskConfig.Node,
	}

	return task, nil
}
