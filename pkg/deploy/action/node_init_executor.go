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
	it "github.com/kpaas-io/kpaas/pkg/deploy/operation/init"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type nodeInitExecutor struct{}

func newNodeInitItem() *NodeInitItem {

	return &NodeInitItem{
		Status: ItemActionPending,
		Err:    &pb.Error{},
	}
}

// due to items, ItemInitScripts exec remote scripts and return std, report, error
func ExecuteInitScript(item it.ItemEnum, node *pb.Node, initItemReport *NodeInitItem) (string, *NodeInitItem, error) {
	logger := logrus.WithFields(logrus.Fields{
		"error_reason": fmt.Sprintf("failed to run script on node: %v", node.Name),
	})

	initItemReport = &NodeInitItem{
		Name:        fmt.Sprintf("init %v", item),
		Description: fmt.Sprintf("init %v", item),
	}

	initItem := it.NewInitOperations().CreateOperations(item)
	if initItem == nil {
		initItemReport.Status = ItemActionFailed
		initItemReport.Err.Reason = ItemErrEmpty
		initItemReport.Err.Detail = ItemErrEmpty
		initItemReport.Err.FixMethods = ItemHelperEmpty
		logger.Errorf("can not create %v operation", item)
		return "", initItemReport, fmt.Errorf("can not create %v's operation for node: %v", item, node.Name)
	}

	op, err := initItem.GetOperations(node)
	if err != nil {
		initItemReport.Status = ItemActionFailed
		initItemReport.Err.Reason = ItemErrOperation
		initItemReport.Err.Detail = err.Error()
		initItemReport.Err.FixMethods = ItemHelperOperation
		logger.Errorf("can not create operation command for %v", item)
		return "", initItemReport, fmt.Errorf("can not create operation command %v for node: %v", item, node.Name)
	}

	stdErr, stdOut, err := op.Do()
	if err != nil {
		initItemReport.Status = ItemActionFailed
		initItemReport.Err.Reason = ItemErrScript
		initItemReport.Err.Detail = string(stdErr)
		initItemReport.Err.FixMethods = ItemHelperScript
		logger.Errorf("can not execute %v operation", item)
		return "", initItemReport, fmt.Errorf("can not execute %v operation command on node: %v", item, node.Name)
	}

	initItemStdOut := string(stdOut)
	return initItemStdOut, initItemReport, nil
}

func (a *nodeInitExecutor) Execute(act Action) error {
	nodeInitAction, ok := act.(*NodeInitAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be node init action, but is %T", act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	logger.Debug("Start to execute node init action")

	// init sample
	// close firewall
	initItemReport := newNodeInitItem()
	initItemReport.Status = ItemActionDoing
	fireWallStdOut, initItemReport, err := ExecuteInitScript(it.FireWall, nodeInitAction.Node, initItemReport)
	if err != nil {
		initItemReport.Status = ItemActionFailed
	}
	UpdateInitItems(nodeInitAction, initItemReport)

	logger.Debugf("firewall stdout: %v", fireWallStdOut)

	// TODO Other Init Items
	// 1. Hostalias
	// 2. Hostname
	// 3. Network
	// 4. Route
	// 5. Swap
	// 6. TimeZone
	// 7. Haproxy
	// 8. Keepalived

	nodeInitAction.Status = ActionDone
	logger.Debug("Finish to execute node init action")
	return nil
}

// update init items with matching name
func UpdateInitItems(initAction *NodeInitAction, report *NodeInitItem) {

	initAction.Lock()
	defer initAction.Unlock()

	updatedFlag := false

	for _, item := range initAction.InitItems {
		if item.Name == report.Name {
			updatedFlag = true
			item.Err = report.Err
			item.Status = report.Status
			item.Description = report.Description
			break
		}
	}

	if updatedFlag == false {
		initAction.InitItems = append(initAction.InitItems, report)
	}
}
