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
	"fmt"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	dockerScript = "/scripts/check_docker_version.sh"
)

type CheckDockerOperation struct {
	operation.BaseOperation
	CheckOperations
	Machine *machine.Machine
}

func (ckops *CheckDockerOperation) getScript() string {
	ckops.Script = dockerScript
	return ckops.Script
}

func (ckops *CheckDockerOperation) getScriptPath() string {
	ckops.ScriptPath = checkRemoteScriptPath
	return ckops.ScriptPath
}

func (ckops *CheckDockerOperation) GetOperations(config *pb.NodeCheckConfig) (operation.Operation, error) {
	ops := &CheckDockerOperation{}
	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}
	ckops.Machine = m

	ops.AddCommands(command.NewShellCommand(m, "docker", fmt.Sprintf(" %v", "version | grep -C1 'Client' | grep -w 'Version:' | awk '{print $2}'")))
	return ops, nil
}

// close ssh client
func (ckops *CheckDockerOperation) CloseSSH() {
	ckops.Machine.Close()
}

// check docker version if version larger or equal than standard version
func CheckDockerVersion(dockerVersion string, standardVersion string, comparedSymbol string) error {
	err := operation.CheckVersion(dockerVersion, standardVersion, comparedSymbol)
	if err != nil {
		return err
	}
	return nil
}
