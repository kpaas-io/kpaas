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
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type AppendTaintConfig struct {
	Machine *deployMachine.Machine
	Logger  *logrus.Entry
	Node    *pb.NodeDeployConfig
	Cluster *pb.ClusterConfig
}

type AppendTaint struct {
	operation.BaseOperation
	logger  *logrus.Entry
	node    *pb.NodeDeployConfig
	cluster *pb.ClusterConfig
	machine *deployMachine.Machine
}

func NewAppendTaint(config *AppendTaintConfig) *AppendTaint {
	return &AppendTaint{
		machine: config.Machine,
		logger:  config.Logger,
		node:    config.Node,
		cluster: config.Cluster,
	}
}

func (operation *AppendTaint) append() *pb.Error {

	taints := make([]string, len(operation.node.GetTaints()))
	for _, taint := range operation.node.GetTaints() {
		taints = append(taints, fmt.Sprintf("%s=%s:%s", taint.GetKey(), taint.GetValue(), taint.GetEffect()))
	}

	return RunCommand(
		command.NewKubectlCommand(operation.machine, consts.KubeConfigPath, "",
			"taint", "node", operation.node.GetNode().GetName(),
			strings.Join(taints, " "),
		),
		"Append taint to node error", // 节点添加Taint错误
		fmt.Sprintf("append taint to node: %s", operation.node.GetNode().GetName()), // 添加Taint到 %s 节点
	)
}

func (operation *AppendTaint) Execute() *pb.Error {

	return operation.append()
}
