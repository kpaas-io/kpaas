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
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/check"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	desiredDockerVersion              = "18.09.0"
	desiredKernelVersion              = "4.19.46"
	desiredCPUCore            float64 = 8
	desiredMemoryByteBase     float64 = 16
	desiredMemory                     = desiredMemoryByteBase * operation.GiByteUnits
	desiredDiskVolumeByteBase float64 = 200
	desiredRootDiskVolume             = desiredDiskVolumeByteBase * operation.GiByteUnits

	ItemActionPending = "pending"
	ItemActionDoing   = "doing"
	ItemActionDone    = "done"
	ItemActionFailed  = "failed"

	ItemErrEmpty     = "empty parameter"
	ItemErrOperation = "failed to generate operations"
	ItemErrScript    = "invalid script"

	ItemHelperEmpty     = "please input suitable check item"
	ItemHelperOperation = "please check your operations"
	ItemHelperScript    = "please check your script"
)

var systemDistributions = [3]string{check.DistributionCentos, check.DistributionUbuntu, check.DistributionRHEL}

type nodeCheckExecutor struct {
}

// due to items, ItemsCheckScripts exec remote scripts and return std, report, error
func ExecuteCheckScript(item check.ItemEnum, config *pb.NodeCheckConfig, checkItemReport *nodeCheckItem) (string, *nodeCheckItem, error) {

	checkItemReport = &nodeCheckItem{
		name:        fmt.Sprintf("%v check", item),
		description: fmt.Sprintf("%v check", item),
	}

	checkItems := check.NewCheckOperations().CreateOperations(item)
	if checkItems == nil {
		checkItemReport.status = ItemActionFailed
		checkItemReport.err.Reason = ItemErrEmpty
		checkItemReport.err.Detail = ItemErrEmpty
		checkItemReport.err.FixMethods = ItemHelperEmpty
	}

	op, err := checkItems.GetOperations(config)
	if err != nil {
		checkItemReport.status = ItemActionFailed
		checkItemReport.err.Reason = ItemErrOperation
		checkItemReport.err.Detail = err.Error()
		checkItemReport.err.FixMethods = ItemHelperOperation
	}

	stdErr, stdOut, err := op.Do()
	if err != nil {
		checkItemReport.status = ItemActionFailed
		checkItemReport.err.Reason = ItemErrScript
		checkItemReport.err.Detail = string(stdErr)
		checkItemReport.err.FixMethods = ItemHelperScript
	}

	checkItemStdOut := string(stdOut)
	return checkItemStdOut, checkItemReport, nil
}

func CheckDockerExecutor(ncAction *nodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.status = ItemActionDoing
	comparedDockerVersion, checkItemReport, err := ExecuteCheckScript(check.Docker, ncAction.nodeCheckConfig, checkItemReport)
	UpdateCheckItems(ncAction, checkItemReport)
	if err != nil {
		checkItemReport.status = ItemActionFailed
	}

	err = check.CheckDockerVersion(comparedDockerVersion, desiredDockerVersion, ">")
	if err != nil {
		checkItemReport.err.Reason = "docker version too low"
		checkItemReport.err.Detail = err.Error()
		checkItemReport.status = ItemActionFailed
		checkItemReport.err.FixMethods = fmt.Sprintf("please upgrade docker version to %v+", desiredDockerVersion)
	} else {
		checkItemReport.status = ItemActionDone
	}
	UpdateCheckItems(ncAction, checkItemReport)

	wg.Done()
}

func newNodeCheckItem() *nodeCheckItem {

	return &nodeCheckItem{
		status: ItemActionPending,
		err:    &pb.Error{},
	}
}

func CheckCPUExecutor(ncAction *nodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.status = ItemActionDoing
	cpuCore, checkItemReport, err := ExecuteCheckScript(check.CPU, ncAction.nodeCheckConfig, checkItemReport)
	UpdateCheckItems(ncAction, checkItemReport)
	if err != nil {
		checkItemReport.status = ItemActionFailed
	}

	err = check.CheckCPUNums(cpuCore, desiredCPUCore)
	if err != nil {
		checkItemReport.err.Reason = "cpu cores not enough"
		checkItemReport.err.Detail = err.Error()
		checkItemReport.status = ItemActionFailed
		checkItemReport.err.FixMethods = fmt.Sprintf("please optimize cpu cores to %v", desiredCPUCore)
	} else {
		checkItemReport.status = ItemActionDone
	}
	UpdateCheckItems(ncAction, checkItemReport)

	wg.Done()
}

func CheckKernelExecutor(ncAction *nodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.status = ItemActionDoing
	kernelVersion, checkItemReport, err := ExecuteCheckScript(check.Kernel, ncAction.nodeCheckConfig, checkItemReport)
	UpdateCheckItems(ncAction, checkItemReport)
	if err != nil {
		checkItemReport.status = ItemActionFailed
	}

	err = check.CheckKernelVersion(kernelVersion, desiredKernelVersion, ">")
	if err != nil {
		checkItemReport.err.Reason = "kernel version too low"
		checkItemReport.err.Detail = err.Error()
		checkItemReport.status = ItemActionFailed
		checkItemReport.err.FixMethods = fmt.Sprintf("please optimize kernel version to %v", desiredKernelVersion)
	} else {
		checkItemReport.status = ItemActionDone
	}
	UpdateCheckItems(ncAction, checkItemReport)

	wg.Done()
}

func CheckMemoryExecutor(ncAction *nodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.status = ItemActionDoing
	memoryCap, checkItemReport, err := ExecuteCheckScript(check.Memory, ncAction.nodeCheckConfig, checkItemReport)
	UpdateCheckItems(ncAction, checkItemReport)
	if err != nil {
		checkItemReport.status = ItemActionFailed
	}

	err = check.CheckMemoryCapacity(memoryCap, desiredMemory)
	if err != nil {
		checkItemReport.err.Reason = "memory capacity not enough"
		checkItemReport.err.Detail = err.Error()
		checkItemReport.status = ItemActionFailed
		checkItemReport.err.FixMethods = fmt.Sprintf("please optimize memory capacity to %v", desiredMemory)
	} else {
		checkItemReport.status = ItemActionDone
	}
	UpdateCheckItems(ncAction, checkItemReport)

	wg.Done()
}

func CheckRootDiskExecutor(ncAction *nodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.status = ItemActionDoing
	rootDiskVolume, checkItemReport, err := ExecuteCheckScript(check.Disk, ncAction.nodeCheckConfig, checkItemReport)
	UpdateCheckItems(ncAction, checkItemReport)
	if err != nil {
		checkItemReport.status = ItemActionFailed
	}

	err = check.CheckRootDiskVolume(rootDiskVolume, desiredRootDiskVolume)
	if err != nil {
		checkItemReport.err.Reason = "root disk volume is not enough"
		checkItemReport.err.Detail = err.Error()
		checkItemReport.status = ItemActionFailed
		checkItemReport.err.FixMethods = fmt.Sprintf("please optimize root disk volume to %v", desiredRootDiskVolume)
	} else {
		checkItemReport.status = ItemActionDone
	}
	UpdateCheckItems(ncAction, checkItemReport)

	wg.Done()
}

func CheckDistributionExecutor(ncAction *nodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.status = ItemActionDoing
	disName, checkItemReport, err := ExecuteCheckScript(check.Distribution, ncAction.nodeCheckConfig, checkItemReport)
	UpdateCheckItems(ncAction, checkItemReport)
	if err != nil {
		checkItemReport.status = ItemActionFailed
	}

	err = check.CheckSystemDistribution(disName)
	if err != nil {
		checkItemReport.err.Reason = "system distribution is not supported"
		checkItemReport.err.Detail = err.Error()
		checkItemReport.status = ItemActionFailed
		checkItemReport.err.FixMethods = fmt.Sprintf("please change suitable distribution to %v", systemDistributions)
	} else {
		checkItemReport.status = ItemActionDone
	}
	UpdateCheckItems(ncAction, checkItemReport)

	wg.Done()
}

func (a *nodeCheckExecutor) Execute(act Action) error {
	nodeCheckAction, ok := act.(*nodeCheckAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be node check action, but is %T", act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	var wg sync.WaitGroup

	logger.Debug("Start to execute node check action")

	// check docker, CPU, kernel, memory, disk, distribution
	wg.Add(6)
	go CheckDockerExecutor(nodeCheckAction, &wg)
	go CheckCPUExecutor(nodeCheckAction, &wg)
	go CheckKernelExecutor(nodeCheckAction, &wg)
	go CheckMemoryExecutor(nodeCheckAction, &wg)
	go CheckRootDiskExecutor(nodeCheckAction, &wg)
	go CheckDistributionExecutor(nodeCheckAction, &wg)
	wg.Wait()

	nodeCheckAction.status = ActionDone
	logger.Debug("Finish to execute node check action")
	return nil
}

// update check items with matching name
func UpdateCheckItems(checkAction *nodeCheckAction, report *nodeCheckItem) {

	checkAction.Lock()
	defer checkAction.Unlock()

	updatedFlag := false

	for _, item := range checkAction.checkItems {
		if item.name == report.name {
			updatedFlag = true
			item.err = report.err
			item.status = report.status
			item.description = report.description
		}
	}

	if updatedFlag == false {
		checkAction.checkItems = append(checkAction.checkItems, report)
	}
}
