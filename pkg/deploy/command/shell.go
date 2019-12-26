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
	"io"
	"strings"
	"time"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/utils"
)

// ShellCommand is a command execute by shell
type ShellCommand struct {
	machine          machine.IMachine
	cmd              string
	args             []string
	executeLogWriter io.Writer
	description      string
}

func NewShellCommand(machine machine.IMachine, cmd string, args ...string) *ShellCommand {
	return &ShellCommand{
		machine: machine,
		cmd:     cmd,
		args:    args,
	}
}

func (c *ShellCommand) WithDescription(desc string) *ShellCommand {
	c.description = desc
	return c
}

func (c *ShellCommand) WithExecuteLogWriter(w io.Writer) *ShellCommand {
	c.executeLogWriter = w
	return c
}

func (c *ShellCommand) Execute() (stdout, stderr []byte, err error) {
	startTime := time.Now()
	stdout, stderr, err = c.machine.Run(c.GetCommand())
	endTime := time.Now()
	if c.executeLogWriter != nil {
		executeLogItem := &utils.ExecuteLogItem{
			StartTime:   startTime,
			EndTime:     endTime,
			Command:     c.cmd + " " + strings.Join(c.args, " "),
			Stdout:      stdout,
			Stderr:      stderr,
			Err:         err,
			Description: c.description,
		}
		utils.WriteExecuteLog(c.executeLogWriter, executeLogItem)
	}
	return
}

func (c *ShellCommand) GetCommand() string {

	cmds := make([]string, 0, len(c.args)+1)
	cmds = append(cmds, c.cmd)
	cmds = append(cmds, c.args...)

	return strings.Join(cmds, " ")
}

func (c *ShellCommand) Exists() (isExist bool, err error) {

	var stderr, stdout []byte
	stdout, stderr, err = c.machine.Run(getCommandExistShell(c.cmd))
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
