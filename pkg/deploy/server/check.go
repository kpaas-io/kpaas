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

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/deploy/task"
)

func (c *controller) getCheckNodeResult(aTask task.Task, withLogs bool) (*pb.GetCheckNodesResultReply, error) {
	if aTask == nil {
		return nil, fmt.Errorf("Task is nil")
	}

	// TODO: handle the logs of task and its actions if withLogs == true

	// Get all actions of the check task
	actions := task.GetAllActions(aTask)
	// Create a pb.NodeCheckResult for each action
	var nodeResults []*pb.NodeCheckResult
	for _, act := range actions {
		nodeCheckAct, ok := act.(*action.NodeCheckAction)
		if !ok {
			logrus.Warnf("Unexpected aciton type: %v", act.GetType())
			continue
		}
		nodeResult := checkActionToNodeCheckResult(nodeCheckAct)
		if nodeResult != nil {
			nodeResults = append(nodeResults, nodeResult)
		}
	}

	result := &pb.GetCheckNodesResultReply{
		Status: string(aTask.GetStatus()),
		Err:    aTask.GetErr(),
		Nodes:  nodeResults,
	}

	logrus.Debugf("Result: %+v", *result)

	return result, nil
}

func checkItemToItemCheckResult(actionItem *action.NodeCheckItem) *pb.ItemCheckResult {
	if actionItem == nil {
		return nil
	}

	return &pb.ItemCheckResult{
		Item: &pb.CheckItem{
			Name:        actionItem.Name,
			Description: actionItem.Description,
		},
		Status: string(actionItem.Status),
		Err:    actionItem.Err,
	}
}

func checkActionToNodeCheckResult(checkAction *action.NodeCheckAction) *pb.NodeCheckResult {
	if checkAction == nil {
		return nil
	}

	node := checkAction.GetNode()
	if node == nil {
		return nil
	}

	var nodeItems []*pb.ItemCheckResult
	for _, actionItem := range checkAction.CheckItems {
		nodeItem := checkItemToItemCheckResult(actionItem)
		if nodeItem != nil {
			nodeItems = append(nodeItems, nodeItem)
		}
	}

	return &pb.NodeCheckResult{
		NodeName: node.GetName(),
		Status:   string(checkAction.GetStatus()),
		Err:      checkAction.GetErr(),
		Items:    nodeItems,
	}
}
