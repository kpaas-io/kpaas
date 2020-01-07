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

	initGroup := make(map[it.ItemEnum]bool)
	initGroup = constructInitGroup(nodeInitAction)
	if len(initGroup) == 0 {
		logger.Error("initialization item group is empty")
	}

	var wg sync.WaitGroup
	for item := range initGroup {
		wg.Add(1)
		if initGroup[item] {
			go InitAsyncExecutor(item, nodeInitAction, &wg)
		}
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

// according to roles, construct an init item group
func constructInitGroup(nodeInitAction *NodeInitAction) map[it.ItemEnum]bool {
	initGroup := make(map[it.ItemEnum]bool)

	var etcdItemEnums []it.ItemEnum
	var masterItemEnums []it.ItemEnum
	var workerItemEnums []it.ItemEnum
	var ingressItemEnums []it.ItemEnum

	baseItemEnums := []it.ItemEnum{it.HostName, it.Swap, it.Route, it.Network, it.FireWall, it.TimeZone, it.HostName, it.HostAlias, it.KubeTool}

	if nodeInitAction.ClusterConfig.KubeAPIServerConnect.Type == "keepalived" {
		masterItemEnums = []it.ItemEnum{it.Haproxy, it.Keepalived}
	}

	if containsRole(nodeInitAction, constant.MachineRoleEtcd) {
		baseItemEnums = append(baseItemEnums, etcdItemEnums...)
	}

	if containsRole(nodeInitAction, constant.MachineRoleMaster) {
		baseItemEnums = append(baseItemEnums, masterItemEnums...)
	}

	if containsRole(nodeInitAction, constant.MachineRoleIngress) {
		baseItemEnums = append(baseItemEnums, ingressItemEnums...)
	}

	if containsRole(nodeInitAction, constant.MachineRoleWorker) {
		baseItemEnums = append(baseItemEnums, workerItemEnums...)
	}

	logrus.Debugf("node: %v, init group: %v", nodeInitAction.Node.Name, initGroup)

	for _, item := range baseItemEnums {
		if _, ok := initGroup[item]; !ok {
			initGroup[item] = true
		}
	}

	return initGroup
}
