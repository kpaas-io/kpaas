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
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	keepalivedScript = "/scripts/init_deploy_haproxy_keepalived/setup_kubernetes_high_availability.sh"
)

func CheckKeepalivedParameter(ipAddress string, ethernet string) error {
	logger := logrus.WithFields(logrus.Fields{
		"error_reason": operation.ErrPara,
	})

	if ipAddress == "" || ethernet == "" {
		logger.Errorf("%v", operation.ErrParaEmpty)
		return fmt.Errorf("%v", operation.ErrParaEmpty)
	}

	if ok := operation.CheckIPValid(ipAddress); ok {
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"error_reason": operation.ErrPara,
	}).Errorf("%v", operation.ErrInvalid)
	return fmt.Errorf(operation.ErrInvalid)
}

type InitKeepalivedOperation struct {
	shellCmd       *command.ShellCommand
	NodeInitAction *operation.NodeInitAction
}

func (itOps *InitKeepalivedOperation) RunCommands(node *pb.Node, initAction *operation.NodeInitAction, logChan chan<- *bytes.Buffer) (stdOut, stdErr []byte, err error) {

	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, nil, err
	}

	logBuffer := &bytes.Buffer{}

	itOps.NodeInitAction = initAction

	// close ssh client if machine is not nil
	if m != nil {
		defer m.Close()
	}

	// acquire floating IP for keepalived
	floatingIP := initAction.ClusterConfig.KubeAPIServerConnect.Keepalived.Vip
	if floatingIP == "" {
		err = fmt.Errorf("floating ip can not be empty")
		return nil, nil, err
	}

	// acquire floating ethernet for keepalived
	floatingEthernet := initAction.ClusterConfig.KubeAPIServerConnect.Keepalived.NetInterfaceName
	if floatingEthernet == "" {
		err = fmt.Errorf("floating ethernet can not be empty")
		return nil, nil, err
	}

	// put setup.sh to machine
	scriptFile, err := assets.Assets.Open(keepalivedScript)
	if err != nil {
		return nil, nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, operation.InitRemoteScriptPath+keepalivedScript); err != nil {
		return nil, nil, err
	}

	// put docker.sh to machine
	scriptFile, err = assets.Assets.Open(HaDockerFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, operation.InitRemoteScriptPath+HaDockerFilePath); err != nil {
		return nil, nil, err
	}

	// put lib.sh to machine
	scriptFile, err = assets.Assets.Open(HaLibFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, operation.InitRemoteScriptPath+HaLibFilePath); err != nil {
		return nil, nil, err
	}

	// put systemd.sh to machine
	scriptFile, err = assets.Assets.Open(HaSystemdFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, operation.InitRemoteScriptPath+HaSystemdFilePath); err != nil {
		return nil, nil, err
	}

	itOps.shellCmd = command.NewShellCommand(m, "bash", fmt.Sprintf("%v -n '%v' -i %v keepalived run", operation.InitRemoteScriptPath+keepalivedScript, floatingIP, floatingEthernet)).
		WithDescription("初始化部署 keepalived 工具").
		WithExecuteLogWriter(logBuffer)

	// run commands
	stdOut, stdErr, err = itOps.shellCmd.Execute()

	// write to log channel
	logChan <- logBuffer

	return
}
