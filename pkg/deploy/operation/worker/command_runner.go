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
	"time"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type CommandRunner struct {
	executeLogWriter io.Writer
}

func NewCommandRunner(executeLogWriter io.Writer) *CommandRunner {

	return &CommandRunner{executeLogWriter: executeLogWriter}
}

// shellCommand is run at remote command
// errorTitle is pb.Error.Reason when error happened
// doSomeThing is describe what the command done
func (runner *CommandRunner) RunCommand(command command.Command, errorTitle, doSomeThing string) *pb.Error {

	var stdout, stderr []byte
	var err error
	startExecuteTime := time.Now()
	stdout, stderr, err = command.Execute()

	runner.log(consts.DashLine + "\n")
	runner.log(fmt.Sprintf("[start time]: %s\n", startExecuteTime.String()))
	runner.log(fmt.Sprintf("[command]: %s\n", command.GetCommand()))

	if len(stderr) > 0 {
		runner.log(fmt.Sprintf("[stderr]:\n%s\n", string(stderr)))
	}

	if len(stdout) > 0 {
		runner.log(fmt.Sprintf("[stdout]:\n%s\n", string(stdout)))
	}

	if err != nil {
		runner.log(fmt.Sprintf("[error]:\n%s\n", err.Error()))
	}

	runner.log(fmt.Sprintf("[end time]: %s\n\n", time.Now().String()))

	if err != nil {
		return &pb.Error{
			Reason:     errorTitle,                                                                                // {$errorTitle}
			Detail:     fmt.Sprintf("We tried to %s, but command run error, error message: %v", doSomeThing, err), // 我们尝试{$doSomeThing}，命令运行出错了，错误信息： %v
			FixMethods: fixMethodSelfAnalyseIt,                                                                    // 请根据错误提示，并且下载日志进行分析，如果遇到困难，可以提issue给我们
		}
	}

	if len(stderr) > 0 {

		return &pb.Error{
			Reason:     errorTitle,                                                                                              // {$errorTitle}
			Detail:     fmt.Sprintf("We tried to %s, but command return error, error message: %s", doSomeThing, string(stderr)), // 我们尝试{$doSomeThing}，但是命令返回出错了，错误信息： %s
			FixMethods: fixMethodSelfAnalyseIt,                                                                                  // 请根据错误提示，并且下载日志进行分析，如果遇到困难，可以提issue给我们
		}
	}
	return nil
}

func (runner *CommandRunner) log(data string) {

	if runner.executeLogWriter == nil {
		return
	}

	if len(data) > 0 {
		_, _ = io.WriteString(runner.executeLogWriter, data)
	}
}
