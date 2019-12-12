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

package server

import (
	"fmt"

	"github.com/sirupsen/logrus"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/deploy/task"
)

func (c *controller) getDeployResult(aTask task.Task, withLogs bool) (*pb.GetDeployResultReply, error) {
	if aTask == nil {
		return nil, fmt.Errorf("Task is nil")
	}

	// TODO: handle the logs of task and its actions if withLogs == true

	// Get all actions of the deploy task
	actions := task.GetAllActions(aTask)
	// Create a pb.DeployItemResult for each action
	var items []*pb.DeployItemResult
	for _, act := range actions {
		if node := act.GetNode(); node != nil {
			items = append(items, &pb.DeployItemResult{
				DeployItem: &pb.DeployItem{
					Role:     string(act.GetType()),
					NodeName: node.GetName(),
				},
				Status: string(act.GetStatus()),
				Err:    act.GetErr(),
			})
		}
	}

	result := &pb.GetDeployResultReply{
		Status: string(aTask.GetStatus()),
		Err:    aTask.GetErr(),
		Items:  items,
	}

	logrus.Debugf("Result: %+v", *result)

	return result, nil
}
