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

package action

import (
	"path/filepath"
	"time"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// Type represents the type of an action
type Type string

const (
	ActionTypeNodeCheck     Type = "NodeCheck"
	ActionTypeDeployEtcd    Type = "DeployEtcd"
	ActionTypeDeployMaster  Type = "DeployMaster"
	ActionTypeDeployWorker  Type = "DeployWorker"
	ActionTypeDeployIngress Type = "DeployIngress"
)

// Status represents the status of an action
type Status string

const (
	ActionPending Status = "Pending"
	ActionDoing   Status = "Doing"
	ActionDone    Status = "Done" // means success
	ActionFailed  Status = "Failed"
)

// Action repsents the definition of executable command(s) in a node,
// multiple actions can be executed concurrently.
type Action interface {
	GetName() string
	GetStatus() Status
	SetStatus(Status)
	GetType() Type
	GetErr() *pb.Error
	SetErr(*pb.Error)
	GetLogFilePath() string
	SetLogFilePath(string)
	GetCreationTimestamp() time.Time
}

// base is the basic metadata of an action
type base struct {
	name              string
	actionType        Type
	status            Status
	err               *pb.Error
	logFilePath       string
	creationTimestamp time.Time
}

func (b *base) GetName() string {
	return b.name
}

func (b *base) GetStatus() Status {
	return b.status
}

func (b *base) SetStatus(status Status) {
	b.status = status
}

func (b *base) GetType() Type {
	return b.actionType
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

func (b *base) GetCreationTimestamp() time.Time {
	return b.creationTimestamp
}

// GenActionLogFilePath is a helper to return a file path based on the base path and aciton name
func GenActionLogFilePath(basePath, actionName string) string {
	if basePath == "" || actionName == "" {
		return ""
	}
	return filepath.Join(basePath, actionName)
}
