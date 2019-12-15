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
	keepalivedScript     = "/scripts/init_deploy_haproxy_keepalived/"
	keepalivedScriptPath = "/tmp"
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
	Machine *machine.Machine
}

func (itOps *InitKeepalivedOperation) getScript() string {
	itOps.Script = routeScript
	return itOps.Script
}

func (itOps *InitKeepalivedOperation) getScriptPath() string {
	itOps.ScriptPath = routeScriptPath
	return itOps.ScriptPath
}

func (itOps *InitKeepalivedOperation) GetOperations(node *pb.Node) (operation.Operation, error) {
	ops := &InitRouteOperation{}
	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, err
	}
	itOps.Machine = m

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

func (itOps *InitKeepalivedOperation) CloseSSH() {
	itOps.Machine.Close()
}
