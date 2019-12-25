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
	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	fireWallScript = "/scripts/init_change_firewall.sh"
)

type InitFireWallOperation struct {
	operation.BaseOperation
	InitOperations
	Machine        *machine.Machine
	NodeInitAction *operation.NodeInitAction
}

func (itOps *InitFireWallOperation) getScript() string {
	itOps.Script = fireWallScript
	return itOps.Script
}

func (itOps *InitFireWallOperation) getScriptPath() string {
	itOps.ScriptPath = RemoteScriptPath
	return itOps.ScriptPath
}

func (itOps *InitFireWallOperation) GetOperations(node *pb.Node, initAction *operation.NodeInitAction) (operation.Operation, error) {
	ops := &InitFireWallOperation{}
	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, err
	}
	itOps.Machine = m
	itOps.NodeInitAction = initAction

	scriptFile, err := assets.Assets.Open(itOps.getScript())
	if err != nil {
		return nil, err
	}

	if err := m.PutFile(scriptFile, itOps.getScriptPath()+itOps.getScript()); err != nil {
		return nil, err
	}

	ops.AddCommands(command.NewShellCommand(m, "bash", itOps.getScriptPath()+itOps.getScript()))
	return ops, nil
}

func (itOps *InitFireWallOperation) CloseSSH() {
	if itOps.Machine.ExecClient.SSHClient != nil {
		itOps.Machine.Close()
	}
}
