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
	"strings"

	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	portOccupiedScript = "/scripts/check_port_occupied.sh"
)

type CheckPortOccupiedOperation struct {
	operation.BaseOperation
	CheckOperations
	Machine machine.IMachine
}

func (ckops *CheckPortOccupiedOperation) GetOperations(config *pb.NodeCheckConfig) (operation.Operation, error) {
	ops := &CheckPortOccupiedOperation{}
	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}
	ckops.Machine = m

	scriptFile, err := assets.Assets.Open(portOccupiedScript)
	if err != nil {
		return nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, checkRemoteScriptPath+portOccupiedScript); err != nil {
		return nil, err
	}

	// bash script should run as `bash /script/check_port_occupied.sh <role1,role2>` which directly return ports split by comma
	var roles string
	for _, role := range config.Roles {
		roles += role + ","
	}

	if roles == "" {
		return nil, fmt.Errorf("roles can not be empty")
	}
	roles = strings.TrimRight(roles, ",")

	ops.AddCommands(command.NewShellCommand(m, "bash", fmt.Sprintf("%v %v", checkRemoteScriptPath+portOccupiedScript, roles)))

	return ops, nil
}

// close ssh client
func (ckops *CheckPortOccupiedOperation) CloseSSH() {
	if ckops.Machine != nil {
		ckops.Machine.Close()
	}
}

// check if port is occupied
func CheckPortOccupied(portSet string) (string, error) {
	if portSet != "" {
		return portSet, fmt.Errorf("port(s) occupied")
	}

	return "", nil
}
