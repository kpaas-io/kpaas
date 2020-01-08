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
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type ItemEnum string

type OperationsGenerator struct {
	Node *pb.Node
}

type CheckOperations struct {
	Script     string
	ScriptPath string
	Machine    machine.IMachine
}

type CheckAction interface {
	CreateCommandAndRun(config *pb.NodeCheckConfig) ([]byte, []byte, error)
}

const (
	checkRemoteScriptPath          = "/tmp"
	Docker                ItemEnum = "docker"
	CPU                   ItemEnum = "cpu"
	Kernel                ItemEnum = "kernel"
	Memory                ItemEnum = "memory"
	Disk                  ItemEnum = "disk"
	Distribution          ItemEnum = "distribution"
	SystemPreference      ItemEnum = "system-preference"
	SystemManager         ItemEnum = "system-manager"
	PortOccupied          ItemEnum = "port-occupied"
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
	case SystemPreference:
		return &CheckSysPrefOperation{}
	case SystemManager:
		return &CheckSystemManagerOperation{}
	case PortOccupied:
		return &CheckPortOccupiedOperation{}
	default:
		return nil
	}
}
