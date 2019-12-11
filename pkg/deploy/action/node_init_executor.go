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

func newNodeInitItem() *nodeInitItem {

	return &nodeInitItem{
		status: ItemActionPending,
		err:    &pb.Error{},
	}
}

// due to items, ItemInitScripts exec remote scripts and return std, report, error
func ExecuteInitScript(item it.ItemEnum, node *pb.Node, initItemReport *nodeInitItem) (string, *nodeInitItem, error) {

	initItemReport = &nodeInitItem{
		name:        fmt.Sprintf("init %v", item),
		description: fmt.Sprintf("init %v", item),
	}

	initItem := it.NewInitOperations().CreateOperations(item)
	if initItem == nil {
		initItemReport.status = ItemActionFailed
		initItemReport.err.Reason = ItemErrEmpty
		initItemReport.err.Detail = ItemErrEmpty
		initItemReport.err.FixMethods = ItemHelperEmpty
	}

	op, err := initItem.GetOperations(node)
	if err != nil {
		initItemReport.status = ItemActionFailed
		initItemReport.err.Reason = ItemErrOperation
		initItemReport.err.Detail = err.Error()
		initItemReport.err.FixMethods = ItemHelperOperation
	}

	stdErr, stdOut, err := op.Do()
	if err != nil {
		initItemReport.status = ItemActionFailed
		initItemReport.err.Reason = ItemErrScript
		initItemReport.err.Detail = string(stdErr)
		initItemReport.err.FixMethods = ItemHelperScript
	}

	initItemStdOut := string(stdOut)
	return initItemStdOut, initItemReport, nil
}

func (a *nodeInitExecutor) Execute(act Action) error {
	nodeInitAction, ok := act.(*nodeInitAction)
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
	initItemReport.status = ItemActionDoing
	fireWallStdOut, initItemReport, err := ExecuteInitScript(it.FireWall, nodeInitAction.node, initItemReport)
	UpdateInitItems(nodeInitAction, initItemReport)
	if err != nil {
		initItemReport.status = ItemActionFailed
	}

	logger.Debugf("firewall std out: %v", fireWallStdOut)
	UpdateInitItems(nodeInitAction, initItemReport)

	// TODO Other Init Items
	// 1. Hostalias
	// 2. Hostname
	// 3. Network
	// 4. Route
	// 5. Swap
	// 6. TimeZone
	// 7. Haproxy
	// 8. Keepalived

	nodeInitAction.status = ActionDone
	logger.Debug("Finish to execute node init action")
	return nil
}

// update init items with matching name
func UpdateInitItems(initAction *nodeInitAction, report *nodeInitItem) {

	initAction.Lock()
	defer initAction.Unlock()

	updatedFlag := false

	for _, item := range initAction.initItems {
		if item.name == report.name {
			updatedFlag = true
			item.err = report.err
			item.status = report.status
			item.description = report.description
		}
	}

	if updatedFlag == false {
		initAction.initItems = append(initAction.initItems, report)
	}
}
