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
	"sync"
)

type nodeInitExecutor struct{}

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

	// close ssh client
	initItem.CloseSSH()

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
	_, initItemReport, err := ExecuteInitScript(item, ncAction.Node, initItemReport)
	if err != nil {
		initItemReport.Status = ItemActionFailed
	}

	UpdateInitItems(ncAction, initItemReport)

	wg.Done()
}

func (a *nodeInitExecutor) Execute(act Action) error {
	nodeInitAction, ok := act.(*NodeInitAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be node init action, but is %T", act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	var wg sync.WaitGroup

	logger.Debug("Start to execute node init action")

	// init events include firewall, hostalias, hostname, network
	// route, swap, timezone, haproxy, keepalived

	itemEnums := []it.ItemEnum{it.Swap, it.Route, it.Network, it.Network, it.FireWall, it.TimeZone, it.HostName, it.HostAlias}
	for _, item := range itemEnums {
		wg.Add(1)
		go InitAsyncExecutor(item, nodeInitAction, &wg)
	}

	// TODO Other Init Items
	// 7. Haproxy
	// 8. Keepalived
	// 9. Install kubeadm kubectl kubelet

	wg.Wait()

	nodeInitAction.Status = ActionDone
	logger.Debug("Finish to execute node init action")
	return nil
}

// update init items with matching name
func UpdateInitItems(initAction *NodeInitAction, report *NodeInitItem) {

	initAction.Lock()
	defer initAction.Unlock()

	initAction.InitItems = append(initAction.InitItems, report)
}
