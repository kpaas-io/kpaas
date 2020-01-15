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
	"bytes"

	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	hostAliasScript = "/scripts/init_change_hostalias.sh"
)

type InitHostaliasOperation struct {
	shellCmd       *command.ShellCommand
	NodeInitAction *operation.NodeInitAction
}

func (itOps *InitHostaliasOperation) RunCommands(node *pb.Node, initAction *operation.NodeInitAction, logChan chan<- *bytes.Buffer) (stdOut, stdErr []byte, err error) {

	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, nil, err
	}

	defer m.Close()

	logBuffer := &bytes.Buffer{}

	itOps.NodeInitAction = initAction

	scriptFile, err := assets.Assets.Open(hostAliasScript)
	if err != nil {
		return nil, nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, operation.InitRemoteScriptPath+hostAliasScript); err != nil {
		return nil, nil, err
	}

	itOps.shellCmd = command.NewShellCommand(m, "bash", operation.InitRemoteScriptPath+hostAliasScript).
		WithDescription("初始化变更 alias").
		WithExecuteLogWriter(logBuffer)

	// run commands
	stdOut, stdErr, err = itOps.shellCmd.Execute()

	// write to log channel
	logChan <- logBuffer

	return
}
