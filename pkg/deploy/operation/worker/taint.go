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
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	deployOperation "github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type AppendTaintConfig struct {
	Machine          deployMachine.IMachine
	Logger           *logrus.Entry
	Node             *pb.NodeDeployConfig
	Cluster          *pb.ClusterConfig
	ExecuteLogWriter io.Writer
}

type AppendTaint struct {
	deployOperation.BaseOperation
	config *AppendTaintConfig
}

func NewAppendTaint(config *AppendTaintConfig) *AppendTaint {
	return &AppendTaint{
		config: config,
	}
}

func (operation *AppendTaint) append() *pb.Error {

	if len(operation.config.Node.GetTaints()) == 0 {
		return nil
	}

	taints := make([]string, 0, len(operation.config.Node.GetTaints()))
	for _, taint := range operation.config.Node.GetTaints() {
		taints = append(taints, fmt.Sprintf("%s=%s:%s", taint.GetKey(), taint.GetValue(), taint.GetEffect()))
	}

	operation.config.Logger.
		WithFields(logrus.Fields{"node": operation.config.Node.GetNode().GetName(), "taints": taints}).
		Debug("append taints")

	return deployOperation.NewCommandRunner(operation.config.ExecuteLogWriter).RunCommand(
		command.NewKubectlCommand(operation.config.Machine, consts.KubeConfigPath, "",
			"taint", "node", operation.config.Node.GetNode().GetName(),
			strings.Join(taints, " "),
		),
		"Append taint to node error", // 节点添加Taint错误
		fmt.Sprintf("append taint to node: %s", operation.config.Node.GetNode().GetName()), // 添加Taint到 %s 节点
	)
}

func (operation *AppendTaint) Execute() *pb.Error {

	return operation.append()
}
