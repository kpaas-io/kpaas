// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
)

type Command interface {
	Execute() ([]byte, []byte, error)
}

// ShellCommand is a command execute by shell
type ShellCommand struct {
	machine *machine.Machine
	cmd     string
	subCmd  string
	options map[string]string
}

func NewShellCommand(machine *machine.Machine, cmd string, subCmd string, options map[string]string) *ShellCommand {
	return &ShellCommand{
		machine: machine,
		cmd:     cmd,
		subCmd:  subCmd,
		options: options,
	}
}

func (c *ShellCommand) Execute() (stderr, stdout []byte, err error) {
	cmds := []string{
		c.cmd,
		c.subCmd,
	}
	for k, v := range c.options {
		cmds = append(cmds, k+" "+v)
	}

	cmd := strings.Join(cmds, " ")

	stderr, stdout, err = c.machine.Run(cmd)

	return
}

type LocalShellCommand struct {
	cmd     string
	subCmds []string
	options map[string]string
}

func NewLocalShellCommand(cmd string, subCmds []string, options map[string]string) *LocalShellCommand {
	return &LocalShellCommand{
		cmd:     cmd,
		subCmds: subCmds,
		options: options,
	}
}

func (c *LocalShellCommand) Execute() (stderr, stdout []byte, err error) {
	args := make([]string, 0)
	args = append(args, c.subCmds...)
	for k, v := range c.options {
		arg := fmt.Sprintf("%v=%v", k, v)
		args = append(args, arg)
	}

	cmd := exec.Command(c.cmd, args...)
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()

	if stderr, err = ioutil.ReadAll(errReader); err != nil {
		return
	}

	if stdout, err = ioutil.ReadAll(outReader); err != nil {
		return
	}

	return
}
