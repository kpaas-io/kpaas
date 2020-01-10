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
	"math"
	"strings"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/check"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// constant value for check
const (
	desiredDockerVersion = "18.09.0"
	desiredKernelVersion = "4.19.46"
	desiredSystemManager = "systemd"

	// CPU factor
	desiredEtcdCPUCore    float64 = 4
	desiredMasterCPUCore  float64 = 4
	desiredWorkerCPUCore  float64 = 4
	desiredIngressCPUCore float64 = 4
	lowestCPUCore         float64 = 4

	// Memory factor
	desiredEtcdMemoryByteBase    float64 = 8
	desiredMasterMemoryByteBase  float64 = 8
	desiredWorkerMemoryByteBase  float64 = 8
	desiredIngressMemoryByteBase float64 = 8
	lowestMemoryByteBase         float64 = 8

	// Root Disk factor
	desiredEtcdDiskVolumeByteBase    float64 = 50
	desiredMasterDiskVolumeByteBase  float64 = 30
	desiredWorkerDiskVolumeByteBase  float64 = 30
	desiredIngressDiskVolumeByteBase float64 = 10
	lowestDiskVolumeByteBase         float64 = 50

	ItemErrEmpty     = "empty parameter"
	ItemErrOperation = "failed to build or run script"
	ItemErrScript    = "invalid script"

	ItemHelperEmpty     = "please input suitable check item"
	ItemHelperOperation = "please check your operation or script"

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

	checkItemReport = newNodeCheckItem(item)

	// create item operation
	checkItems := check.NewCheckOperations().CreateOperations(item)
	if checkItems == nil {
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = ItemErrEmpty
		checkItemReport.Err.Detail = ItemErrEmpty
		checkItemReport.Err.FixMethods = ItemHelperEmpty
		return "", checkItemReport, fmt.Errorf("fail to construct check %v operation", item)
	}

	// create command and run on remote node
	stdOut, stdErr, err := checkItems.RunCommands(config)
	if err != nil {
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = ItemErrOperation
		checkItemReport.Err.Detail = fmt.Sprintf("stdErr: %v, err: %v", stdErr, err.Error())
		checkItemReport.Err.FixMethods = ItemHelperOperation
		return "", checkItemReport, fmt.Errorf("fail to run check %v scripts", item)
	}

	checkItemStdOut := strings.Trim(string(stdOut), "\n")
	return checkItemStdOut, checkItemReport, nil
}

func newNodeCheckItem(item check.ItemEnum) *NodeCheckItem {

	return &NodeCheckItem{
		Status:      ItemDoing,
		Name:        fmt.Sprintf("check %v", item),
		Description: fmt.Sprintf("检查 %v 环境", item),
	}
}

// goroutine as executor for check docker
func CheckDockerExecutor(ncAction *NodeCheckAction, ch chan<- *NodeCheckItem) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "docker",
	})

	logger.Debug("Start to execute check docker")

	checkItemReport := newNodeCheckItem(check.Docker)

	comparedDockerVersion, checkItemReport, err := ExecuteCheckScript(check.Docker, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check docker failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	err = check.CheckDockerVersion(comparedDockerVersion, desiredDockerVersion, ">")
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "docker version too low"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = fmt.Sprintf("please upgrade docker version to %v+", desiredDockerVersion)
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ch <- checkItemReport
}

// goroutine as executor for check CPU
func CheckCPUExecutor(ncAction *NodeCheckAction, ch chan<- *NodeCheckItem) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "cpu",
	})

	logrus.Debug("Start to execute check cpu")

	checkItemReport := newNodeCheckItem(check.CPU)

	cpuCore, checkItemReport, err := ExecuteCheckScript(check.CPU, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check cpu failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	var desiredCPUCore float64
	for _, role := range ncAction.NodeCheckConfig.Roles {
		switch role {
		case string(constant.MachineRoleMaster):
			desiredCPUCore = math.Max(desiredCPUCore, desiredMasterCPUCore)
		case string(constant.MachineRoleWorker):
			desiredCPUCore = math.Max(desiredCPUCore, desiredWorkerCPUCore)
		case string(constant.MachineRoleEtcd):
			desiredCPUCore = math.Max(desiredCPUCore, desiredEtcdCPUCore)
		case string(constant.MachineRoleIngress):
			desiredCPUCore = math.Max(desiredCPUCore, desiredIngressCPUCore)
		}
	}

	// compare with lowest standard
	desiredCPUCore = math.Max(desiredCPUCore, lowestCPUCore)

	err = check.CheckCPUNums(cpuCore, desiredCPUCore)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "cpu cores not enough"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize cpu cores to %v", desiredCPUCore)
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ch <- checkItemReport
}

// goroutine as executor for check kernel
func CheckKernelExecutor(ncAction *NodeCheckAction, ch chan<- *NodeCheckItem) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "kernel",
	})

	logrus.Debug("Start to execute check kernel")

	checkItemReport := newNodeCheckItem(check.Kernel)

	kernelVersion, checkItemReport, err := ExecuteCheckScript(check.Kernel, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check kernel failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	err = check.CheckKernelVersion(kernelVersion, desiredKernelVersion, ">")
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "kernel version too low"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize kernel version to %v", desiredKernelVersion)
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ch <- checkItemReport
}

// goroutine as executor for check memory
func CheckMemoryExecutor(ncAction *NodeCheckAction, ch chan<- *NodeCheckItem) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "memory",
	})

	logrus.Debug("Start to execute check memory")

	checkItemReport := newNodeCheckItem(check.Memory)

	memoryCap, checkItemReport, err := ExecuteCheckScript(check.Memory, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check memory failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	var desiredMemory float64
	for _, role := range ncAction.NodeCheckConfig.Roles {
		switch role {
		case string(constant.MachineRoleMaster):
			desiredMemory = math.Max(desiredMemory, desiredMasterMemoryByteBase)
		case string(constant.MachineRoleWorker):
			desiredMemory = math.Max(desiredMemory, desiredWorkerMemoryByteBase)
		case string(constant.MachineRoleEtcd):
			desiredMemory = math.Max(desiredMemory, desiredEtcdMemoryByteBase)
		case string(constant.MachineRoleIngress):
			desiredMemory = math.Max(desiredMemory, desiredIngressMemoryByteBase)
		}
	}

	// compare with lowest standard
	desiredMemory = math.Max(desiredMemory, lowestMemoryByteBase)
	desiredMemory = desiredMemory * operation.GiByteUnits

	err = check.CheckMemoryCapacity(memoryCap, desiredMemory)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "memory capacity not enough"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize memory capacity to %v", deploy.ReturnWithUnit(desiredMemory))
	} else {
		logger.Debug(CheckPassed)
		logrus.Debug("memory check passed")
		checkItemReport.Status = ItemDone
	}

	ch <- checkItemReport
}

// goroutine as executor for check disk
func CheckRootDiskExecutor(ncAction *NodeCheckAction, ch chan<- *NodeCheckItem) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "root disk",
	})

	logrus.Debug("Start to execute check disk volume")

	checkItemReport := newNodeCheckItem(check.Disk)

	rootDiskVolume, checkItemReport, err := ExecuteCheckScript(check.Disk, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check root disk failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	var desiredRootDiskVolume float64
	for _, role := range ncAction.NodeCheckConfig.Roles {
		switch role {
		case string(constant.MachineRoleMaster):
			desiredRootDiskVolume = math.Max(desiredRootDiskVolume, desiredMasterDiskVolumeByteBase)
		case string(constant.MachineRoleWorker):
			desiredRootDiskVolume = math.Max(desiredRootDiskVolume, desiredWorkerDiskVolumeByteBase)
		case string(constant.MachineRoleEtcd):
			desiredRootDiskVolume = math.Max(desiredRootDiskVolume, desiredEtcdDiskVolumeByteBase)
		case string(constant.MachineRoleIngress):
			desiredRootDiskVolume = math.Max(desiredRootDiskVolume, desiredIngressDiskVolumeByteBase)
		}
	}

	// compare with lowest standard
	desiredRootDiskVolume = math.Max(desiredRootDiskVolume, lowestDiskVolumeByteBase)
	desiredRootDiskVolume = desiredRootDiskVolume * operation.GiByteUnits

	err = check.CheckRootDiskVolume(rootDiskVolume, desiredRootDiskVolume)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "root disk volume is not enough"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = fmt.Sprintf("please optimize root disk volume to %s", deploy.ReturnWithUnit(desiredRootDiskVolume))
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ch <- checkItemReport
}

// goroutine as executor for check distribution
func CheckDistributionExecutor(ncAction *NodeCheckAction, ch chan<- *NodeCheckItem) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "distro",
	})

	logrus.Debug("Start to execute check distro")

	checkItemReport := newNodeCheckItem(check.Distribution)

	disName, checkItemReport, err := ExecuteCheckScript(check.Distribution, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check distro failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	disName = strings.Trim(disName, "\"")
	err = check.CheckSystemDistribution(disName)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "system distribution is not supported"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = fmt.Sprintf("please change suitable distribution to %v", systemDistributions)
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ch <- checkItemReport
}

// goroutine as executor for check system preference
func CheckSysPrefExecutor(ncAction *NodeCheckAction, ch chan<- *NodeCheckItem) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "system preference",
	})

	logrus.Debug("Start to execute check system preference")

	checkItemReport := newNodeCheckItem(check.SystemPreference)

	_, checkItemReport, err := ExecuteCheckScript(check.SystemPreference, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "system preference is not supported"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = fmt.Sprint("please modify system preference")
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ch <- checkItemReport
}

// goroutine as executor for check system manager
func CheckSysManagerExecutor(ncAction *NodeCheckAction, ch chan<- *NodeCheckItem) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "system manager",
	})

	logrus.Debug("Start to execute check system manager")

	checkItemReport := newNodeCheckItem(check.SystemManager)

	systemManager, checkItemReport, err := ExecuteCheckScript(check.SystemManager, ncAction.NodeCheckConfig, checkItemReport)
	if err != nil {
		logger.Errorf("check system manager failed, err: %v", err)
		checkItemReport.Status = ItemFailed
	}

	err = check.CheckSystemManager(systemManager, desiredSystemManager)
	if err != nil {
		logger.Debugf("%v: %v", CheckFailed, err)
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "system manager is not clear"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = fmt.Sprint("please check system manager is systemd")
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ch <- checkItemReport
}

// goroutine as executor for port occupied check
func CheckPortOccupiedExecutor(ncAction *NodeCheckAction, ch chan<- *NodeCheckItem) {

	logger := logrus.WithFields(logrus.Fields{
		"node":       ncAction.Node.Name,
		"check_item": "port occupied",
	})

	logrus.Debug("Start to execute check port occupied")

	checkItemReport := newNodeCheckItem(check.PortOccupied)

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
		checkItemReport.Status = ItemFailed
		checkItemReport.Err = new(pb.Error)
		checkItemReport.Err.Reason = "port occupied check failed"
		checkItemReport.Err.Detail = err.Error()
		checkItemReport.Err.FixMethods = fmt.Sprintf("please close the process which occupied port: %v", portResult)
	} else {
		logger.Debug(CheckPassed)
		checkItemReport.Status = ItemDone
	}

	ch <- checkItemReport
}

func (a *nodeCheckExecutor) Execute(act Action) *pb.Error {
	nodeCheckAction, ok := act.(*NodeCheckAction)
	if !ok {
		return errOfTypeMismatched(new(NodeCheckAction), act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	logger.Debug("Start to execute node check action")

	// make enough length of check items
	channel := make(chan *NodeCheckItem, 9)

	// check docker, CPU, kernel, memory, disk, distribution, system preference
	// system manager, port occupied
	go CheckDockerExecutor(nodeCheckAction, channel)
	go CheckCPUExecutor(nodeCheckAction, channel)
	go CheckKernelExecutor(nodeCheckAction, channel)
	go CheckMemoryExecutor(nodeCheckAction, channel)
	go CheckRootDiskExecutor(nodeCheckAction, channel)
	go CheckDistributionExecutor(nodeCheckAction, channel)
	go CheckSysPrefExecutor(nodeCheckAction, channel)
	go CheckSysManagerExecutor(nodeCheckAction, channel)
	go CheckPortOccupiedExecutor(nodeCheckAction, channel)

	for report := range channel {
		nodeCheckAction.CheckItems = append(nodeCheckAction.CheckItems, report)

		if len(nodeCheckAction.CheckItems) == 9 {
			break
		}
	}

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
