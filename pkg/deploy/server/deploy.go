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
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/deploy/task"
)

var allRoles = []constant.MachineRole{constant.MachineRoleEtcd, constant.MachineRoleMaster,
	constant.MachineRoleWorker, constant.MachineRoleIngress}

func (c *controller) getDeployResult(aTask task.Task) (*pb.GetDeployResultReply, error) {
	if aTask == nil {
		return nil, fmt.Errorf("Task is nil")
	}

	deployTask, ok := aTask.(*task.DeployTask)
	if !ok {
		return nil, fmt.Errorf("invalid task")
	}
	roleNodes := groupNodesByRole(deployTask.NodeConfigs)
	// If the task is already failed, set the default status in deploy item result as "aborted",
	// otherewise, set the default status to "pending". The final status of them would be updated
	// in the following process.
	initStatus := string(constant.OperationStatusPending)
	if aTask.GetStatus() == task.TaskFailed {
		initStatus = string(constant.OperationStatusAborted)
	}
	// Create a pb.DeployItemResult for each {role, node}
	roleNodeDeployItemResult := make(map[constant.MachineRole]map[string]*pb.DeployItemResult)
	for role, nodes := range roleNodes {
		roleNodeDeployItemResult[role] = make(map[string]*pb.DeployItemResult)
		for _, node := range nodes {
			roleNodeDeployItemResult[role][node] = &pb.DeployItemResult{
				DeployItem: &pb.DeployItem{
					Role:     string(role),
					NodeName: node,
				},
				Status: initStatus,
			}
		}
	}

	// Get all actions of the deploy task
	actions := task.GetAllActions(aTask)

	// Firstly, iterate the node init action, if any of the node init action is not done
	// update the related pb.DeployItemResult in the collecton of {role, node} to the
	// node init aciton's status
	initNotDone := false
	for _, act := range actions {
		if act.GetType() != action.ActionTypeNodeInit || act.GetStatus() == action.ActionDone {
			continue
		}

		initNotDone = true

		node := act.GetNode()
		if node == nil || node.Name == "" {
			logrus.Warn("Invalid node")
			continue
		}

		for _, role := range allRoles {
			if _, ok := roleNodeDeployItemResult[role]; !ok {
				continue
			}
			if _, ok := roleNodeDeployItemResult[role][node.Name]; !ok {
				continue
			}
			itemResult := roleNodeDeployItemResult[role][node.Name]
			itemResult.Status = string(actionStatusToOperationStatus(act.GetStatus()))
			itemResult.Err = act.GetErr()
		}

	}

	// If all node init action are done, update deploy item results with non node init actions
	if !initNotDone {
		for _, act := range actions {
			if act.GetType() == action.ActionTypeNodeInit {
				continue
			}

			node := act.GetNode()
			if node == nil || node.Name == "" {
				logrus.Warn("Invalid node")
				continue
			}

			role := actionTypeToRole(act.GetType())
			if _, ok := roleNodeDeployItemResult[role]; !ok {
				logrus.Warnf("Didn't find the role %q in the map", role)
				continue
			}
			if _, ok := roleNodeDeployItemResult[role][node.Name]; !ok {
				logrus.Warnf("Didn't find the node %q with the role %q in the map", node.Name, role)
				continue
			}
			itemResult := roleNodeDeployItemResult[role][node.Name]
			itemResult.Status = string(actionStatusToOperationStatus(act.GetStatus()))
			itemResult.Err = act.GetErr()
		}
	}

	// Update the reply's status according to the deploy task's status.
	result := &pb.GetDeployResultReply{
		Status: string(taskStatusToOperationStatus(aTask.GetStatus())),
		Err:    aTask.GetErr(),
		Items:  sortResultByRole(roleNodeDeployItemResult),
	}

	logrus.Debugf("Result: %+v", *result)

	return result, nil
}
func sortMap(m map[string]*pb.DeployItemResult) []*pb.DeployItemResult {
	// get all keys
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	// sort the keys
	sort.Strings(keys)

	var items []*pb.DeployItemResult
	for _, key := range keys {
		items = append(items, m[key])
	}
	return items
}

func sortResultByRole(m map[constant.MachineRole]map[string]*pb.DeployItemResult) []*pb.DeployItemResult {
	var items []*pb.DeployItemResult
	// Sort the result by roles
	for _, role := range allRoles {
		if roleResultMap, ok := m[role]; ok {
			items = append(items, sortMap(roleResultMap)...)
		}
	}
	return items
}

func groupNodesByRole(cfgs []*pb.NodeDeployConfig) map[constant.MachineRole][]string {
	roleNodes := make(map[constant.MachineRole][]string)
	for _, nodeCfg := range cfgs {
		if nodeCfg == nil || nodeCfg.Node == nil || nodeCfg.Node.Name == "" || len(nodeCfg.Roles) == 0 {
			logrus.Warnf("Invalid nodeCfg")
			continue
		}
		for _, role := range nodeCfg.Roles {
			roleName := constant.MachineRole(role)
			roleNodes[roleName] = append(roleNodes[roleName], nodeCfg.Node.Name)
		}
	}
	return roleNodes
}

func actionTypeToRole(actionType action.Type) constant.MachineRole {
	switch actionType {
	case action.ActionTypeDeployEtcd:
		return constant.MachineRoleEtcd
	case action.ActionTypeInitMaster, action.ActionTypeJoinMaster:
		return constant.MachineRoleMaster
	case action.ActionTypeDeployNode:
		return constant.MachineRoleWorker
	// treat node init action as ectd role
	case action.ActionTypeNodeInit:
		return constant.MachineRoleEtcd
	default:
		logrus.Warnf("unknown action type: %v", actionType)
		return "unknown"
	}
}
