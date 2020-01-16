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

package action

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/worker"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// With current desigen, a "static" instanc of each type of executor is registered globally,
// so the executor can't hold data to avoid concurrent issues.
type deployNodeExecutor struct {
}

type deployNodeExecutionData struct {
	logger           *logrus.Entry
	machine          deployMachine.IMachine
	executeLogWriter io.Writer
	config           *DeployNodeActionConfig
	action           Action
}

type DeployNodeActionConfig struct {
	NodeCfg         *protos.NodeDeployConfig
	ClusterConfig   *protos.ClusterConfig
	MasterNodes     []*protos.Node
	LogFileBasePath string
}

func (executor *deployNodeExecutor) Deploy(act Action, config *DeployNodeActionConfig) *protos.Error {

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
		consts.LogFieldNode:   config.NodeCfg.GetNode().GetName(),
	})

	logger.Info("start to execute deploy node executor")

	machine, err := executor.connectSSH(config.NodeCfg.GetNode(), logger)
	if err != nil {
		return err
	}
	defer machine.Close()

	executorData := &deployNodeExecutionData{
		action:           act,
		config:           config,
		logger:           logger,
		executeLogWriter: act.GetExecuteLogBuffer(),
		machine:          machine,
	}

	operations := []func(*deployNodeExecutionData) *protos.Error{
		executor.startKubelet,
		executor.joinCluster,
	}

	for _, operation := range operations {
		err := operation(executorData)
		if err != nil {
			return err
		}
	}

	logger.Info("deploy node finished")

	return nil
}

func (executor *deployNodeExecutor) connectSSH(node *protos.Node, logger *logrus.Entry) (deployMachine.IMachine, *protos.Error) {

	logger.Debug("Start to connect ssh")

	machine, err := deployMachine.NewMachine(node)
	if err != nil {
		pbError := &protos.Error{
			Reason:     "Connect ssh error",                                                                         // 连接SSH失败。
			Detail:     fmt.Sprintf("SSH connect to %s(%s) failed , error: %v.", node.GetName(), node.GetIp(), err), // 连接%s(%s)失败，失败原因：%v。
			FixMethods: "Please check node reliability, make SSH service is available.",                             // 请检查节点的可用性，确保SSH服务可用。
		}
		logger.WithField("error", pbError).Error("connect ssh error")
		return nil, pbError
	}

	logger.Debug("ssh connected")
	return machine, nil
}

func (executor *deployNodeExecutor) startKubelet(data *deployNodeExecutionData) *protos.Error {

	data.logger.Debug("Start to install kubelet")

	operation := worker.NewStartKubelet(
		&worker.StartKubeletConfig{
			Machine:          data.machine,
			Node:             data.config.NodeCfg,
			Logger:           data.logger,
			ExecuteLogWriter: data.executeLogWriter,
		},
	)

	if err := operation.Execute(); err != nil {
		data.logger.WithField("error", err).Error("install kubelet error")
		return err
	}

	data.logger.Info("Finish to install kubelet action")
	return nil
}

func (executor *deployNodeExecutor) joinCluster(data *deployNodeExecutionData) *protos.Error {

	data.logger.Debug("Start to join cluster")

	operation := worker.NewJoinCluster(
		&worker.JoinClusterConfig{
			Machine:          data.machine,
			Node:             data.config.NodeCfg,
			Logger:           data.logger,
			Cluster:          data.config.ClusterConfig,
			MasterNodes:      data.config.MasterNodes,
			ExecuteLogWriter: data.executeLogWriter,
		},
	)

	if err := operation.Execute(); err != nil {
		data.logger.WithField("error", err).Error("join cluster error")
		return err
	}

	data.logger.Info("Finish to join cluster action")
	return nil
}
