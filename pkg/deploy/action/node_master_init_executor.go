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
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	it "github.com/kpaas-io/kpaas/pkg/deploy/operation/init"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type nodeMasterInitExecutor struct{}

// due to items, ItemInitScripts exec remote scripts and return std, report, error
func ExecuteMasterInitScript(item it.ItemEnum, node *pb.Node, initItemReport *NodeInitItem) (string, *NodeInitItem, error) {
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

	// close ssh client
	initItem.CloseSSH()

	initItemStdOut := string(stdOut)
	return initItemStdOut, initItemReport, nil
}

func newNodeMasterInitItem() *NodeMasterInitItem {

	return &NodeMasterInitItem{
		Status: ItemActionPending,
		Err:    &pb.Error{},
	}
}

// goroutine exec item init event
func InitMasterAsyncExecutor(item it.ItemEnum, ncAction *NodeMasterInitAction, wg *sync.WaitGroup) {

	initItemReport := newNodeInitItem()
	initItemReport.Status = ItemActionDoing
	_, initItemReport, err := ExecuteMasterInitScript(item, ncAction.Node, initItemReport)
	if err != nil {
		initItemReport.Status = ItemActionFailed
	}

	UpdateInitMasterItems(ncAction, initItemReport)

	wg.Done()
}

func (a *nodeMasterInitExecutor) Execute(act Action) *pb.Error {
	nodeMasterInitAction, ok := act.(*NodeMasterInitAction)
	if !ok {
		return errOfTypeMismatched(new(NodeMasterInitAction), act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	var wg sync.WaitGroup

	logger.Debug("Start to execute node init action")

	// init master node events include haproxy, keepalived

	itemEnums := []it.ItemEnum{it.Haproxy, it.Keepalived}
	for _, item := range itemEnums {
		wg.Add(1)
		go InitMasterAsyncExecutor(item, nodeMasterInitAction, &wg)
	}

	wg.Wait()

	// If any of init item was failed, we should return an error
	failedItems := getFailedInitMasterItems(nodeMasterInitAction)
	if len(failedItems) > 0 {
		return &pb.Error{
			Reason: fmt.Sprintf("%d init item(s) failed", len(failedItems)),
			Detail: fmt.Sprintf("failed init item list: %v", failedItems),
		}
	}

	logger.Debug("Finish to execute node init action")
	return nil
}

// update init items with matching name
func UpdateInitMasterItems(initAction *NodeMasterInitAction, report *NodeInitItem) {

	initAction.Lock()
	defer initAction.Unlock()

	initAction.InitItems = append(initAction.InitItems, report)
}

func getFailedInitMasterItems(initAction *NodeMasterInitAction) []string {
	var failedItemName []string
	for _, item := range initAction.InitItems {
		if item.Status != nodeInitItemDone {
			failedItemName = append(failedItemName, item.Name)
		}
	}
	return failedItemName
}
