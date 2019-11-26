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
}

// Type represents the type of a task
type Type string

const (
	TaskTypeNodeCheck Type = "NodeCheck"
	TaskTypeDeploy    Type = "Deploy"
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

type base struct {
	name              string
	taskType          Type
	actions           []action.Action
	status            Status
	err               *pb.Error
	logFilePath       string
	creationTimestamp time.Time
}

func (b *base) GetName() string {
	return b.name
}

func (b *base) GetType() Type {
	return b.taskType
}

func (b *base) GetStatus() Status {
	return b.status
}

func (b *base) SetStatus(status Status) {
	b.status = status
}

func (b *base) GetErr() *pb.Error {
	return b.err
}

func (b *base) SetErr(err *pb.Error) {
	b.err = err
}

func (b *base) GetLogFilePath() string {
	return b.logFilePath
}

func (b *base) SetLogFilePath(path string) {
	b.logFilePath = path
}

func (b *base) GetActions() []action.Action {
	return b.actions
}

func (b *base) GetCreationTimestamp() time.Time {
	return b.creationTimestamp
}

// GenTaskLogFilePath is a helper to return the log file path based on base path and task name
func GenTaskLogFilePath(basePath, taskName string) string {
	if basePath == "" || taskName == "" {
		return ""
	}
	return filepath.Join(basePath, taskName)
}
