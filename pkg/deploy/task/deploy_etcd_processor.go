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
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/etcd"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
)

// deployEtcdProcessor implements the specific logic to deploy etcd
type deployEtcdProcessor struct {
}

// Spilt the task into one or more node check actions
func (p *deployEtcdProcessor) SplitTask(t Task) error {
	if err := p.verifyTask(t); err != nil {
		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: t.GetName(),
	})

	logger.Debug("Start to split deploy etcd task")

	etcdTask := t.(*deployEtcdTask)

	// generate etcd ca cert and key and put it into every action
	caCertConfig := etcd.GetCaCrtConfig()
	caCrt, cakey, err := etcd.CreateAsCA(caCertConfig)
	if err != nil {
		return fmt.Errorf("failed to get etcd-ca key and cert, error: %v", err)
	}

	// split task into actions: will create a action for every node, the action type
	// is ActionTypeDeployEtcd
	actions := make([]action.Action, 0, len(etcdTask.nodes))
	for _, node := range etcdTask.nodes {
		actionCfg := &action.DeployEtcdActionConfig{
			CaCrt:           caCrt,
			CaKey:           cakey,
			Node:            node,
			ClusterNodes:    etcdTask.nodes,
			LogFileBasePath: etcdTask.logFilePath,
		}
		act, err := action.NewDeployEtcdAction(actionCfg)
		if err != nil {
			return err
		}
		actions = append(actions, act)
	}
	etcdTask.actions = actions

	logger.Debugf("Finish to split deploy etcd task: %d actions", len(actions))

	return nil
}

// Verify if the task is valid.
func (p *deployEtcdProcessor) verifyTask(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	etcdTask, ok := t.(*deployEtcdTask)
	if !ok {
		return fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, t)
	}

	if len(etcdTask.nodes) == 0 {
		return fmt.Errorf("nodes is empty")
	}

	return nil
}
