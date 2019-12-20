package action

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

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// Executor represents the interface of an action executor.
// Concrete executors implements the logic of actions.
type Executor interface {
	Execute(act Action) *pb.Error
}

var _executorRegistry map[Type]Executor

// RegisterExecutor is to register an Executor for an action type
func RegisterExecutor(actionType Type, exec Executor) error {
	if _executorRegistry == nil {
		_executorRegistry = make(map[Type]Executor)
	}
	if exec == nil {
		err := fmt.Errorf("the Executor to be registered is nil")
		logrus.Error(err)
		return err
	}
	if _, ok := _executorRegistry[actionType]; ok {
		err := fmt.Errorf("the Executor for type %v has already been registered", actionType)
		logrus.Error(err)
		return err
	}
	_executorRegistry[actionType] = exec
	return nil
}

// NewExecutor is a simple factory method to return an action executor based on action type.
func NewExecutor(actionType Type) (Executor, error) {
	exec, ok := _executorRegistry[actionType]
	if !ok {
		return nil, fmt.Errorf("%s: %s", consts.MsgActionTypeUnsupported, actionType)
	}

	return exec, nil
}

// ExecuteAction creates and run the executor for an action,
// a *sync.WaitGroup should be passed in.
func ExecuteAction(act Action, wg *sync.WaitGroup) {
	defer wg.Done()

	if act == nil {
		return
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction:     act.GetName(),
		consts.LogFieldActionType: act.GetType(),
	})

	logger.Debug("Start to execute action")

	executor, err := NewExecutor(act.GetType())
	if err != nil {
		act.SetStatus(ActionFailed)
		act.SetErr(&pb.Error{
			Reason: consts.MsgActionExecutorCreationFailed,
			Detail: err.Error(),
		})
		deploy.PBErrLogger(act.GetErr(), logger).Error()
		return
	}

	act.SetStatus(ActionDoing)

	if err := setup(act); err != nil {
		act.SetStatus(ActionFailed)
		act.SetErr(&pb.Error{
			Reason: "failed to setup",
			Detail: err.Error(),
		})
		deploy.PBErrLogger(act.GetErr(), logger).Error()
		return
	}

	defer writeExecuteLogs(act)

	if exeErr := executor.Execute(act); exeErr != nil {
		act.SetStatus(ActionFailed)
		act.SetErr(exeErr)
		deploy.PBErrLogger(act.GetErr(), logger).Error()
		return
	}

	act.SetStatus(ActionDone)
	logger.Debug("Finish to execute action")
}

func errOfTypeMismatched(expected, actual interface{}) *pb.Error {
	return &pb.Error{
		Reason: consts.MsgActionTypeMismatched,
		Detail: fmt.Sprintf(consts.MsgActionTypeMismatchedDetail, expected, actual),
	}
}

// Do some setup work before execut the action, like check and create log file...
func setup(act Action) error {
	if act == nil {
		return consts.ErrEmptyAction
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	// Create the log file. If failed to create the log file, just
	// log a warning and go on.
	logFilePath := act.GetLogFilePath()
	if logFilePath == "" {
		logger.Warn("The 'LogFilePath' field is empty")
	} else {
		file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
		if err != nil {
			logger.Warnf("Failed to create the log file: %s", err)
		} else {
			// Write the header info into the log file
			if err = writeActionLogHeader(file, act); err != nil {
				logger.Warnf("Failed to write log header: %s", err)
			}
			// Close the file
			if err = file.Close(); err != nil {
				logger.Warnf("Failed to close the log file: %s", err)
			}
		}
	}

	// TODO: setup exeute log buffer with a thread safe ReadWriter.
	if act.GetExecuteLogBuffer() == nil {
		act.SetExecuteLogBuffer(&bytes.Buffer{})
	}
	return nil
}

// Write action information into the log file
func writeActionLogHeader(file *os.File, act Action) error {
	if file == nil {
		return fmt.Errorf("invalid file descriptor")
	}
	if act == nil {
		return consts.ErrEmptyAction
	}

	_, err := file.WriteString("# action logs \n")
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("# action type: %v\n", act.GetType()))
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("# action name: %v\n", act.GetName()))
	if err != nil {
		return err
	}
	var nodeName string
	if node := act.GetNode(); node != nil {
		nodeName = node.GetName()
	}
	_, err = file.WriteString(fmt.Sprintf("# action node: %v\n", nodeName))
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("# action creation timestamp: %v\n", act.GetCreationTimestamp()))
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("# action log file path: %v\n", act.GetLogFilePath()))
	if err != nil {
		return err
	}

	return nil
}

func writeExecuteLogs(act Action) {
	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})
	logFilePath := act.GetLogFilePath()
	if logFilePath == "" {
		logger.Warning("action log file not specified")
		return
	}
	file, err := os.OpenFile(
		logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(0644))
	if err != nil {
		logger.WithField("error", err).WithField("file", logFilePath).
			Warningf("failed to open log file %s", logFilePath)
		return
	}
	defer file.Close()
	buf := act.GetExecuteLogBuffer()
	if buf == nil {
		logger.Warning("action did not write any execute log")
		return
	}

	_, err = file.WriteString("# action execute logs: \n")
	if err != nil {
		logger.WithField("error", err).WithField("file", logFilePath).
			Warningf("failed to write to log files %s", logFilePath)
		return
	}
	_, err = io.Copy(file, buf)
	if err != nil {
		logger.WithField("error", err).WithField("file", logFilePath).
			Warningf("failed to write to log files %s", logFilePath)
		return
	}
	// TODO: run buf.Flush() here?
	return
}
