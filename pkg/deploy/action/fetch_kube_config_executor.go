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
)

type fetchKubeConfigExecutor struct {
}

func (a *fetchKubeConfigExecutor) Execute(act Action) error {
	kubeCfgAction, ok := act.(*FetchKubeConfigAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be fetch kube config action, but is %T", act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	logger.Debug("Start to execute action")

	// TODO: ssh to fetch kube config file from kubeCfgAction.node

	// Update action
	kubeCfgAction.status = ActionDone
	kubeCfgAction.KubeConfig = "todo: the content of kube config file"

	logger.Debug("Finsih to execute action")
	return nil
}
