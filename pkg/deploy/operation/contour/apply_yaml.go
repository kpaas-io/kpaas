// Copyright 2020 Shanghai JingDuo Information Technology co., Ltd.
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

package contour

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	deployOperation "github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type ApplyYAMLConfig struct {
	Node             deployMachine.IMachine
	Logger           *logrus.Entry
	ExecuteLogWriter io.Writer
	FilePath         string
}

type ApplyYAML struct {
	deployOperation.BaseOperation
	config *ApplyYAMLConfig
}

func NewApplyYAML(config *ApplyYAMLConfig) *ApplyYAML {
	return &ApplyYAML{
		config: config,
	}
}

func (operation *ApplyYAML) Execute() *pb.Error {

	if operation.config.FilePath == "" {
		return &pb.Error{
			Reason:     "Cannot apply yaml",               // 无法应用YAML
			Detail:     "It's not specify yaml file path", // 没有指定YAML文件路径
			FixMethods: "Please contact us",               // 请联系我们
		}
	}

	operation.config.Logger.
		WithFields(logrus.Fields{"node": operation.config.Node.GetName(), "filePath": operation.config.FilePath}).
		Debugf("apply yaml")

	return deployOperation.NewCommandRunner(operation.config.ExecuteLogWriter).RunCommand(
		command.NewKubectlCommand(operation.config.Node, consts.KubeConfigPath, "",
			"apply", "-f", operation.config.FilePath,
		),
		"Apply YAML error", // 应用YAML错误
		fmt.Sprintf("Apply YAML %s", operation.config.FilePath), // 应用 %s YAML
	)
}
