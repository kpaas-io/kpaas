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
	"io"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	fixMethodSelfAnalyseIt = "Please follow the error message and download deploy log to analyse it. Please create issues if you find any problem."
)

type StartKubeletConfig struct {
	Machine          *deployMachine.Machine
	Node             *pb.NodeDeployConfig
	Logger           *logrus.Entry
	ExecuteLogWriter io.Writer
}

type StartKubelet struct {
	operation.BaseOperation
	config *StartKubeletConfig
}

func NewStartKubelet(config *StartKubeletConfig) *StartKubelet {
	return &StartKubelet{
		config: config,
	}
}

func (operation *StartKubelet) RunKubelet() *pb.Error {

	operation.config.Logger.WithField("node", operation.config.Node.GetNode().GetName()).Info("Start kubelet service")

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
func (operation *StartKubelet) runCommand(shellCommand string, errorTitle string, doSomeThing string) *pb.Error {

	return NewCommandRunner(operation.config.ExecuteLogWriter).RunCommand(
		command.NewShellCommand(operation.config.Machine, shellCommand), errorTitle, doSomeThing,
	)
}

func (operation *StartKubelet) Execute() *pb.Error {

	if err := operation.RunKubelet(); err != nil {
		return err
	}

	return nil
}
