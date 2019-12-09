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

package system

import (
	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
)

const (
	kernelScript    = "/scripts/check_kernel_version.sh"
	kernelRemoteDir = "/tmp"
)

type CheckKernelOperation struct {
	operation.BaseOperation
}

func NewCheckKernelOperation(config *pb.NodeCheckConfig) (operation.Operation, error) {
	ops := &CheckKernelOperation{}
	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}

	scriptFile, err := assets.Assets.Open(kernelScript)
	if err != nil {
		return nil, err
	}

	if err := m.PutFile(scriptFile, kernelRemoteDir+kernelScript); err != nil {
		return nil, err
	}

	ops.AddCommands(command.NewShellCommand(m, "bash", kernelRemoteDir+kernelScript, nil))
	return ops, nil
}

// check if kernel version larger or equal than standard version
func CheckKernelVersion(kernelVersion string, standardVersion string, checkStandard string) error {
	err := operation.CheckVersion(kernelVersion, standardVersion, checkStandard)
	if err != nil {
		return err
	}
	return nil
}
