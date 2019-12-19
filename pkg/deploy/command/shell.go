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

package command

import (
	"strings"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
)

// ShellCommand is a command execute by shell
type ShellCommand struct {
	machine *machine.Machine
	cmd     string
	args    []string
}

func NewShellCommand(machine *machine.Machine, cmd string, args ...string) *ShellCommand {
	return &ShellCommand{
		machine: machine,
		cmd:     cmd,
		args:    args,
	}
}

func (c *ShellCommand) Execute() (stderr, stdout []byte, err error) {

	return c.machine.Run(c.GetCommand())
}

func (c *ShellCommand) GetCommand() string {

	cmds := make([]string, 0, len(c.args)+1)
	cmds = append(cmds, c.cmd)
	cmds = append(cmds, c.args...)

	return strings.Join(cmds, " ")
}

func (c *ShellCommand) Exists() (isExist bool, err error) {

	var stderr, stdout []byte
	stderr, stdout, err = c.machine.Run(getCommandExistShell(c.cmd))
	if err != nil {
		return false, err
	}

	if len(stderr) > 0 {
		return false, nil
	}

	if len(stdout) > 0 {
		return true, nil
	}

	return
}
