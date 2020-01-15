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
	"fmt"
	"strings"

	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	portOccupiedScript = "/scripts/check_port_occupied.sh"
)

type CheckPortOccupiedOperation struct {
	shellCmd *command.ShellCommand
}

func (ckops *CheckPortOccupiedOperation) RunCommands(config *pb.NodeCheckConfig, logChan chan<- *bytes.Buffer) (stdOut, stdErr []byte, err error) {

	itemBuffer := &bytes.Buffer{}

	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, nil, err
	}

	defer m.Close()

	scriptFile, err := assets.Assets.Open(portOccupiedScript)
	if err != nil {
		return nil, nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, checkRemoteScriptPath+portOccupiedScript); err != nil {
		return nil, nil, err
	}

	// bash script should run as `bash /script/check_port_occupied.sh <role1,role2>` which directly return ports split by comma
	var roles string
	for _, role := range config.Roles {
		roles += role + ","
	}

	if roles == "" {
		return nil, nil, fmt.Errorf("roles can not be empty")
	}
	roles = strings.TrimRight(roles, ",")

	// construct command for check docker
	ckops.shellCmd = command.NewShellCommand(m, "bash", fmt.Sprintf("%v %v", checkRemoteScriptPath+portOccupiedScript, roles)).
		WithDescription("检查机器端口是否被占用").
		WithExecuteLogWriter(itemBuffer)

	// run commands
	stdOut, stdErr, err = ckops.shellCmd.Execute()

	// write buffer to channel
	logChan <- itemBuffer

	return
}

// check if port is occupied
func CheckPortOccupied(portSet string) (string, error) {
	if portSet != "" {
		return portSet, fmt.Errorf("port(s) occupied")
	}

	return "", nil
}
