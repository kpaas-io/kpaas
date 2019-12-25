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
	"strings"
	"sync"

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
		"node":      action.Node.GetName(),
		"init_item": item,
	})

	initItemReport = &NodeInitItem{
		Name:        fmt.Sprintf("init %v", item),
		Description: fmt.Sprintf("init %v", item),
		Err:         new(pb.Error),
	}

	initAction := &operation.NodeInitAction{
		NodeInitConfig: action.NodeInitConfig,
		NodesConfig:    action.NodesConfig,
		ClusterConfig:  action.ClusterConfig,
	}

	initItem := it.NewInitOperations().CreateOperations(item, initAction)
	if initItem == nil {
		logger.Errorf("can not create %v operation", item)
		initItemReport.Status = ItemActionFailed
		initItemReport.Err.Reason = ItemErrEmpty
		initItemReport.Err.Detail = ItemErrEmpty
		initItemReport.Err.FixMethods = ItemHelperEmpty
		return "", initItemReport, fmt.Errorf("can not create %v's operation for node: %v", item, action.Node.Name)
	}

	// close ssh client
	defer initItem.CloseSSH()

	op, err := initItem.GetOperations(action.Node, initAction)
	if err != nil {
		logger.Errorf("can not create %v operation command for %v", item, err)
		initItemReport.Status = ItemActionFailed
		initItemReport.Err.Reason = ItemErrOperation
		initItemReport.Err.Detail = err.Error()
		initItemReport.Err.FixMethods = ItemHelperOperation
		return "", initItemReport, fmt.Errorf("can not create operation command %v for node: %v", item, action.Node.Name)
	}

	stdOut, stdErr, err := op.Do()
	if err != nil {
		logger.Errorf("can not execute %v operation command for %v", item, err)
		initItemReport.Status = ItemActionFailed
		initItemReport.Err.Reason = ItemErrScript
		initItemReport.Err.Detail = string(stdErr)
		initItemReport.Err.FixMethods = ItemHelperScript
		return "", initItemReport, fmt.Errorf("can not execute %v operation command on node: %v", item, action.Node.Name)
	}

	initItemStdOut := strings.Trim(string(stdOut), "\n")

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

	logger := logrus.WithFields(logrus.Fields{
		"node":      ncAction.Node.GetName(),
		"init_item": item,
	})

	logrus.Debugf("Start to execute init %v, node: %v", item, ncAction.Node.GetName())

	initItemReport := newNodeInitItem()
	initItemReport.Status = ItemActionDoing
	_, initItemReport, err := ExecuteInitScript(item, ncAction, initItemReport)
	if err != nil {
		logger.Errorf("init %v failed, err: %v", item, err)
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
	itemEnums := []it.ItemEnum{it.Swap, it.Route, it.Network, it.Network, it.FireWall, it.TimeZone, it.HostName, it.HostAlias, it.KubeTool} // TODO it.Haproxy, it.Keepalived,
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
