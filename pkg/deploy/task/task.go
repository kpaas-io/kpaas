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
	"path/filepath"
	"time"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// Task represents something to do and typically includes one or more actions.
type Task interface {
	GetName() string
	GetType() Type
	GetStatus() Status
	SetStatus(Status)
	GetErr() *pb.Error
	SetErr(*pb.Error)
	GetLogFilePath() string
	SetLogFilePath(string)
	GetActions() []action.Action
	GetCreationTimestamp() time.Time
	// Sub tasks are Task too.
	GetSubTasks() []Task
	// GetPriority returns the priority of the task: smaller value means higher prioirty.
	// A task should wait until all higher priority tasks are done
	GetPriority() int
	// If a task is not a sub task, this will return ""
	GetParent() string
}

// Type represents the type of a task
type Type string

const (
	TaskTypeNodeCheck                Type = "NodeCheck"
	TaskTypeInit                     Type = "init"
	TaskTypeDeploy                   Type = "Deploy"
	TaskTypeDeployEtcd               Type = "DeployEtcd"
	TaskTypeDeployMaster             Type = "DeployMaster"
	TaskTypeDeployWorker             Type = "DeployWorker"
	TaskTypeDeployIngress            Type = "DeployIngess"
	TaskTypeFetchKubeConfig          Type = "FetchKubeConfig"
	TaskTypeCheckNetworkRequirements Type = "CheckNetworkRequirements"
)

// Status represents the status of a task
type Status string

const (
	TaskPending   Status = "Pending"
	TaskSplitting Status = "Splitting"
	TaskSplitted  Status = "Splitted"
	TaskDoing     Status = "Doing"
	TaskDone      Status = "Done" // means success
	TaskFailed    Status = "Failed"
)

type Base struct {
	Name              string
	TaskType          Type
	Actions           []action.Action
	Status            Status
	Err               *pb.Error
	LogFileBasePath   string
	LogFilePath       string
	CreationTimestamp time.Time
	SubTasks          []Task
	Priority          int
	Parent            string
}

func (b *Base) GetName() string {
	return b.Name
}

func (b *Base) GetType() Type {
	return b.TaskType
}

func (b *Base) GetStatus() Status {
	return b.Status
}

func (b *Base) SetStatus(status Status) {
	b.Status = status
}

func (b *Base) GetErr() *pb.Error {
	return b.Err
}

func (b *Base) SetErr(err *pb.Error) {
	b.Err = err
}

func (b *Base) GetLogFilePath() string {
	return b.LogFilePath
}

func (b *Base) SetLogFilePath(path string) {
	b.LogFilePath = path
}

func (b *Base) GetActions() []action.Action {
	return b.Actions
}

func (b *Base) GetCreationTimestamp() time.Time {
	return b.CreationTimestamp
}

func (b *Base) GetSubTasks() []Task {
	return b.SubTasks
}

func (b *Base) GetPriority() int {
	return b.Priority
}

func (b *Base) GetParent() string {
	return b.Parent
}

// GenTaskLogFilePath is a helper to return the log file path based on base path and task name
func GenTaskLogFilePath(basePath, taskName string) string {
	if basePath == "" || taskName == "" {
		return ""
	}
	return filepath.Join(basePath, taskName)
}

// GetAllActions returns all actions of a task, including its direct actions and
// its subtasks' actions recursively.
func GetAllActions(aTask Task) []action.Action {
	var actions []action.Action
	// Collect actions from sub tasks.
	for _, subTask := range aTask.GetSubTasks() {
		actions = append(actions, GetAllActions(subTask)...)
	}

	// Collect direct actions
	actions = append(actions, aTask.GetActions()...)
	return actions
}
