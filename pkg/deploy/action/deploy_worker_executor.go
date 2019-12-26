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
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/worker"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	RegisterExecutor(ActionTypeDeployWorker, new(deployWorkerExecutor))
}

type deployWorkerExecutor struct {
	logger           *logrus.Entry
	machine          deployMachine.IMachine
	action           *DeployWorkerAction
	executeLogWriter io.Writer
}

func (executor *deployWorkerExecutor) Execute(act Action) *protos.Error {

	action, ok := act.(*DeployWorkerAction)
	if !ok {
		return errOfTypeMismatched(new(DeployWorkerAction), act)
	}

	executor.action = action

	executor.initLogger()
	executor.initExecuteLogWriter()
	defer executor.closeExecuteLogWriter()

	executor.logger.Info("start to execute deploy worker executor")

	if err := executor.connectSSH(); err != nil {

		return err
	}
	defer executor.disconnectSSH()

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

	executor.logger.Info("deploy worker finished")

	return nil
}

func (executor *deployWorkerExecutor) connectSSH() *protos.Error {

	executor.logger.Debug("Start to connect ssh")

	var err error
	executor.machine, err = deployMachine.NewMachine(executor.action.config.NodeCfg.GetNode())
	if err != nil {
		pbError := &protos.Error{
			Reason:     "Connect ssh error",                                                                                                                                                 // 连接SSH失败。
			Detail:     fmt.Sprintf("SSH connect to %s(%s) failed , error: %v.", executor.action.config.NodeCfg.GetNode().GetName(), executor.action.config.NodeCfg.GetNode().GetIp(), err), // 连接%s(%s)失败，失败原因：%v。
			FixMethods: "Please check node reliability, make SSH service is available.",                                                                                                     // 请检查节点的可用性，确保SSH服务可用。
		}
		executor.logger.WithField("error", pbError).Error("connect ssh error")
		return pbError
	}

	executor.logger.Debug("ssh connected")
	return nil
}

func (executor *deployWorkerExecutor) initLogger() {
	executor.logger = logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: executor.action.GetName(),
		"nodeName":            executor.action.config.NodeCfg.GetNode().GetName(),
		"nodeIP":              executor.action.config.NodeCfg.GetNode().GetIp(),
	})
}

func (executor *deployWorkerExecutor) startKubelet() *protos.Error {

	executor.logger.Debug("Start to install kubelet")

	operation := worker.NewStartKubelet(
		&worker.StartKubeletConfig{
			Machine:          executor.machine,
			Node:             executor.action.config.NodeCfg,
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

func (executor *deployWorkerExecutor) joinCluster() *protos.Error {

	executor.logger.Debug("Start to join cluster")

	operation := worker.NewJoinCluster(
		&worker.JoinClusterConfig{
			Machine:          executor.machine,
			Node:             executor.action.config.NodeCfg,
			Logger:           executor.logger,
			Cluster:          executor.action.config.ClusterConfig,
			MasterNodes:      executor.action.config.MasterNodes,
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

func (executor *deployWorkerExecutor) appendLabel() *protos.Error {

	executor.logger.Debug("Start to append label")

	operation := worker.NewAppendLabel(
		&worker.AppendLabelConfig{
			Machine:          executor.machine,
			Logger:           executor.logger,
			Node:             executor.action.config.NodeCfg,
			Cluster:          executor.action.config.ClusterConfig,
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

func (executor *deployWorkerExecutor) appendAnnotation() *protos.Error {

	executor.logger.Debug("Start to append annotation")

	operation := worker.NewAppendAnnotation(
		&worker.AppendAnnotationConfig{
			Machine:          executor.machine,
			Logger:           executor.logger,
			Node:             executor.action.config.NodeCfg,
			Cluster:          executor.action.config.ClusterConfig,
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

func (executor *deployWorkerExecutor) appendTaint() *protos.Error {

	executor.logger.Debug("Start to append taint")

	operation := worker.NewAppendTaint(
		&worker.AppendTaintConfig{
			Machine:          executor.machine,
			Logger:           executor.logger,
			Node:             executor.action.config.NodeCfg,
			Cluster:          executor.action.config.ClusterConfig,
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

func (executor *deployWorkerExecutor) disconnectSSH() {

	executor.logger.Debug("Start to disconnect ssh")

	executor.machine.Close()

	executor.logger.Debug("ssh disconnected")
}

func (executor *deployWorkerExecutor) initExecuteLogWriter() {

	if executor.action.LogFilePath == "" {
		return
	}

	var err error
	// LogFilePath /app/deploy/logs/unknown/deploy-worker/deploy-worker-{nodeName}.log
	err = os.MkdirAll(filepath.Dir(executor.action.LogFilePath), os.FileMode(0755))
	if err != nil {
		return
	}
	executor.executeLogWriter, err = os.OpenFile(executor.action.LogFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		executor.logger.Errorf("init deploy worker execute log writer error, error message: %v", err)
		executor.executeLogWriter = bytes.NewBuffer([]byte{})
		return
	}
}

func (executor *deployWorkerExecutor) closeExecuteLogWriter() {

	switch writer := executor.executeLogWriter.(type) {
	case *os.File:
		err := writer.Close()
		if err != nil {
			executor.logger.Errorf("close deploy worker execute log writer error, error message: %v", err)
		}
	case *bytes.Buffer:
		// Open executed log file handle error, so write to logrus(at least we can find the log)
		logrus.Infof("%s\n", writer.String())
	}
}
