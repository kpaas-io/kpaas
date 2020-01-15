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

package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type AppendTaintConfig struct {
	MasterMachine    deployMachine.IMachine
	Logger           *logrus.Entry
	Node             *pb.NodeDeployConfig
	Cluster          *pb.ClusterConfig
	ExecuteLogWriter io.Writer
}

type AppendTaint struct {
	config *AppendTaintConfig
}

func NewAppendTaint(config *AppendTaintConfig) *AppendTaint {
	return &AppendTaint{
		config: config,
	}
}

func (a *AppendTaint) append() *pb.Error {

	if len(a.config.Node.GetTaints()) == 0 {
		return nil
	}

	taints := make([]string, 0, len(a.config.Node.GetTaints()))
	for _, taint := range a.config.Node.GetTaints() {
		taints = append(taints, fmt.Sprintf("%s=%s:%s", taint.GetKey(), taint.GetValue(), taint.GetEffect()))
	}

	a.config.Logger.
		WithFields(logrus.Fields{"node": a.config.Node.GetNode().GetName(), "taints": taints}).
		Debug("append taints")

	return operation.NewCommandRunner(a.config.ExecuteLogWriter).RunCommand(
		command.NewKubectlCommand(a.config.MasterMachine, consts.KubeConfigPath, "",
			"taint", "node", a.config.Node.GetNode().GetName(),
			strings.Join(taints, " "),
		),
		"Append taint to node error", // 节点添加Taint错误
		fmt.Sprintf("append taint to node: %s", a.config.Node.GetNode().GetName()), // 添加Taint到 %s 节点
	)
}

func (a *AppendTaint) Execute() *pb.Error {

	return a.append()
}
