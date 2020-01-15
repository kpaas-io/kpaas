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

type AppendAnnotationConfig struct {
	MasterMachine    deployMachine.IMachine
	Logger           *logrus.Entry
	Node             *pb.NodeDeployConfig
	Cluster          *pb.ClusterConfig
	ExecuteLogWriter io.Writer
}

type AppendAnnotation struct {
	config *AppendAnnotationConfig
}

func NewAppendAnnotation(config *AppendAnnotationConfig) *AppendAnnotation {
	return &AppendAnnotation{
		config: config,
	}
}

func (a *AppendAnnotation) append() *pb.Error {

	if len(a.config.Cluster.GetNodeAnnotations()) == 0 {

		a.config.Logger.
			WithFields(logrus.Fields{"node": a.config.Node.GetNode().GetName()}).
			Debug("Not have annotation")
		return nil
	}

	annotations := make([]string, 0, len(a.config.Cluster.GetNodeAnnotations()))
	for annotationKey, annotationValue := range a.config.Cluster.GetNodeAnnotations() {
		annotations = append(annotations, fmt.Sprintf("%s='%s'", annotationKey, annotationValue))
	}

	a.config.Logger.
		WithFields(logrus.Fields{"node": a.config.Node.GetNode().GetName(), "annotations": annotations}).
		Debug("append annotation")

	return operation.NewCommandRunner(a.config.ExecuteLogWriter).RunCommand(
		command.NewKubectlCommand(a.config.MasterMachine, consts.KubeConfigPath, "",
			"annotate", "node", a.config.Node.GetNode().GetName(),
			strings.Join(annotations, " "),
		),
		"Append annotation to node error",                                               // 节点添加Annotation错误
		fmt.Sprintf("append annotation to node: %s", a.config.Node.GetNode().GetName()), // 添加Annotation到 %s 节点
	)
}

func (a *AppendAnnotation) Execute() *pb.Error {

	return a.append()
}
