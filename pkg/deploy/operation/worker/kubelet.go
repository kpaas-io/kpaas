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

package worker

import (
	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	fixMethodSelfAnalyseIt = "Please follow the error message and download deploy log to analyse it. Please create issues if you find any problem."
)

type InstallKubeletConfig struct {
	Machine *deployMachine.Machine
	Logger  *logrus.Entry
	Node    *pb.NodeDeployConfig
	Cluster *pb.ClusterConfig
}

type InstallKubelet struct {
	operation.BaseOperation
	logger      *logrus.Entry
	node        *pb.NodeDeployConfig
	cluster     *pb.ClusterConfig
	machine     *deployMachine.Machine
	isInstalled bool
}

func NewInstallKubelet(config *InstallKubeletConfig) *InstallKubelet {
	return &InstallKubelet{
		machine: config.Machine,
		logger:  config.Logger,
		node:    config.Node,
		cluster: config.Cluster,
	}
}

func (operation *InstallKubelet) RunKubelet() *pb.Error {

	operation.logger.Info("Start kubelet service")

	if err := operation.runCommand(
		"systemctl restart kubelet",
		"Restart kubelet service error", // 重启kubelet服务错误
		"restart kubelet service",       // 重启kubelet服务
	); err != nil {
		return err
	}

	return nil
}

// shellCommand is run at remote command
// errorTitle is pb.Error.Reason when error happened
// doSomeThing is describe what the command done
func (operation *InstallKubelet) runCommand(shellCommand string, errorTitle string, doSomeThing string) *pb.Error {

	return RunCommand(command.NewShellCommand(operation.machine, shellCommand), errorTitle, doSomeThing)
}

func (operation *InstallKubelet) Execute() *pb.Error {

	if err := operation.RunKubelet(); err != nil {
		return err
	}

	return nil
}
