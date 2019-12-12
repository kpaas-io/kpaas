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
	"fmt"
	"path/filepath"
	"time"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/utils/idcreator"
)

// Type represents the type of an action
type Type string

const (
	ActionTypeNodeCheck         Type = "NodeCheck"
	ActionTypeDeployEtcd        Type = "DeployEtcd"
	ActionTypeDeployMaster      Type = "DeployMaster"
	ActionTypeDeployWorker      Type = "DeployWorker"
	ActionTypeDeployIngress     Type = "DeployIngress"
	ActionTypeFetchKubeConfig   Type = "FetchKubeConfig"
	ActionTypeConnectivityCheck Type = "ConnectivityCheck"
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
	GetNode() *pb.Node
}

// Base is the basic metadata of an action
type Base struct {
	Name              string
	ActionType        Type
	Status            Status
	Err               *pb.Error
	LogFilePath       string
	CreationTimestamp time.Time
	Node              *pb.Node
}

func (b *Base) GetName() string {
	return b.Name
}

func (b *Base) GetStatus() Status {
	return b.Status
}

func (b *Base) SetStatus(status Status) {
	b.Status = status
}

func (b *Base) GetType() Type {
	return b.ActionType
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

func (b *Base) GetCreationTimestamp() time.Time {
	return b.CreationTimestamp
}

func (b *Base) GetNode() *pb.Node {
	return b.Node
}

// GenActionLogFilePath is a helper to return a file path based on the base path and aciton name
func GenActionLogFilePath(basePath, actionName string) string {
	if basePath == "" || actionName == "" {
		return ""
	}
	return filepath.Join(basePath, actionName)
}

// GenActionName generates a unique action name with the action type as prefix.
func GenActionName(actionType Type) (string, error) {
	str, err := idcreator.NextString()
	if err != nil {
		return "", fmt.Errorf("failed to generate action name: %s", err)
	}
	return fmt.Sprintf("%s-%s", actionType, str), nil
}
