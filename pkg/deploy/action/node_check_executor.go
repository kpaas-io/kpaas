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
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
)

const (
	desireDockerVersion               = "18.06.0"
	desireKernelVersion               = "4.19.46"
	desireCPUCore int                 = 8
	desiredMemoryBase float64         = 16
	desiredMemory                     = desiredMemoryBase * operation.GiByteUnits
	desiredRootDiskVolumeBase float64 = 200
	desiredRootDiskVolume     float64 = desiredRootDiskVolumeBase * operation.GiByteUnits
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
		detail string
		status nodeCheckItemStatus
		fixmethod string
	)

	op, err := docker.NewCheckDockerOperation(nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return fmt.Errorf("failed to create docker check operation, error: %v", err)
	}

	stdErr, stdOut, err := op.Do()
	if err != nil {
		reason = "run command failed"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = "please check your scripts"
	}

	comparedDockerVersion := string(stdOut[:])
	err = docker.NewDockerVersionCheck(comparedDockerVersion, desireDockerVersion, ".", ">")
	if err != nil {
		reason = "docker version not satisfied"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = fmt.Sprintf("please upgrade docker version to %v+", desireDockerVersion)
	} else {
		reason = ""
		detail = ""
		status = nodeCheckItemSucessful
		fixmethod = ""
	}

	dockerVersionItem := &nodeCheckItem{
		name:        "docker version check",
		description: "docker version check",
		status:      status,
		err: &pb.Error{
			Reason:     reason,
			Detail:     detail,
			FixMethods: fixmethod,
		},
	}
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, dockerVersionItem)

	// TODO: other checks

	// TODO: update action status
	nodeCheckAction.status = ActionFailed
	nodeCheckAction.err = dockerVersionItem.err

	logger.Debug("Finish to execute node check action")
	return nil
}
