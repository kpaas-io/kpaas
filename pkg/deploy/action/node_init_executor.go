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
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	it "github.com/kpaas-io/kpaas/pkg/deploy/operation/init"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	InitPassed = "init passed"
	InitFailed = "init failed"
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
		Description: fmt.Sprintf("初始化 %v 环境", item),
		Err:         new(pb.Error),
	}

	initAction := &operation.NodeInitAction{
		NodeInitConfig: action.NodeInitConfig,
		NodesConfig:    action.NodesConfig,
		ClusterConfig:  action.ClusterConfig,
	}

	initItem := it.NewInitOperations().CreateOperations(item, initAction)
	if initItem == nil {
		logger.Error("can not create operation")
		initItemReport.Status = ItemFailed
		initItemReport.Err.Reason = ItemErrEmpty
		initItemReport.Err.Detail = ItemErrEmpty
		initItemReport.Err.FixMethods = ItemHelperEmpty
		return "", initItemReport, fmt.Errorf("can not create %v's operation for node: %v", item, action.Node.Name)
	}

	// close ssh client
	defer initItem.CloseSSH()

	op, err := initItem.GetOperations(action.Node, initAction)
	if err != nil {
		logger.Errorf("can not create operation command, err: %v", err)
		initItemReport.Status = ItemFailed
		initItemReport.Err.Reason = ItemErrOperation
		initItemReport.Err.Detail = err.Error()
		initItemReport.Err.FixMethods = ItemHelperOperation
		return "", initItemReport, fmt.Errorf("can not create operation command %v for node: %v", item, action.Node.Name)
	}

	stdOut, stdErr, err := op.Do()
	if err != nil {
		logger.Errorf("can not execute operation command, err: %v", err)
		initItemReport.Status = ItemFailed
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
		Status: ItemPending,
		Err:    &pb.Error{},
	}
}

// goroutine exec item init event
func InitAsyncExecutor(item it.ItemEnum, ncAction *NodeInitAction, wg *sync.WaitGroup) {

	logger := logrus.WithFields(logrus.Fields{
		"node":      ncAction.Node.GetName(),
		"init_item": item,
	})

	logger.Debugf("Start to execute init")

	initItemReport := newNodeInitItem()
	initItemReport.Status = ItemDoing
	_, initItemReport, err := ExecuteInitScript(item, ncAction, initItemReport)
	if err != nil {
		logger.Errorf("%v: %v", InitFailed, err)
		initItemReport.Status = ItemFailed
	} else {
		initItemReport.Status = ItemDone
		logger.Info(InitPassed)
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
	logger.Debug("Start to execute node init action")

	// build init group contains multi roles all items
	var wg sync.WaitGroup
	var initGroup []it.ItemEnum
	initMap := make(map[it.ItemEnum]bool)

	initGroup = constructInitGroup(nodeInitAction, initMap)
	if len(initGroup) == 0 {
		logger.Error("init items group is empty")
	}

	for _, item := range initGroup {
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

// check if contains role
func containsRole(initAction *NodeInitAction, wantRole constant.MachineRole) bool {
	for _, role := range initAction.NodeInitConfig.Roles {
		if role == string(wantRole) {
			return true
		}
	}
	return false
}

// construct an init group contains items for one or more roles initiation
func constructInitGroup(nodeInitAction *NodeInitAction, itMap map[it.ItemEnum]bool) []it.ItemEnum {
	var initGroup []it.ItemEnum

	regularItemEnums := []it.ItemEnum{it.HostName, it.Swap, it.Route, it.Network, it.FireWall, it.TimeZone, it.HostName, it.HostAlias, it.KubeTool}

	// add init items by roles is supported based on regular items
	etcdItemEnums := regularItemEnums
	workerItemEnums := regularItemEnums
	ingressItemEnums := regularItemEnums
	masterItemEnums := regularItemEnums // cloud machine can not test it.Haproxy, it.Keepalived}

	if containsRole(nodeInitAction, constant.MachineRoleEtcd) {
		initGroup = addNotContainsItems(etcdItemEnums, itMap, initGroup)
	}

	if containsRole(nodeInitAction, constant.MachineRoleMaster) {
		initGroup = addNotContainsItems(masterItemEnums, itMap, initGroup)
	}

	if containsRole(nodeInitAction, constant.MachineRoleIngress) {
		initGroup = addNotContainsItems(ingressItemEnums, itMap, initGroup)
	}

	if containsRole(nodeInitAction, constant.MachineRoleWorker) {
		initGroup = addNotContainsItems(workerItemEnums, itMap, initGroup)
	}

	logrus.Debugf("node: %v, init group: %v", nodeInitAction.Node.Name, initGroup)

	return initGroup
}

// add items into array if not contains in it
func addNotContainsItems(initItems []it.ItemEnum, initMap map[it.ItemEnum]bool, initGroup []it.ItemEnum) []it.ItemEnum {
	for _, value := range initItems {
		if _, ok := initMap[value]; !ok {
			initMap[value] = true
			initGroup = append(initGroup, value)
		}
	}
	return initGroup
}
