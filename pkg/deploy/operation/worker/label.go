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

type AppendLabelConfig struct {
	Machine *deployMachine.Machine
	Logger  *logrus.Entry
	Node    *pb.NodeDeployConfig
	Cluster *pb.ClusterConfig
}

type AppendLabel struct {
	operation.BaseOperation
	logger  *logrus.Entry
	node    *pb.NodeDeployConfig
	cluster *pb.ClusterConfig
	machine *deployMachine.Machine
	labels  map[string]string
}

func NewAppendLabel(config *AppendLabelConfig) *AppendLabel {
	return &AppendLabel{
		machine: config.Machine,
		logger:  config.Logger,
		node:    config.Node,
		cluster: config.Cluster,
		labels:  map[string]string{},
	}
}

func (operation *AppendLabel) computeLabels() {

	for labelKey, labelValue := range operation.cluster.NodeLabels {

		operation.labels[labelKey] = labelValue
	}

	for labelKey, labelValue := range operation.node.GetLabels() {

		operation.labels[labelKey] = labelValue
	}
}

func (operation *AppendLabel) append() *pb.Error {

	labels := make([]string, len(operation.labels))
	for labelKey, labelValue := range operation.labels {
		labels = append(labels, fmt.Sprintf("%s=%s", labelKey, labelValue))
	}

	return RunCommand(
		command.NewKubectlCommand(operation.machine, consts.KubeConfigPath, "",
			"label", "node", operation.node.GetNode().GetName(),
			strings.Join(labels, " "),
		),
		"Append label to node error", // 节点添加Label错误
		fmt.Sprintf("append label to node: %s", operation.node.GetNode().GetName()), // 添加Label到 %s 节点
	)
}

func (operation *AppendLabel) Execute() *pb.Error {

	operation.computeLabels()
	return operation.append()
}
