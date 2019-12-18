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
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	it "github.com/kpaas-io/kpaas/pkg/deploy/operation/init"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	RegisterExecutor(ActionTypeNodeInit, new(nodeInitExecutor))
}

type nodeInitExecutor struct{}

// due to items, ItemInitScripts exec remote scripts and return std, report, error
func ExecuteInitScript(item it.ItemEnum, action *NodeInitAction, initItemReport *NodeInitItem) (string, *NodeInitItem, error) {
	logger := logrus.WithFields(logrus.Fields{
		"error_reason": fmt.Sprintf("failed to run script on node: %v", action.Node.Name),
	})

	initItemReport = &NodeInitItem{
		Name:        fmt.Sprintf("init %v", item),
		Description: fmt.Sprintf("init %v", item),
	}

	initAction := &operation.NodeInitAction{
		NodeInitConfig: action.NodeInitConfig,
		NodesConfig:    action.NodesConfig,
		ClusterConfig:  action.ClusterConfig,
	}

	initItem := it.NewInitOperations().CreateOperations(item, initAction)
	if initItem == nil {
		initItemReport.Status = ItemActionFailed
		initItemReport.Err.Reason = ItemErrEmpty
		initItemReport.Err.Detail = ItemErrEmpty
		initItemReport.Err.FixMethods = ItemHelperEmpty
		logger.Errorf("can not create %v operation", item)
		return "", initItemReport, fmt.Errorf("can not create %v's operation for node: %v", item, action.Node.Name)
	}

	// close ssh client
	defer initItem.CloseSSH()

	op, err := initItem.GetOperations(action.Node, initAction)
	if err != nil {
		initItemReport.Status = ItemActionFailed
		initItemReport.Err.Reason = ItemErrOperation
		initItemReport.Err.Detail = err.Error()
		initItemReport.Err.FixMethods = ItemHelperOperation
		logger.Errorf("can not create operation command for %v", item)
		return "", initItemReport, fmt.Errorf("can not create operation command %v for node: %v", item, action.Node.Name)
	}

	stdErr, stdOut, err := op.Do()
	if err != nil {
		initItemReport.Status = ItemActionFailed
		initItemReport.Err.Reason = ItemErrScript
		initItemReport.Err.Detail = string(stdErr)
		initItemReport.Err.FixMethods = ItemHelperScript
		logger.Errorf("can not execute %v operation", item)
		return "", initItemReport, fmt.Errorf("can not execute %v operation command on node: %v", item, action.Node.Name)
	}

	initItemStdOut := string(stdOut)
	return initItemStdOut, initItemReport, nil
}

func newNodeInitItem() *NodeInitItem {

	return &NodeInitItem{
		Status: ItemActionPending,
		Err:    &pb.Error{},
	}
}

// goroutine exec item init event
func InitAsyncExecutor(item it.ItemEnum, ncAction *NodeInitAction, wg *sync.WaitGroup) {

	initItemReport := newNodeInitItem()
	initItemReport.Status = ItemActionDoing
	_, initItemReport, err := ExecuteInitScript(item, ncAction, initItemReport)
	if err != nil {
		initItemReport.Status = ItemActionFailed
	}

	UpdateInitItems(ncAction, initItemReport)

	wg.Done()
}

func (a *nodeInitExecutor) Execute(act Action) *pb.Error {
	nodeInitAction, ok := act.(*NodeInitAction)
	if !ok {
		return errOfTypeMismatched(new(NodeInitAction), act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	var wg sync.WaitGroup

	logger.Debug("Start to execute node init action")

	// init events include firewall, hostalias, hostname, network
	// route, swap, timezone， kubetool
	itemEnums := []it.ItemEnum{it.Swap, it.Route, it.Network, it.Network, it.FireWall, it.TimeZone,
		it.HostName, it.HostAlias, it.Haproxy, it.Keepalived, it.KubeTool}
	for _, item := range itemEnums {
		wg.Add(1)
		go InitAsyncExecutor(item, nodeInitAction, &wg)
	}

	wg.Wait()

	// If any of init item was failed, we should return an error
	failedItems := getFailedInitItems(nodeInitAction)
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
func UpdateInitItems(initAction *NodeInitAction, report *NodeInitItem) {

	initAction.Lock()
	defer initAction.Unlock()

	initAction.InitItems = append(initAction.InitItems, report)
}

func getFailedInitItems(initAction *NodeInitAction) []string {
	var failedItemName []string
	for _, item := range initAction.InitItems {
		if item.Status != nodeInitItemDone {
			failedItemName = append(failedItemName, item.Name)
		}
	}
	return failedItemName
}
