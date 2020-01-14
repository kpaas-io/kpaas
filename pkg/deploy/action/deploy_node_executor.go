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

type deployNodeExecutor struct {
	logger           *logrus.Entry
	machine          deployMachine.IMachine
	masterMachine    deployMachine.IMachine
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

	executor.action = act
	executor.config = config

	executor.initLogger()
	executor.initExecuteLogWriter()

	executor.logger.Info("start to execute deploy node executor")

	if err := executor.connectSSH(); err != nil {

		return err
	}
	defer executor.disconnectSSH()

	if err := executor.connectMasterNode(); err != nil {
		return err
	}
	defer executor.disconnectMasterNode()

	operations := []func() *protos.Error{
		executor.startKubelet,
		executor.joinCluster,
		executor.appendLabel,
		executor.appendAnnotation,
		executor.appendTaint,
	}

	for _, operation := range operations {
		err := operation()
		if err != nil {
			return err
		}
	}

	executor.logger.Info("deploy node finished")

	return nil
}

func (executor *deployNodeExecutor) connectSSH() *protos.Error {

	executor.logger.Debug("Start to connect ssh")

	var err error
	executor.machine, err = deployMachine.NewMachine(executor.config.NodeCfg.GetNode())
	if err != nil {
		pbError := &protos.Error{
			Reason:     "Connect ssh error",                                                                                                                                   // 连接SSH失败。
			Detail:     fmt.Sprintf("SSH connect to %s(%s) failed , error: %v.", executor.config.NodeCfg.GetNode().GetName(), executor.config.NodeCfg.GetNode().GetIp(), err), // 连接%s(%s)失败，失败原因：%v。
			FixMethods: "Please check node reliability, make SSH service is available.",                                                                                       // 请检查节点的可用性，确保SSH服务可用。
		}
		executor.logger.WithField("error", pbError).Error("connect ssh error")
		return pbError
	}

	executor.logger.Debug("ssh connected")
	return nil
}

func (executor *deployNodeExecutor) connectMasterNode() *protos.Error {
	var err error
	executor.masterMachine, err = deployMachine.NewMachine(executor.config.MasterNodes[0])
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("failed to connect master node")
		return &protos.Error{
			Reason:     "connecting failed",
			Detail:     fmt.Sprintf("failed to connect master node, err: %s", err),
			FixMethods: "please check deploy node config to ensure master node can be connected successfully",
		}
	}
	return nil
}

func (executor *deployNodeExecutor) disconnectMasterNode() {
	if executor.masterMachine != nil {
		executor.masterMachine.Close()
	}
}

func (executor *deployNodeExecutor) initLogger() {
	executor.logger = logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: executor.action.GetName(),
		"nodeName":            executor.config.NodeCfg.GetNode().GetName(),
		"nodeIP":              executor.config.NodeCfg.GetNode().GetIp(),
	})
}

func (executor *deployNodeExecutor) startKubelet() *protos.Error {

	executor.logger.Debug("Start to install kubelet")

	operation := worker.NewStartKubelet(
		&worker.StartKubeletConfig{
			Machine:          executor.machine,
			Node:             executor.config.NodeCfg,
			Logger:           executor.logger,
			ExecuteLogWriter: executor.executeLogWriter,
		},
	)

	if err := operation.Execute(); err != nil {
		executor.logger.WithField("error", err).Error("install kubelet error")
		return err
	}

	executor.logger.Info("Finish to install kubelet action")
	return nil
}

func (executor *deployNodeExecutor) joinCluster() *protos.Error {

	executor.logger.Debug("Start to join cluster")

	operation := worker.NewJoinCluster(
		&worker.JoinClusterConfig{
			Machine:          executor.machine,
			Node:             executor.config.NodeCfg,
			Logger:           executor.logger,
			Cluster:          executor.config.ClusterConfig,
			MasterNodes:      executor.config.MasterNodes,
			ExecuteLogWriter: executor.executeLogWriter,
		},
	)

	if err := operation.Execute(); err != nil {
		executor.logger.WithField("error", err).Error("join cluster error")
		return err
	}

	executor.logger.Info("Finish to join cluster action")
	return nil
}

func (executor *deployNodeExecutor) appendLabel() *protos.Error {

	executor.logger.Debug("Start to append label")

	operation := worker.NewAppendLabel(
		&worker.AppendLabelConfig{
			MasterMachine:    executor.masterMachine,
			Logger:           executor.logger,
			Node:             executor.config.NodeCfg,
			Cluster:          executor.config.ClusterConfig,
			ExecuteLogWriter: executor.executeLogWriter,
		},
	)

	if err := operation.Execute(); err != nil {
		executor.logger.WithField("error", err).Error("append label error")
		return err
	}

	executor.logger.Info("Finish to append label action")
	return nil
}

func (executor *deployNodeExecutor) appendAnnotation() *protos.Error {

	executor.logger.Debug("Start to append annotation")

	operation := worker.NewAppendAnnotation(
		&worker.AppendAnnotationConfig{
			MasterMachine:    executor.masterMachine,
			Logger:           executor.logger,
			Node:             executor.config.NodeCfg,
			Cluster:          executor.config.ClusterConfig,
			ExecuteLogWriter: executor.executeLogWriter,
		},
	)

	if err := operation.Execute(); err != nil {
		executor.logger.WithField("error", err).Error("append annotation error")
		return err
	}

	executor.logger.Info("Finish to append annotation action")
	return nil
}

func (executor *deployNodeExecutor) appendTaint() *protos.Error {

	executor.logger.Debug("Start to append taint")

	operation := worker.NewAppendTaint(
		&worker.AppendTaintConfig{
			Machine:          executor.masterMachine,
			Logger:           executor.logger,
			Node:             executor.config.NodeCfg,
			Cluster:          executor.config.ClusterConfig,
			ExecuteLogWriter: executor.executeLogWriter,
		},
	)

	if err := operation.Execute(); err != nil {
		executor.logger.WithField("error", err).Error("append taint error")
		return err
	}

	executor.logger.Info("Finish to append taint action")
	return nil
}

func (executor *deployNodeExecutor) disconnectSSH() {

	if executor.machine == nil {
		return
	}

	executor.logger.Debug("Start to disconnect ssh")

	executor.machine.Close()

	executor.logger.Debug("ssh disconnected")
}

func (executor *deployNodeExecutor) initExecuteLogWriter() {

	executor.executeLogWriter = executor.action.GetExecuteLogBuffer()
}
