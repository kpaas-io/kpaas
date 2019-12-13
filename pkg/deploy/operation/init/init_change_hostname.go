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

package init

import (
	"fmt"

	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	hostNameScript     = "/scripts/init_change_hostname.sh"
	hostNameScriptPath = "/tmp"
)

type InitHostNameOperation struct {
	operation.BaseOperation
	InitOperations
	Machine *machine.Machine
}

func (itOps *InitHostNameOperation) getScript() string {
	itOps.Script = hostNameScript
	return itOps.Script
}

func (itOps *InitHostNameOperation) getScriptPath() string {
	itOps.ScriptPath = hostNameScriptPath
	return itOps.ScriptPath
}

func (itOps *InitHostNameOperation) GetOperations(node *pb.Node) (operation.Operation, error) {
	ops := &InitHostNameOperation{}
	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, err
	}
	itOps.Machine = m

	scriptFile, err := assets.Assets.Open(itOps.getScript())
	if err != nil {
		return nil, err
	}

	if err := m.PutFile(scriptFile, itOps.getScriptPath()+itOps.getScript()); err != nil {
		return nil, err
	}

	currentName := node.Name

	ops.AddCommands(command.NewShellCommand(m, "bash", fmt.Sprintf("%v %v", itOps.getScriptPath()+itOps.getScript(), currentName)))
	return ops, nil
}

func (itOps *InitHostNameOperation) CloseSSH() {
	itOps.Machine.Close()
}
