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

package check

import (
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type ItemEnum int

type OperationsGenerator struct{}

type CheckOperations struct {
	Script     string
	ScriptPath string
}

type CheckAction interface {
	GetOperations(config *pb.NodeCheckConfig) (operation.Operation, error)
	getScript() string
	getScriptPath() string
}

const (
	Docker       ItemEnum = 0
	CPU          ItemEnum = 1
	Kernel       ItemEnum = 2
	Memory       ItemEnum = 3
	Disk         ItemEnum = 4
	Distribution ItemEnum = 5
)

func NewCheckOperations() *OperationsGenerator {
	return &OperationsGenerator{}
}

func (og *OperationsGenerator) CreateOperations(item ItemEnum) CheckAction {
	switch item {
	case Docker:
		return &CheckDockerOperation{}
	case CPU:
		return &CheckCPUOperation{}
	case Kernel:
		return &CheckKernelOperation{}
	case Memory:
		return &CheckMemoryOperation{}
	case Disk:
		return &CheckRootDiskOperation{}
	case Distribution:
		return &CheckDistributionOperation{}
	default:
		return nil
	}
}
