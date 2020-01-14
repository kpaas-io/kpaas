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
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	deployOperation "github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type WriteFileConfig struct {
	Node             deployMachine.IMachine
	Logger           *logrus.Entry
	ExecuteLogWriter io.Writer
	FilePath         string
	FileContent      string
}

type WriteFile struct {
	deployOperation.BaseOperation
	config *WriteFileConfig
}

func NewWriteFile(config *WriteFileConfig) *WriteFile {
	return &WriteFile{
		config: config,
	}
}

func (operation *WriteFile) Execute() *pb.Error {

	if operation.config.FilePath == "" {
		return &pb.Error{
			Reason:     "Cannot write file",                 // 无法写文件
			Detail:     "It's not a specify file path",      // 没有指定文件路径
			FixMethods: consts.MsgFixMethodsPleaseContactUs, // 请联系我们
		}
	}
	if operation.config.FileContent == "" {
		return &pb.Error{
			Reason:     "Cannot write file",                 // 无法写文件
			Detail:     "It's not a specify file content",   // 没有指定文件内容
			FixMethods: consts.MsgFixMethodsPleaseContactUs, // 请联系我们
		}
	}

	operation.config.Logger.
		WithFields(logrus.Fields{"node": operation.config.Node.GetName(), "filePath": operation.config.FilePath}).
		Debugf("write file")

	_, _ = operation.config.ExecuteLogWriter.Write([]byte("Write file\n"))
	_, _ = operation.config.ExecuteLogWriter.Write([]byte(fmt.Sprintf("node: %s\n", operation.config.Node.GetName())))
	_, _ = operation.config.ExecuteLogWriter.Write([]byte(fmt.Sprintf("filePath: %s\n", operation.config.FilePath)))
	_, _ = operation.config.ExecuteLogWriter.Write([]byte(fmt.Sprintf("fileContent:\n%s\n\n", operation.config.FileContent)))

	err := operation.config.Node.PutFile(strings.NewReader(operation.config.FileContent), operation.config.FilePath)
	if err != nil {
		_, _ = operation.config.ExecuteLogWriter.Write([]byte(fmt.Sprintf("Error: %v\n", err)))

		return &pb.Error{
			Reason:     "Cannot write file",                                 // 无法写文件
			Detail:     "When writing file content to %s, we got error: %s", // 我们在写文件（%s）时，发生了些问题。%s
			FixMethods: deployOperation.FixMethodSelfAnalyseIt,              // 请通过部署日志进行排查，如果有其他问题，请联系我们
		}
	}

	return nil
}
