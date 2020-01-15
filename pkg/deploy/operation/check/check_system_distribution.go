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

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	DistributionCentos string = "centos"
	DistributionUbuntu string = "ubuntu"
	DistributionRHEL   string = "rhel"
)

type CheckDistributionOperation struct {
	shellCmd *command.ShellCommand
}

func (ckops *CheckDistributionOperation) RunCommands(config *pb.NodeCheckConfig, logChan chan<- *bytes.Buffer) (stdOut, stdErr []byte, err error) {

	itemBuffer := &bytes.Buffer{}

	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, nil, err
	}

	defer m.Close()

	ckops.shellCmd = command.NewShellCommand(m, "cat", "/etc/*-release | grep -w 'ID' | awk '/ID/{print $1}' | awk -F '=' '{print $2}'").
		WithDescription("检查机器发行版是否满足最低要求").
		WithExecuteLogWriter(itemBuffer)

	// run commands
	stdOut, stdErr, err = ckops.shellCmd.Execute()

	logChan <- itemBuffer

	return
}

// check if system distribution can be supported
func CheckSystemDistribution(disName string) error {
	logger := logrus.WithFields(logrus.Fields{
		"actual_value":  disName,
		"desired_value": fmt.Sprintf("supported distribution: '%v' or '%v' or '%v'", DistributionCentos, DistributionUbuntu, DistributionRHEL),
	})

	if disName == "" {
		logger.Errorf("%v", operation.ErrParaInput)
		return fmt.Errorf("%v, can not be empty", operation.ErrParaInput)
	}

	if disName == DistributionCentos || disName == DistributionUbuntu || disName == DistributionRHEL {
		return nil
	} else {
		logger.Errorf("distribution unclear")
		return fmt.Errorf("unclear distribution, support below: (%v, %v, %v)", DistributionCentos, DistributionUbuntu, DistributionRHEL)
	}
}
