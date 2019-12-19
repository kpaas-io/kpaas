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

package task

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
)

func init() {
	RegisterProcessor(TaskTypeTestConnection, new(testConnectionProcessor))
}

// testConnectionProcessor implements the specific logic for the test-connection task.
type testConnectionProcessor struct {
}

// Spilt the task into one test-connecton action
func (p *testConnectionProcessor) SplitTask(t Task) error {
	testConnTask, err := p.verifyTask(t)
	if err != nil {
		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: t.GetName(),
	})
	logger.Debug("Start to split test connection task")

	// split the task into one action
	actionCfg := &action.TestConnectionActionConfig{
		Node:            testConnTask.Node,
		LogFileBasePath: testConnTask.LogFileDir,
	}
	act, err := action.NewTestConnectionAction(actionCfg)
	if err != nil {
		return err
	}
	testConnTask.Actions = append(testConnTask.Actions, act)

	logrus.Debugf("Finish to split node check task: %d actions", len(testConnTask.Actions))
	return nil
}

// Verify if the task is valid.
func (p *testConnectionProcessor) verifyTask(t Task) (*TestConnectionTask, error) {
	if t == nil {
		return nil, consts.ErrEmptyTask
	}

	testConnectionTask, ok := t.(*TestConnectionTask)
	if !ok {
		return nil, fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, t)
	}

	if testConnectionTask.Node == nil {
		return nil, fmt.Errorf("Node field is empty")
	}

	return testConnectionTask, nil
}
