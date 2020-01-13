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

type AppendLabelConfig struct {
	MasterMachine    deployMachine.IMachine
	Logger           *logrus.Entry
	Node             *pb.NodeDeployConfig
	Cluster          *pb.ClusterConfig
	ExecuteLogWriter io.Writer
}

type AppendLabel struct {
	deployOperation.BaseOperation
	config *AppendLabelConfig
	labels map[string]string
}

func NewAppendLabel(config *AppendLabelConfig) *AppendLabel {
	return &AppendLabel{
		config: config,
		labels: map[string]string{},
	}
}

func (operation *AppendLabel) computeLabels() {

	operation.computeClusterLabels()
	operation.computeNodeLabels()
}

func (operation *AppendLabel) computeClusterLabels() {

	if len(operation.config.Cluster.GetNodeLabels()) <= 0 {

		operation.config.Logger.
			WithFields(logrus.Fields{"node": operation.config.Node.GetNode().GetName()}).
			Debug("Not have cluster label")
		return
	}

	for labelKey, labelValue := range operation.config.Cluster.GetNodeLabels() {

		operation.labels[labelKey] = labelValue
	}
}

func (operation *AppendLabel) computeNodeLabels() {

	if len(operation.config.Node.GetLabels()) <= 0 {

		operation.config.Logger.
			WithFields(logrus.Fields{"node": operation.config.Node.GetNode().GetName()}).
			Debug("Not have node label")
		return
	}

	for labelKey, labelValue := range operation.config.Node.GetLabels() {

		operation.labels[labelKey] = labelValue
	}
}

func (operation *AppendLabel) append() *pb.Error {

	if len(operation.labels) <= 0 {

		operation.config.Logger.
			WithFields(logrus.Fields{"node": operation.config.Node.GetNode().GetName()}).
			Info("No label need patch")
		return nil
	}

	labels := make([]string, 0, len(operation.labels))
	for labelKey, labelValue := range operation.labels {
		labels = append(labels, fmt.Sprintf("%s=%s", labelKey, labelValue))
	}

	operation.config.Logger.
		WithFields(logrus.Fields{"node": operation.config.Node.GetNode().GetName(), "labels": labels}).
		Debugf("append labels")

	return deployOperation.NewCommandRunner(operation.config.ExecuteLogWriter).RunCommand(
		command.NewKubectlCommand(operation.config.MasterMachine, consts.KubeConfigPath, "",
			"label", "node", operation.config.Node.GetNode().GetName(),
			strings.Join(labels, " "),
		),
		"Append label to node error", // 节点添加Label错误
		fmt.Sprintf("append label to node: %s", operation.config.Node.GetNode().GetName()), // 添加Label到 %s 节点
	)
}

func (operation *AppendLabel) Execute() *pb.Error {

	operation.computeLabels()
	return operation.append()
}
