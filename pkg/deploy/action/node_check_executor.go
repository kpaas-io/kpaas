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
	desiredSystemManager              = "systemd"
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

func init() {
	RegisterExecutor(ActionTypeNodeCheck, new(nodeCheckExecutor))
}

type nodeCheckExecutor struct {
}

// due to items, ItemsCheckScripts exec remote scripts and return std, report, error
func ExecuteCheckScript(item check.ItemEnum, config *pb.NodeCheckConfig, checkItemReport *NodeCheckItem) (string, *NodeCheckItem, error) {

	checkItemReport = &NodeCheckItem{
		Name:        fmt.Sprintf("%v check", item),
		Description: fmt.Sprintf("%v check", item),
	}

	// create item operation
	checkItems := check.NewCheckOperations().CreateOperations(item)
	if checkItems == nil {
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.Reason = ItemErrEmpty
		checkItemReport.Err.Detail = ItemErrEmpty
		checkItemReport.Err.FixMethods = ItemHelperEmpty
	}

	// close ssh client
	defer checkItems.CloseSSH()

	// create operation commands for specific item
	op, err := checkItems.GetOperations(config)
	if err != nil {
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.Reason = ItemErrOperation
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = ItemHelperOperation
	}

	// exec operations commands
	stdErr, stdOut, err := op.Do()
	if err != nil {
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.Reason = ItemErrScript
		checkItemReport.Err.Detail = string(stdErr)
		checkItemReport.Err.FixMethods = ItemHelperScript
	}

	checkItemStdOut := string(stdOut)
	return checkItemStdOut, checkItemReport, nil
}

func newNodeCheckItem() *NodeCheckItem {

	return &NodeCheckItem{
		Status: ItemActionPending,
		Err:    &pb.Error{},
	}
}

// goroutine as executor for check docker
func CheckDockerExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemActionDoing
	comparedDockerVersion, checkItemReport, err := ExecuteCheckScript(check.Docker, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		checkItemReport.Status = ItemActionFailed
	}

	err = check.CheckDockerVersion(comparedDockerVersion, desiredDockerVersion, ">")
	if err != nil {
		checkItemReport.Err.Reason = "docker version too low"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please upgrade docker version to %v+", desiredDockerVersion)
	} else {
		checkItemReport.Status = ItemActionDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check CPU
func CheckCPUExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemActionDoing
	cpuCore, checkItemReport, err := ExecuteCheckScript(check.CPU, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		checkItemReport.Status = ItemActionFailed
	}

	err = check.CheckCPUNums(cpuCore, desiredCPUCore)
	if err != nil {
		checkItemReport.Err.Reason = "cpu cores not enough"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize cpu cores to %v", desiredCPUCore)
	} else {
		checkItemReport.Status = ItemActionDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check kernel
func CheckKernelExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemActionDoing
	kernelVersion, checkItemReport, err := ExecuteCheckScript(check.Kernel, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		checkItemReport.Status = ItemActionFailed
	}

	err = check.CheckKernelVersion(kernelVersion, desiredKernelVersion, ">")
	if err != nil {
		checkItemReport.Err.Reason = "kernel version too low"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize kernel version to %v", desiredKernelVersion)
	} else {
		checkItemReport.Status = ItemActionDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check memory
func CheckMemoryExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemActionDoing
	memoryCap, checkItemReport, err := ExecuteCheckScript(check.Memory, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		checkItemReport.Status = ItemActionFailed
	}

	err = check.CheckMemoryCapacity(memoryCap, desiredMemory)
	if err != nil {
		checkItemReport.Err.Reason = "memory capacity not enough"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize memory capacity to %v", desiredMemory)
	} else {
		checkItemReport.Status = ItemActionDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check disk
func CheckRootDiskExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemActionDoing
	rootDiskVolume, checkItemReport, err := ExecuteCheckScript(check.Disk, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		checkItemReport.Status = ItemActionFailed
	}

	err = check.CheckRootDiskVolume(rootDiskVolume, desiredRootDiskVolume)
	if err != nil {
		checkItemReport.Err.Reason = "root disk volume is not enough"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize root disk volume to %v", desiredRootDiskVolume)
	} else {
		checkItemReport.Status = ItemActionDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check distribution
func CheckDistributionExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemActionDoing
	disName, checkItemReport, err := ExecuteCheckScript(check.Distribution, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		checkItemReport.Status = ItemActionFailed
	}

	err = check.CheckSystemDistribution(disName)
	if err != nil {
		checkItemReport.Err.Reason = "system distribution is not supported"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please change suitable distribution to %v", systemDistributions)
	} else {
		checkItemReport.Status = ItemActionDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check system preference
func CheckSysPrefExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemActionDoing
	_, checkItemReport, err := ExecuteCheckScript(check.SystemPreference, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		checkItemReport.Err.Reason = "system preference is not supported"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.FixMethods = fmt.Sprint("please modify system preference")
	} else {
		checkItemReport.Status = ItemActionDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check system components
func CheckSysComponentExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemActionDoing
	systemManager, checkItemReport, err := ExecuteCheckScript(check.SystemComponent, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		checkItemReport.Status = ItemActionFailed
	}

	err = check.CheckSysComponent(systemManager, desiredSystemManager)
	if err != nil {
		checkItemReport.Err.Reason = "system component is not clear"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemActionFailed
		checkItemReport.Err.FixMethods = fmt.Sprint("please check system component is available")
	} else {
		checkItemReport.Status = ItemActionDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

func (a *nodeCheckExecutor) Execute(act Action) *pb.Error {
	nodeCheckAction, ok := act.(*NodeCheckAction)
	if !ok {
		return errOfTypeMismatched(new(NodeCheckAction), act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	var wg sync.WaitGroup

	logger.Debug("Start to execute node check action")

	// check docker, CPU, kernel, memory, disk, distribution, system preference
	wg.Add(8)
	go CheckDockerExecutor(nodeCheckAction, &wg)
	go CheckCPUExecutor(nodeCheckAction, &wg)
	go CheckKernelExecutor(nodeCheckAction, &wg)
	go CheckMemoryExecutor(nodeCheckAction, &wg)
	go CheckRootDiskExecutor(nodeCheckAction, &wg)
	go CheckDistributionExecutor(nodeCheckAction, &wg)
	go CheckSysPrefExecutor(nodeCheckAction, &wg)
	go CheckSysComponentExecutor(nodeCheckAction, &wg)
	wg.Wait()

	// If any of check item was failed, we should return an error
	failedItems := getFailedCheckItems(nodeCheckAction)
	if len(failedItems) > 0 {
		return &pb.Error{
			Reason: fmt.Sprintf("%d check item(s) failed", len(failedItems)),
			Detail: fmt.Sprintf("failed check item list: %v", failedItems),
		}
	}

	logger.Debug("Finish to execute node check action")
	return nil
}

func getFailedCheckItems(checkAction *NodeCheckAction) []string {
	var failedItemName []string
	for _, item := range checkAction.CheckItems {
		if item.Status != ItemActionDone {
			failedItemName = append(failedItemName, item.Name)
		}
	}
	return failedItemName
}
