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
	NodeInitAction *operation.NodeInitAction
}

func (itOps *InitFireWallOperation) RunCommands(node *pb.Node, initAction *operation.NodeInitAction) (stdOut, stdErr []byte, err error) {

	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, nil, err
	}

	itOps.NodeInitAction = initAction

	// close ssh client if machine is not nil
	if m != nil {
		defer m.Close()
	}

	scriptFile, err := assets.Assets.Open(fireWallScript)
	if err != nil {
		return nil, nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, operation.InitRemoteScriptPath+fireWallScript); err != nil {
		return nil, nil, err
	}

	itOps.AddCommands(command.NewShellCommand(m, "bash", operation.InitRemoteScriptPath+fireWallScript))

	stdOut, stdErr, err = itOps.Do()

	return
}
