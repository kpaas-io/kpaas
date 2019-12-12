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

	"github.com/docker/docker/api/types/versions"
	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type (
	InstallKubeletConfig struct {
		Machine *deployMachine.Machine
		Logger  *logrus.Entry
		Node    *pb.NodeDeployConfig
		Cluster *pb.ClusterConfig
	}
	InstallKubelet struct {
		operation.BaseOperation
		logger      *logrus.Entry
		node        *pb.NodeDeployConfig
		cluster     *pb.ClusterConfig
		machine     *deployMachine.Machine
		isInstalled bool
	}
)

func NewInstallKubelet(config *InstallKubeletConfig) *InstallKubelet {
	return &InstallKubelet{
		machine: config.Machine,
		logger:  config.Logger,
		node:    config.Node,
		cluster: config.Cluster,
	}
}

func (operation *InstallKubelet) ValidateKubelet() *pb.Error {

	operation.logger.Debug("start to check kubelet installation")

	kubelet := command.NewShellCommand(operation.machine, "kubelet",
		"--version",
		"| sed 's/-/_/' ",
		"| awk -Fv '{print $2}'")

	var err error
	var isExist bool
	isExist, err = kubelet.Exists()
	if err != nil {

		return &pb.Error{
			Reason:     "Can not confirm kubelet installation",                                                                               // 无法确认kubelet是否已经安装
			Detail:     fmt.Sprintf("We're try to check the kubelet installation, but execute command result error, error message: %v", err), // 我们尝试确认kubelet是否已经安装，但是执行检查命令时发生错误，错误信息： %v
			FixMethods: "Please follow the error message and download deploy log to analyse it. If any problem can make issue for us.",       // 请根据错误提示，并且下载日志进行分析，如果遇到困难，可以提issue给我们
		}
	}

	// Next step to install kubelet
	if !isExist {
		return nil
	}

	stdout, stderr, err := kubelet.Execute()
	if err != nil {

		return &pb.Error{
			Reason:     "Can not confirm kubelet version",                                                                               // 无法确认kubelet版本号是否符合要求
			Detail:     fmt.Sprintf("We're try to check the kubelet version, but execute command result error, error message: %v", err), // 我们尝试确认kubelet是否符合要求，但是执行检查命令时发生错误，错误信息： %v
			FixMethods: "Please follow the error message and download deploy log to analyse it. If any problem can make issue for us.",  // 请根据错误提示，并且下载日志进行分析，如果遇到困难，可以提issue给我们
		}
	}

	if len(stderr) > 0 {

		return &pb.Error{
			Reason:     "Can not confirm kubelet version",                                                                                          // 无法确认kubelet版本号是否符合要求
			Detail:     fmt.Sprintf("We're try to check the kubelet version, but execute command result error, error message: %s", string(stderr)), // 我们尝试确认kubelet是否符合要求，但是执行检查命令时发生错误，错误信息： %v
			FixMethods: "Please follow the error message and download deploy log to analyse it. If any problem can make issue for us.",             // 请根据错误提示，并且下载日志进行分析，如果遇到困难，可以提issue给我们
		}
	}

	if versions.LessThan(string(stdout), constant.KubeletVersionMinimum) ||
		versions.GreaterThan(string(stdout), constant.KubeletVersionMaximum) {

		return &pb.Error{
			Reason:     "The kubelet version not expected",                                                                                                                                                                                 // kubelet的版本号不满足
			Detail:     fmt.Sprintf("We need the kubelet version between %s and %s, but current kubelet version is %s, please uninstall it and try again", constant.KubeletVersionMinimum, constant.KubeletVersionMaximum, string(stdout)), // 我们需要的Kubelet版本号范围（%s 至 %s），当前检查到版本为： %s，请卸载后再尝试安装
			FixMethods: "Please reinstall the kubelet specific version, or uninstall it later and let us automatic install it.",                                                                                                            // 请尝试重新安装kubelet到指定版本，或者卸载kubelet后，让我们来帮您安装它
		}
	}

	operation.isInstalled = true

	operation.logger.Debug("check kubelet installation completed")
	return nil
}

func (operation *InstallKubelet) InstallKubelet() *pb.Error {

	// TODO Lucky Implement
	return nil
}

func (operation *InstallKubelet) ConfigurateKubelet() *pb.Error {

	// TODO Lucky Implement
	return nil
}

func (operation *InstallKubelet) RunKubelet() *pb.Error {

	// TODO Lucky Implement
	return nil
}

func (operation *InstallKubelet) Execute() *pb.Error {

	if err := operation.ValidateKubelet(); err != nil {
		return err
	}

	if err := operation.InstallKubelet(); err != nil {
		return err
	}

	if err := operation.ConfigurateKubelet(); err != nil {
		return err
	}

	if err := operation.RunKubelet(); err != nil {
		return err
	}

	return nil
}
