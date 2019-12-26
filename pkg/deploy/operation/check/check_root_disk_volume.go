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
	rootDiskScript = "/scripts/check_root_disk_volume.sh"
)

type CheckRootDiskOperation struct {
	operation.BaseOperation
	CheckOperations
	Machine machine.IMachine
}

func (ckops *CheckRootDiskOperation) getScript() string {
	ckops.Script = rootDiskScript
	return ckops.Script
}

func (ckops *CheckRootDiskOperation) getScriptPath() string {
	ckops.ScriptPath = checkRemoteScriptPath
	return ckops.ScriptPath
}

func (ckops *CheckRootDiskOperation) GetOperations(config *pb.NodeCheckConfig) (operation.Operation, error) {
	ops := &CheckRootDiskOperation{}
	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}
	ckops.Machine = m

	ops.AddCommands(command.NewShellCommand(m, "df", "-B1 / | awk '/\\//{print $2}'"))
	return ops, nil
}

// close ssh client
func (ckops *CheckRootDiskOperation) CloseSSH() {
	ckops.Machine.Close()
}

// check if root disk volume satisfied with desired disk volume
func CheckRootDiskVolume(rootDiskVolume string, desiredDiskVolume float64) error {
	err := operation.CheckEntity(rootDiskVolume, desiredDiskVolume)
	if err != nil {
		return err
	}
	return nil
}
