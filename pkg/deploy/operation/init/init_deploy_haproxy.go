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
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	HaproxyPort       uint16 = 6443
	haproxyScript            = "/scripts/init_deploy_haproxy_keepalived/"
	haproxyScriptPath        = "/scripts/init_deploy_haproxy_keepalived/setup_kubernetes_high_availability.sh"
)

func CheckHaproxyParameter(ipAddresses ...string) error {
	logger := logrus.WithFields(logrus.Fields{
		"error_reason": operation.ErrPara,
	})

	if len(ipAddresses) == 0 {
		logger.Errorf("%v", operation.ErrParaEmpty)
		return fmt.Errorf("%v", operation.ErrParaEmpty)
	}

	for _, ip := range ipAddresses {
		if ok := operation.CheckIPValid(ip); ok {
			continue
		}

		logrus.WithFields(logrus.Fields{
			"error_reason": operation.ErrPara,
		}).Errorf("%v", operation.ErrInvalid)
		return fmt.Errorf("%v", operation.ErrInvalid)
	}

	return nil
}

type InitHaproxyOperation struct {
	operation.BaseOperation
	InitOperations
	Machine        *machine.Machine
	NodeInitAction *operation.NodeInitAction
}

func (itOps *InitHaproxyOperation) getScript() string {
	itOps.Script = haproxyScriptPath
	return itOps.Script
}

func (itOps *InitHaproxyOperation) getScriptPath() string {
	itOps.ScriptPath = RemoteScriptPath
	return itOps.ScriptPath
}

func (itOps *InitHaproxyOperation) GetOperations(node *pb.Node, initAction *operation.NodeInitAction) (operation.Operation, error) {

	ops := &InitHaproxyOperation{}
	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, err
	}
	itOps.Machine = m
	itOps.NodeInitAction = initAction

	if masterIps := itOps.getMastersIP(); len(masterIps) == 0 {
		err = fmt.Errorf("master ip can not be empty")
		return nil, err
	}
	haproxyStr := buildHaproxyStr(itOps.getMastersIP(), HaproxyPort)
	if haproxyStr == "" {
		err = fmt.Errorf("haproxy string can not be built, please check")
		return nil, err
	}

	err = m.PutDir(haproxyScript, RemoteScriptPath, deploy.AllFilesNeeded)
	if err != nil {
		return nil, err
	}

	ops.AddCommands(command.NewShellCommand(m, "bash", fmt.Sprintf("%v -u '%v' haproxy run", itOps.getScriptPath()+itOps.getScript(), haproxyStr)))
	return ops, nil
}

func (itOps *InitHaproxyOperation) CloseSSH() {
	if itOps.Machine == nil {
		return
	}
	itOps.Machine.Close()
}

// construct haproxy parameter
func buildHaproxyStr(masterIps []string, port uint16) string {
	haproxyStr := ""
	if len(masterIps) == 0 {
		return ""
	}
	for _, ip := range masterIps {
		haproxyStr += haproxyStr + fmt.Sprintf("%v:%v ", ip, port)
	}
	haproxyStr = strings.TrimSpace(haproxyStr)
	return haproxyStr
}

// get master IP with config
func (itOps *InitHaproxyOperation) getMastersIP() []string {
	masterIps := []string{}
	for _, node := range itOps.NodeInitAction.NodesConfig {
		if groupByRole(node.Roles, "master"); true {
			err := CheckHaproxyParameter(node.Node.Ip)
			if err != nil {
				return []string{}
			}
			masterIps = append(masterIps, node.Node.Ip)
		}
	}
	if len(masterIps) < 3 {
		return []string{}
	}
	return masterIps
}
