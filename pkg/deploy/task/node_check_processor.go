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
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	RegisterProcessor(TaskTypeNodeCheck, new(nodeCheckProcessor))
}

// nodeCheckProcessor implements the specific logic for the node check task.
type nodeCheckProcessor struct {
}

// Spilt the task into one or more node check actions
func (p *nodeCheckProcessor) SplitTask(t Task) error {
	if err := p.verifyTask(t); err != nil {
		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldTask: t.GetName(),
	})

	logger.Debug("Start to split node check task")

	checkTask := t.(*NodeCheckTask)

	// split task into actions: will create a action for every node, the action type
	// is NodeCheckAction
	actions := make([]action.Action, 0, len(checkTask.NodeConfigs))
	for _, subConfig := range checkTask.NodeConfigs {
		actionCfg := &action.NodeCheckActionConfig{
			NodeCheckConfig: subConfig,
			LogFileBasePath: checkTask.LogFileDir,
		}
		act, err := action.NewNodeCheckAction(actionCfg)
		if err != nil {
			return err
		}
		actions = append(actions, act)
	}

	if checkTask.NetworkOptions == nil {
		logger.Debugf("skip checking network requirements since networkOptions is empty")
	} else if len(checkTask.NodeConfigs) > 1 {
		// split into connectivity check actions
		numNodes := len(checkTask.NodeConfigs)
		for i, subConfig := range checkTask.NodeConfigs {
			randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
			// choose the index of peer. If index of itself is chosen, use the last node instead.
			peerIndex := randGen.Intn(numNodes - 1)
			if peerIndex == i {
				peerIndex = numNodes - 1
			}
			// make a connectivity check action for the pair.
			calicoOptions := checkTask.NetworkOptions.CalicoOptions
			if calicoOptions == nil {
				calicoOptions = &protos.CalicoOptions{
					EncapsulationMode: "vxlan",
					VxlanPort:         4789,
				}
			}
			act, err := makeConnectivityCheckActionCalico(
				subConfig.Node, checkTask.NodeConfigs[peerIndex].Node,
				calicoOptions, checkTask.GetLogFileDir())
			if err != nil {
				logger.WithField("node", subConfig.Node.Name).
					WithField("peer-node", checkTask.NodeConfigs[peerIndex].Node.Name).
					WithField("error", err.Error()).Warningf("failed to make connectivity check action")
			} else {
				actions = append(actions, act)
			}
		}
	}

	checkTask.Actions = actions
	logrus.Debugf("Finish to split node check task: %d actions", len(actions))
	return nil
}

// Verify if the task is valid.
func (p *nodeCheckProcessor) verifyTask(t Task) error {
	if t == nil {
		return consts.ErrEmptyTask
	}

	nodeCheckTask, ok := t.(*NodeCheckTask)
	if !ok {
		return fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, t)
	}

	if len(nodeCheckTask.NodeConfigs) == 0 {
		return fmt.Errorf("nodeConfigs is empty")
	}

	return nil
}
