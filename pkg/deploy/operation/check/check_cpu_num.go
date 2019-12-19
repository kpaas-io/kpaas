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

const (
	cpuScript = "/scripts/check_cpu_num.sh"
)

type CheckCPUOperation struct {
	operation.BaseOperation
	CheckOperations
	Machine *machine.Machine
}

func (ckops *CheckCPUOperation) getScript() string {
	ckops.Script = cpuScript
	return ckops.Script
}

func (ckops *CheckCPUOperation) getScriptPath() string {
	ckops.ScriptPath = checkRemoteScriptPath
	return ckops.ScriptPath
}

func (ckops *CheckCPUOperation) GetOperations(config *pb.NodeCheckConfig) (operation.Operation, error) {
	ops := &CheckCPUOperation{}
	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}
	ckops.Machine = m

	ops.AddCommands(command.NewShellCommand(m, "cat", "/proc/cpuinfo | grep -w 'processor' | awk '{print $NF}' | wc -l"))
	return ops, nil
}

// close ssh client
func (ckops *CheckCPUOperation) CloseSSH() {
	ckops.Machine.Close()
}

// check if CPU numbers larger or equal than desired cores
func CheckCPUNums(cpuCore string, desiredCPUCore float64) error {
	err := operation.CheckEntity(cpuCore, desiredCPUCore)
	if err != nil {
		return err
	}
	return nil
}
