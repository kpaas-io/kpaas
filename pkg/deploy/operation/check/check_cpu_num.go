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
	"bytes"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type CheckCPUOperation struct {
	shellCmd *command.ShellCommand
}

func (ckops *CheckCPUOperation) RunCommands(config *pb.NodeCheckConfig, logChan chan<- *bytes.Buffer) (stdOut, stdErr []byte, err error) {

	itemBuffer := &bytes.Buffer{}

	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, nil, err
	}

	defer m.Close()

	// construct command for check cpu
	ckops.shellCmd = command.NewShellCommand(m, "cat", "/proc/cpuinfo | grep -w 'processor' | awk '{print $NF}' | wc -l").
		WithDescription("检查机器 cpu 数量是否满足最低要求").
		WithExecuteLogWriter(itemBuffer)

	// run commands
	stdOut, stdErr, err = ckops.shellCmd.Execute()

	// write buffer to channel
	logChan <- itemBuffer

	return
}

// check if CPU numbers larger or equal than desired cores
func CheckCPUNums(cpuCore string, desiredCPUCore float64) error {
	err := operation.CheckEntity(cpuCore, desiredCPUCore)
	if err != nil {
		return err
	}
	return nil
}
