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
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/check"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	desiredDockerVersion              = "18.09.0"
	desiredKernelVersion              = "4.19.46"
	desiredCPUCore            float64 = 8
	desiredMemoryBase         float64 = 16
	desiredMemory                     = desiredMemoryBase * operation.GiByteUnits
	desiredRootDiskVolumeBase float64 = 200
	desiredRootDiskVolume             = desiredRootDiskVolumeBase * operation.GiByteUnits
)

var systemDistributions = [3]string{"centos", "ubuntu", "rhel"}

type nodeCheckExecutor struct {
}

func ItemsCheckScripts(items string, config *pb.NodeCheckConfig) (string, *nodeCheckItem, error) {

	var (
		reason    string
		detail    string
		status    nodeCheckItemStatus
		fixMethod string
	)

	checkItemReport := &nodeCheckItem{
		name:        fmt.Sprintf("%v check", items),
		description: fmt.Sprintf("%v check", items),
		status:      status,
		err: &pb.Error{
			Reason:     reason,
			Detail:     detail,
			FixMethods: fixMethod,
		},
	}

	checkItems := check.NewCheckOperations().CreateOperations(items)
	op, err := checkItems.GetOperations(config)
	if err != nil {
		return "", checkItemReport, fmt.Errorf("failed to create %v check operation, error: %v", items, err)
	}

	stdErr, stdOut, err := op.Do()
	if err != nil {
		checkItemReport.err.Reason = fmt.Sprintf("run check %v command failed", items)
		checkItemReport.err.Detail = string(stdErr)
		checkItemReport.status = nodeCheckItemFailed
		checkItemReport.err.FixMethods = "please check your scripts"
		return "", checkItemReport, fmt.Errorf("failed to run check %v scripts", items)
	}

	checkItemStdOut := string(stdOut[:])
	return checkItemStdOut, checkItemReport, nil
}

func (a *nodeCheckExecutor) Execute(act Action) error {
	nodeCheckAction, ok := act.(*nodeCheckAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be node check action, but is %T", act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	logger.Debug("Start to execute node check action")

	// check docker
	comparedDockerVersion, report, err := ItemsCheckScripts("docker", nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return err
	}

	err = check.CheckDockerVersion(comparedDockerVersion, desiredDockerVersion, ">")
	if err != nil {
		report.err.Reason = "docker version too low"
		report.err.Detail = string(report.err.Detail)
		report.status = nodeCheckItemFailed
		report.err.FixMethods = fmt.Sprintf("please upgrade docker version to %v+", desiredDockerVersion)
		return err
	}

	report.status = nodeCheckItemSucessful
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, report)

	//// check CPU
	cpuCore, report, err := ItemsCheckScripts("cpu", nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return err
	}

	err = check.CheckCPUNums(cpuCore, desiredCPUCore)
	if err != nil {
		report.err.Reason = "cpu cores not enough"
		report.err.Detail = string(report.err.Detail)
		report.status = nodeCheckItemFailed
		report.err.FixMethods = fmt.Sprintf("please optimize cpu cores to %v", desiredCPUCore)
		return err
	}

	report.status = nodeCheckItemSucessful
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, report)

	// check kernel version
	kernelVersion, report, err := ItemsCheckScripts("kernel", nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return err
	}

	err = check.CheckKernelVersion(kernelVersion, desiredKernelVersion, ">")
	if err != nil {
		report.err.Reason = "kernel version too low"
		report.err.Detail = string(report.err.Detail)
		report.status = nodeCheckItemFailed
		report.err.FixMethods = fmt.Sprintf("please optimize kernel version to %v", desiredKernelVersion)
		return err
	}

	report.status = nodeCheckItemSucessful
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, report)

	// check memory capacity
	memoryCap, report, err := ItemsCheckScripts("memory", nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return err
	}

	err = check.CheckMemoryCapacity(memoryCap, desiredMemory)
	if err != nil {
		report.err.Reason = "memory capacity not enough"
		report.err.Detail = string(report.err.Detail)
		report.status = nodeCheckItemFailed
		report.err.FixMethods = fmt.Sprintf("please optimize memory capacity to %v", desiredMemory)
		return err
	}

	report.status = nodeCheckItemSucessful
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, report)

	// check root disk volume
	rootDiskVolume, report, err := ItemsCheckScripts("disk", nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return err
	}

	err = check.CheckRootDiskVolume(rootDiskVolume, desiredRootDiskVolume)
	if err != nil {
		report.err.Reason = "root disk volume is not enough"
		report.err.Detail = string(report.err.Detail)
		report.status = nodeCheckItemFailed
		report.err.FixMethods = fmt.Sprintf("please optimize root disk volume to %v", desiredRootDiskVolume)
		return err
	}

	report.status = nodeCheckItemSucessful
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, report)

	// check system distribution
	disName, report, err := ItemsCheckScripts("distribution", nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return err
	}

	err = check.CheckSystemDistribution(disName)
	if err != nil {
		report.err.Reason = "system distribution is not supported"
		report.err.Detail = string(report.err.Detail)
		report.status = nodeCheckItemFailed
		report.err.FixMethods = fmt.Sprintf("please change suitable distribution to %v", systemDistributions)
		return err
	}

	report.status = nodeCheckItemSucessful
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, report)

	nodeCheckAction.status = ActionDone
	logger.Debug("Finish to execute node check action")
	return nil
}
