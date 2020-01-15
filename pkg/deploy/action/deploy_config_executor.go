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

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/config"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	RegisterExecutor(ActionTypeDeployConfig, new(deployConfigExecutor))
}

type deployConfigExecutor struct {
}

func (e *deployConfigExecutor) Execute(act Action) *pb.Error {
	configAction, ok := act.(*DeployConfigAction)
	if !ok {
		return errOfTypeMismatched(new(DeployConfigAction), act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
		consts.LogFieldNode:   act.GetNode().GetName(),
	})
	logger.Debug("Start to execute deploy config action")

	masterMachine, err := machine.NewMachine(configAction.MasterNodes[0])
	if err != nil {
		pbErr := &pb.Error{
			Reason: "failed to connect to target node",
			Detail: err.Error(),
		}
		return pbErr
	}
	defer masterMachine.Close()

	operations := []func(*DeployConfigAction, machine.IMachine, *logrus.Entry) *protos.Error{
		e.appendLabel,
		e.appendAnnotation,
		e.appendTaint,
	}

	var errs []*pb.Error
	for _, operation := range operations {
		err := operation(configAction, masterMachine, logger)
		if err != nil {
			// If operation failed, record the error and continue.
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		pbErr := &pb.Error{
			Reason:     "failed to do deploy config",
			Detail:     fmt.Sprint(errs),
			FixMethods: consts.FixMethodSelfAnalyseIt,
		}
		return pbErr
	}

	logger.Debug("Finish to execute deploy config action")
	return nil
}

func (e *deployConfigExecutor) appendLabel(act *DeployConfigAction, machine machine.IMachine, logger *logrus.Entry) *protos.Error {
	logger.Debug("Start to append label")

	appendLabel := config.NewAppendLabel(
		&config.AppendLabelConfig{
			MasterMachine:    machine,
			Logger:           logger,
			Node:             act.NodeConfig,
			Cluster:          act.ClusterConfig,
			ExecuteLogWriter: act.GetExecuteLogBuffer(),
		},
	)

	if err := appendLabel.Execute(); err != nil {
		logger.WithField("error", err).Error("append label error")
		return err
	}

	logger.Debug("Finish to append label")
	return nil
}

func (e *deployConfigExecutor) appendAnnotation(act *DeployConfigAction, machine machine.IMachine, logger *logrus.Entry) *protos.Error {
	logger.Debug("Start to append annotation")

	operation := config.NewAppendAnnotation(
		&config.AppendAnnotationConfig{
			MasterMachine:    machine,
			Logger:           logger,
			Node:             act.NodeConfig,
			Cluster:          act.ClusterConfig,
			ExecuteLogWriter: act.GetExecuteLogBuffer(),
		},
	)

	if err := operation.Execute(); err != nil {
		logger.WithField("error", err).Error("append annotation error")
		return err
	}

	logger.Debug("Finish to append annotation")
	return nil
}

func (e *deployConfigExecutor) appendTaint(act *DeployConfigAction, machine machine.IMachine, logger *logrus.Entry) *protos.Error {
	logger.Debug("Start to append taint")

	operation := config.NewAppendTaint(
		&config.AppendTaintConfig{
			MasterMachine:    machine,
			Logger:           logger,
			Node:             act.NodeConfig,
			Cluster:          act.ClusterConfig,
			ExecuteLogWriter: act.GetExecuteLogBuffer(),
		},
	)

	if err := operation.Execute(); err != nil {
		logger.WithField("error", err).Error("append taint error")
		return err
	}

	logger.Info("Finish to append taint")
	return nil
}
