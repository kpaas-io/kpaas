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

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// FetchKubeConfigTaskConfig represents the config for a fetch-kube-config task.
type FetchKubeConfigTaskConfig struct {
	Node            *pb.Node
	LogFileBasePath string
	Priority        int
}

type FetchKubeConfigTask struct {
	base
	node *pb.Node

	// KubeConfig stores the task result: content of kube config file.
	KubeConfig string
}

// NewFetchKubeConfigTask returns a fetch-kube-config task based on the config.
// User should use this function to create a fetch-kube-config task.
func NewFetchKubeConfigTask(taskName string, taskConfig *FetchKubeConfigTaskConfig) (Task, error) {
	if taskName == "" {
		return nil, fmt.Errorf("taskName can't be empty")
	}
	if taskConfig == nil {
		return nil, fmt.Errorf("invalid task config: nil")
	}
	if taskConfig.Node == nil {
		return nil, fmt.Errorf("invalid task config: Node field is nil")

	}

	task := &FetchKubeConfigTask{
		base: base{
			name:              taskName,
			taskType:          TaskTypeFetchKubeConfig,
			status:            TaskPending,
			logFilePath:       GenTaskLogFilePath(taskConfig.LogFileBasePath, taskName),
			creationTimestamp: time.Now(),
			priority:          taskConfig.Priority,
		},
		node: taskConfig.Node,
	}

	return task, nil
}
