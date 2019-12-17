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

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func RunCommand(command command.Command, errorTitle, doSomeThing string) *pb.Error {
	var stderr []byte
	var err error
	_, stderr, err = command.Execute()
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
