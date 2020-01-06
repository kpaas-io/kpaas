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
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type CheckKernelOperation struct {
	operation.BaseOperation
	CheckOperations
	Machine machine.IMachine
}

func (ckops *CheckKernelOperation) GetOperations(config *pb.NodeCheckConfig) (operation.Operation, error) {
	ops := &CheckKernelOperation{}
	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}
	ckops.Machine = m

	ops.AddCommands(command.NewShellCommand(m, "uname", "-r"))
	return ops, nil
}

// close ssh client
func (ckops *CheckKernelOperation) CloseSSH() {
	if ckops.Machine != nil {
		ckops.Machine.Close()
	}
}

// check if kernel version larger or equal than standard version
func CheckKernelVersion(kernelVersion string, standardVersion string, checkStandard string) error {
	err := operation.CheckVersion(kernelVersion, standardVersion, checkStandard)
	if err != nil {
		return err
	}
	return nil
}
