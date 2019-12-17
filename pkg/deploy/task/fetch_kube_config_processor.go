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
	RegisterProcessor(TaskTypeFetchKubeConfig, new(fetchKubeConfigProcessor))
}

// fetchKubeConfigProcessor implements the specific logic for the fetch-kube-config task.
type fetchKubeConfigProcessor struct {
}

// Spilt the task into one fetch-kube-config action
func (p *fetchKubeConfigProcessor) SplitTask(t Task) error {
	if err := p.verifyTask(t); err != nil {
		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: t.GetName(),
	})

	logger.Debug("Start to split task")

	kubeCfgTask := t.(*FetchKubeConfigTask)

	// split task into an actions
	act, err := action.NewFetchKubeConfigAction(&action.FetchKubeConfigActionConfig{
		Node:            kubeCfgTask.Node,
		LogFileBasePath: kubeCfgTask.LogFilePath,
	})
	if err != nil {
		return err
	}
	kubeCfgTask.Actions = []action.Action{act}

	logrus.Debugf("Finish to split task")
	return nil
}

func (p *fetchKubeConfigProcessor) ProcessExtraResult(t Task) error {
	if err := p.verifyTask(t); err != nil {
		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: t.GetName(),
	})

	kubeCfgTask := t.(*FetchKubeConfigTask)
	if len(kubeCfgTask.Actions) == 0 {
		logger.Debug("Task has no action")
		return nil
	}

	kubeCfgAction, ok := kubeCfgTask.Actions[0].(*action.FetchKubeConfigAction)
	if !ok {
		return fmt.Errorf("%s: %T", consts.MsgActionTypeMismatched, kubeCfgTask.Actions[0])
	}

	kubeCfgTask.KubeConfig = kubeCfgAction.KubeConfig
	logger.Debugf("KubeConfig: %v", kubeCfgTask.KubeConfig)
	return nil
}

// Verify if the task is valid.
func (p *fetchKubeConfigProcessor) verifyTask(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	kubeCfgTask, ok := t.(*FetchKubeConfigTask)
	if !ok {
		return fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, t)
	}

	if kubeCfgTask.Node == nil {
		return fmt.Errorf("node field is nil")
	}

	return nil
}
