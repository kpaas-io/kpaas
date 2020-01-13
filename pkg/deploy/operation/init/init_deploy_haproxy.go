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

	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	HaproxyPort   uint16 = 6443
	haproxyScript        = "/scripts/init_deploy_haproxy_keepalived/setup_kubernetes_high_availability.sh"
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
	NodeInitAction *operation.NodeInitAction
}

func (itOps *InitHaproxyOperation) RunCommands(node *pb.Node, initAction *operation.NodeInitAction) (stdOut, stdErr []byte, err error) {

	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, nil, err
	}

	itOps.NodeInitAction = initAction

	// close ssh client if machine is not nil
	if m != nil {
		defer m.Close()
	}

	if masterIps := itOps.getMastersIP(); len(masterIps) == 0 {
		err = fmt.Errorf("master ip can not be empty")
		return nil, nil, err
	}

	haproxyStr := buildHaproxyStr(itOps.getMastersIP(), HaproxyPort)
	if haproxyStr == "" {
		err = fmt.Errorf("haproxy string can not be built, please check")
		return nil, nil, err
	}

	// put setup.sh to machine
	scriptFile, err := assets.Assets.Open(haproxyScript)
	if err != nil {
		return nil, nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, operation.InitRemoteScriptPath+haproxyScript); err != nil {
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

	itOps.AddCommands(command.NewShellCommand(m, "bash", fmt.Sprintf("%v -u '%v' haproxy run", operation.InitRemoteScriptPath+haproxyScript, haproxyStr)))

	// run commands
	stdOut, stdErr, err = itOps.Do()

	return
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
