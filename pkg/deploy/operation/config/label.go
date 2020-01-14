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

type AppendLabelConfig struct {
	MasterMachine    deployMachine.IMachine
	Logger           *logrus.Entry
	Node             *pb.NodeDeployConfig
	Cluster          *pb.ClusterConfig
	ExecuteLogWriter io.Writer
}

type AppendLabel struct {
	config *AppendLabelConfig
	labels map[string]string
}

func NewAppendLabel(config *AppendLabelConfig) *AppendLabel {
	return &AppendLabel{
		config: config,
		labels: map[string]string{},
	}
}

func (a *AppendLabel) computeLabels() {

	a.computeClusterLabels()
	a.computeNodeLabels()
}

func (a *AppendLabel) computeClusterLabels() {

	if len(a.config.Cluster.GetNodeLabels()) == 0 {

		a.config.Logger.
			WithFields(logrus.Fields{"node": a.config.Node.GetNode().GetName()}).
			Debug("Not have cluster label")
		return
	}

	for labelKey, labelValue := range a.config.Cluster.GetNodeLabels() {

		a.labels[labelKey] = labelValue
	}
}

func (a *AppendLabel) computeNodeLabels() {

	if len(a.config.Node.GetLabels()) == 0 {

		a.config.Logger.
			WithFields(logrus.Fields{"node": a.config.Node.GetNode().GetName()}).
			Debug("Not have node label")
		return
	}

	for labelKey, labelValue := range a.config.Node.GetLabels() {

		a.labels[labelKey] = labelValue
	}
}

func (a *AppendLabel) append() *pb.Error {

	if len(a.labels) == 0 {

		a.config.Logger.
			WithFields(logrus.Fields{"node": a.config.Node.GetNode().GetName()}).
			Info("No label need patch")
		return nil
	}

	labels := make([]string, 0, len(a.labels))
	for labelKey, labelValue := range a.labels {
		labels = append(labels, fmt.Sprintf("%s=%s", labelKey, labelValue))
	}

	a.config.Logger.
		WithFields(logrus.Fields{"node": a.config.Node.GetNode().GetName(), "labels": labels}).
		Debugf("append labels")

	return operation.NewCommandRunner(a.config.ExecuteLogWriter).RunCommand(
		command.NewKubectlCommand(a.config.MasterMachine, consts.KubeConfigPath, "",
			"label", "node", a.config.Node.GetNode().GetName(),
			strings.Join(labels, " "),
		),
		"Append label to node error", // 节点添加Label错误
		fmt.Sprintf("append label to node: %s", a.config.Node.GetNode().GetName()), // 添加Label到 %s 节点
	)
}

func (a *AppendLabel) Execute() *pb.Error {

	a.computeLabels()
	return a.append()
}
