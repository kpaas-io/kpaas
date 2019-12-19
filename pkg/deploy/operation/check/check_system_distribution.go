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
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	DistributionCentos       string = "centos"
	DistributionUbuntu       string = "ubuntu"
	DistributionRHEL         string = "rhel"
	systemDistributionScript        = "/scripts/check_system_distribution.sh"
)

type CheckDistributionOperation struct {
	operation.BaseOperation
	CheckOperations
	Machine *machine.Machine
}

func (ckops *CheckDistributionOperation) getScript() string {
	ckops.Script = systemDistributionScript
	return ckops.Script
}

func (ckops *CheckDistributionOperation) getScriptPath() string {
	ckops.ScriptPath = checkRemoteScriptPath
	return ckops.ScriptPath
}

func (ckops *CheckDistributionOperation) GetOperations(config *pb.NodeCheckConfig) (operation.Operation, error) {
	ops := &CheckDistributionOperation{}
	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}
	ckops.Machine = m

	ops.AddCommands(command.NewShellCommand(m, "cat", fmt.Sprintf("%v", "/etc/*-release | grep -w 'ID' | awk '/ID/{print $1}' | awk -F '=' '{print $2}'")))
	return ops, nil
}

// close ssh client
func (ckops *CheckDistributionOperation) CloseSSH() {
	ckops.Machine.Close()
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
		logrus.Infof("%v", operation.InfoPassed)
		return nil
	} else {
		logger.Errorf("distribution unclear")
		return fmt.Errorf("unclear distribution, support below: (%v, %v, %v)", DistributionCentos, DistributionUbuntu, DistributionRHEL)
	}
}
