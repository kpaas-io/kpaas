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

	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/check"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// constant value for check
const (
	desiredDockerVersion              = "18.09.0"
	desiredKernelVersion              = "4.19.46"
	desiredSystemManager              = "systemd"
	desiredCPUCore            float64 = 4
	desiredMemoryByteBase     float64 = 8
	desiredMemory                     = desiredMemoryByteBase * operation.GiByteUnits
	desiredDiskVolumeByteBase float64 = 50
	desiredRootDiskVolume             = desiredDiskVolumeByteBase * operation.GiByteUnits

	ItemErrEmpty     = "empty parameter"
	ItemErrOperation = "failed to generate operations"
	ItemErrScript    = "invalid script"

	ItemHelperEmpty     = "please input suitable check item"
	ItemHelperOperation = "please check your operations"
	ItemHelperScript    = "please check your script"

	CheckPassed = "check passed"
	CheckFailed = "check failed"
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
		Name:        fmt.Sprintf("check %v", item),
		Description: fmt.Sprintf("检查 %v 环境", item),
	}

	// create item operation
	checkItems := check.NewCheckOperations().CreateOperations(item)
	if checkItems == nil {
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = ItemErrEmpty
		checkItemReport.Err.Detail = ItemErrEmpty
		checkItemReport.Err.FixMethods = ItemHelperEmpty
		return "", checkItemReport, fmt.Errorf("fail to construct %v operation", item)
	}

	// close ssh client
	defer checkItems.CloseSSH()

	// create operation commands for specific item
	op, err := checkItems.GetOperations(config)
	if err != nil {
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = ItemErrOperation
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = ItemHelperOperation
		return "", checkItemReport, fmt.Errorf("fail to construct %v commands", item)
	}

	// exec operations commands
	stdOut, stdErr, err := op.Do()
	if err != nil {
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = ItemErrScript
		checkItemReport.Err.Detail = string(stdErr)
		checkItemReport.Err.FixMethods = ItemHelperScript
		return "", checkItemReport, fmt.Errorf("fail to run %v commands", item)
	}

	checkItemStdOut := strings.Trim(string(stdOut), "\n")
	return checkItemStdOut, checkItemReport, nil
}

func newNodeCheckItem() *NodeCheckItem {

	return &NodeCheckItem{
		Status: ItemPending,
		Err:    &pb.Error{},
	}
}

// goroutine as executor for check docker
func CheckDockerExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "docker",
	})

	logger.Debug("Start to execute check docker")

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemDoing
	comparedDockerVersion, checkItemReport, err := ExecuteCheckScript(check.Docker, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check docker failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	err = check.CheckDockerVersion(comparedDockerVersion, desiredDockerVersion, ">")
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "docker version too low"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please upgrade docker version to %v+", desiredDockerVersion)
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check CPU
func CheckCPUExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "cpu",
	})

	logrus.Debug("Start to execute check cpu")

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemDoing
	cpuCore, checkItemReport, err := ExecuteCheckScript(check.CPU, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check cpu failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	err = check.CheckCPUNums(cpuCore, desiredCPUCore)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "cpu cores not enough"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize cpu cores to %v", desiredCPUCore)
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check kernel
func CheckKernelExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "kernel",
	})

	logrus.Debug("Start to execute check kernel")

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemDoing
	kernelVersion, checkItemReport, err := ExecuteCheckScript(check.Kernel, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check kernel failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	err = check.CheckKernelVersion(kernelVersion, desiredKernelVersion, ">")
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "kernel version too low"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize kernel version to %v", desiredKernelVersion)
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check memory
func CheckMemoryExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "memory",
	})

	logrus.Debug("Start to execute check memory")

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemDoing
	memoryCap, checkItemReport, err := ExecuteCheckScript(check.Memory, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check memory failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	err = check.CheckMemoryCapacity(memoryCap, desiredMemory)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "memory capacity not enough"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize memory capacity to %v", deploy.ReturnWithUnit(desiredMemory))
	} else {
		logger.Debug(CheckPassed)
		logrus.Debug("memory check passed")
		checkItemReport.Status = ItemDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check disk
func CheckRootDiskExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "root disk",
	})

	logrus.Debug("Start to execute check disk volume")

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemDoing
	rootDiskVolume, checkItemReport, err := ExecuteCheckScript(check.Disk, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check root disk failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	err = check.CheckRootDiskVolume(rootDiskVolume, desiredRootDiskVolume)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "root disk volume is not enough"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize root disk volume to %s", deploy.ReturnWithUnit(desiredRootDiskVolume))
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check distribution
func CheckDistributionExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "distro",
	})

	logrus.Debug("Start to execute check distro")

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemDoing
	disName, checkItemReport, err := ExecuteCheckScript(check.Distribution, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check distro failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	disName = strings.Trim(disName, "\"")
	err = check.CheckSystemDistribution(disName)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "system distribution is not supported"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please change suitable distribution to %v", systemDistributions)
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check system preference
func CheckSysPrefExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "system preference",
	})

	logrus.Debug("Start to execute check system preference")

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemDoing
	_, checkItemReport, err := ExecuteCheckScript(check.SystemPreference, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check system preference failed, err: %v", err)
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "system preference is not supported"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemFailed
		checkItemReport.Err.FixMethods = fmt.Sprint("please modify system preference")
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for check system manager
func CheckSysManagerExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "system manager",
	})

	logrus.Debug("Start to execute check system manager")

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemDoing
	systemManager, checkItemReport, err := ExecuteCheckScript(check.SystemManager, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check system manager failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	err = check.CheckSystemManager(systemManager, desiredSystemManager)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "system manager is not clear"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemFailed
		checkItemReport.Err.FixMethods = fmt.Sprint("please check system manager is systemd")
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ncAction.Lock()
	defer ncAction.Unlock()
	ncAction.CheckItems = append(ncAction.CheckItems, checkItemReport)

	wg.Done()
}

// goroutine as executor for port occupied check
func CheckPortOccupiedExecutor(ncAction *NodeCheckAction, wg *sync.WaitGroup) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "port occupied",
	})

	logrus.Debug("Start to execute check port occupied")

	checkItemReport := newNodeCheckItem()
	checkItemReport.Status = ItemDoing
	portOccupied, checkItemReport, err := ExecuteCheckScript(check.PortOccupied, ncAction.NodeCheckConfig, checkItemReport)

	// trim can be done whatever error occurs
	portOccupied = strings.TrimRight(portOccupied, ",")
	if err != nil {
		logger.Errorf("check port occupied failed, err: %v, occupied port: %v", err, portOccupied)
		checkItemReport.Status = ItemFailed
	}

	portResult, err := check.CheckPortOccupied(portOccupied)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "port occupied check is failed"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Status = ItemFailed
		checkItemReport.Err.FixMethods = fmt.Sprintf("please solve port occupied problem, occupied port: %s", portResult)
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
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
	// system manager, port occupied
	wg.Add(9)
	go CheckDockerExecutor(nodeCheckAction, &wg)
	go CheckCPUExecutor(nodeCheckAction, &wg)
	go CheckKernelExecutor(nodeCheckAction, &wg)
	go CheckMemoryExecutor(nodeCheckAction, &wg)
	go CheckRootDiskExecutor(nodeCheckAction, &wg)
	go CheckDistributionExecutor(nodeCheckAction, &wg)
	go CheckSysPrefExecutor(nodeCheckAction, &wg)
	go CheckSysManagerExecutor(nodeCheckAction, &wg)
	go CheckPortOccupiedExecutor(nodeCheckAction, &wg)
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
		if item.Status != ItemDone {
			failedItemName = append(failedItemName, item.Name)
		}
	}
	return failedItemName
}
