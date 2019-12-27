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

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	// we use master IP as keepalived listen IP, later we need use VIP to replace it
	keepalivedEthernet = "eth0"
	keepalivedScript   = "/scripts/init_deploy_haproxy_keepalived/setup_kubernetes_high_availability.sh"
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
	operation.BaseOperation
	InitOperations
	Machine        machine.IMachine
	NodeInitAction *operation.NodeInitAction
}

func (itOps *InitKeepalivedOperation) getScript() string {
	itOps.Script = keepalivedScript
	return itOps.Script
}

func (itOps *InitKeepalivedOperation) getScriptPath() string {
	itOps.ScriptPath = RemoteScriptPath
	return itOps.ScriptPath
}

func (itOps *InitKeepalivedOperation) GetOperations(node *pb.Node, initAction *operation.NodeInitAction) (operation.Operation, error) {
	ops := &InitKeepalivedOperation{}
	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, err
	}
	itOps.Machine = m
	itOps.NodeInitAction = initAction

	masterIP := itOps.getMastersIP()
	if masterIP == "" {
		err = fmt.Errorf("master ip can not be empty")
		return nil, err
	}

	// put setup.sh to machine
	scriptFile, err := assets.Assets.Open(itOps.getScript())
	if err != nil {
		return nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, itOps.getScriptPath()+itOps.getScript()); err != nil {
		return nil, err
	}

	// put docker.sh to machine
	scriptFile, err = assets.Assets.Open(HaDockerFilePath)
	if err != nil {
		return nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, itOps.getScriptPath()+HaDockerFilePath); err != nil {
		return nil, err
	}

	// put lib.sh to machine
	scriptFile, err = assets.Assets.Open(HaLibFilePath)
	if err != nil {
		return nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, itOps.getScriptPath()+HaLibFilePath); err != nil {
		return nil, err
	}

	// put systemd.sh to machine
	scriptFile, err = assets.Assets.Open(HaSystemdFilePath)
	if err != nil {
		return nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, itOps.getScriptPath()+HaSystemdFilePath); err != nil {
		return nil, err
	}

	ops.AddCommands(command.NewShellCommand(m, "bash", fmt.Sprintf("%v -n '%v' -i %v keepalived run", itOps.getScriptPath()+itOps.getScript(), masterIP, keepalivedEthernet)))
	return ops, nil
}

func (itOps *InitKeepalivedOperation) CloseSSH() {
	if itOps.Machine != nil {
		itOps.Machine.Close()
	}
}

// get master IP with config
func (itOps *InitKeepalivedOperation) getMastersIP() string {
	for _, node := range itOps.NodeInitAction.NodesConfig {
		if groupByRole(node.Roles, "master"); true {
			err := CheckHaproxyParameter(node.Node.Ip)
			if err != nil {
				return ""
			}
			return node.Node.Ip
		}
	}
	return ""
}
