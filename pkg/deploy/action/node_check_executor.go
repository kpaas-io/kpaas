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
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/check/docker"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type nodeCheckExecutor struct {
}

func (a *nodeCheckExecutor) Execute(act Action) error {
	nodeCheckAction, ok := act.(*nodeCheckAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be node check action, but is %T", act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	logger.Debug("Start to execute node check action")

	var (
		reason string
		status nodeCheckItemStatus
	)

	op, err := docker.NewCheckDockerOperation(nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return fmt.Errorf("failed to create docker check operation, error: %v", err)
	}

	_, errOut, err := op.Do()
	if err != nil {
		reason = string(errOut)
		status = nodeCheckItemFailed
	}

	dockerVersionItem := &nodeCheckItem{
		name:        "docker version check",
		description: "docker version check",
		status:      status,
		err: &pb.Error{
			Reason:     reason,
			Detail:     "",
			FixMethods: "upgrade docker version to 17.03.02+",
		},
	}
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, dockerVersionItem)

	// TODO: other checks

	// TODO: update action status
	nodeCheckAction.status = ActionFailed
	nodeCheckAction.err = dockerVersionItem.err

	logger.Debug("Finsih to execute node check action")
	return nil
}
