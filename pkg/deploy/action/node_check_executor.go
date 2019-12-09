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
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/check/docker"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/check/system"
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

func (a *nodeCheckExecutor) Execute(act Action) error {
	nodeCheckAction, ok := act.(*nodeCheckAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be node check action, but is %T", act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	logger.Debug("Start to execute node check action")

	var (
		reason    string
		detail    string
		status    nodeCheckItemStatus
		fixmethod string
	)

	// check docker
	op, err := docker.NewCheckDockerOperation(nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return fmt.Errorf("failed to create docker check operation, error: %v", err)
	}

	stdErr, stdOut, err := op.Do()
	if err != nil {
		reason = "run check docker command failed"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = "please check your scripts"
		return err
	}

	comparedDockerVersion := string(stdOut[:])
	err = docker.CheckDockerVersion(comparedDockerVersion, desiredDockerVersion, ">")
	if err != nil {
		reason = "docker version not enough"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = fmt.Sprintf("please upgrade docker version to %v+", desiredDockerVersion)
		return err
	}

	status = nodeCheckItemSucessful

	dockerVersionItem := &nodeCheckItem{
		name:        "docker version check",
		description: "docker version check",
		status:      status,
		err: &pb.Error{
			Reason:     reason,
			Detail:     detail,
			FixMethods: fixmethod,
		},
	}
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, dockerVersionItem)

	// check CPU
	op, err = system.NewCheckCPUOperation(nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return fmt.Errorf("failed to create check cpu operation, error: %v", err)
	}

	stdErr, stdOut, err = op.Do()
	if err != nil {
		reason = "run check cpu command failed"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = "please check your scripts"
		return err
	}

	cpuCore := string(stdOut[:])
	err = system.CheckCPUNums(cpuCore, desiredCPUCore)
	if err != nil {
		reason = "cpu cores not enough"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = fmt.Sprintf("please optimize cpu cores to %v", desiredCPUCore)
		return err
	}

	status = nodeCheckItemSucessful

	cpuCoreItem := &nodeCheckItem{
		name:        "cpu core check",
		description: "cpu core check",
		status:      status,
		err: &pb.Error{
			Reason:     reason,
			Detail:     detail,
			FixMethods: fixmethod,
		},
	}
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, cpuCoreItem)

	// check kernel version
	op, err = system.NewCheckKernelOperation(nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return fmt.Errorf("failed to create check kernel operation, error: %v", err)
	}

	stdErr, stdOut, err = op.Do()
	if err != nil {
		reason = "run check kernel command failed"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = "please check your scripts"
		return err
	}

	kernelVersion := string(stdOut[:])
	err = system.CheckKernelVersion(kernelVersion, desiredKernelVersion, ">")
	if err != nil {
		reason = "kernel version not enough"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = fmt.Sprintf("please upgrade your kernel version to %v", desiredKernelVersion)
		return err
	}

	status = nodeCheckItemSucessful

	kernelItem := &nodeCheckItem{
		name:        "kernel version check",
		description: "kernel version check",
		status:      status,
		err: &pb.Error{
			Reason:     reason,
			Detail:     detail,
			FixMethods: fixmethod,
		},
	}
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, kernelItem)

	// check memory capacity
	op, err = system.NewCheckMemoryOperation(nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return fmt.Errorf("failed to create check memory operation, error: %v", err)
	}

	stdErr, stdOut, err = op.Do()
	if err != nil {
		reason = "run check memory command failed"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = "please check your scripts"
		return err
	}

	memoryCap := string(stdOut[:])
	err = system.CheckMemoryCapacity(memoryCap, desiredMemory)
	if err != nil {
		reason = "memory capacity not enough"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = fmt.Sprintf("please optimize your memory capacity to %v", desiredMemory)
		return err
	}

	status = nodeCheckItemSucessful

	memoryItem := &nodeCheckItem{
		name:        "memory capacity check",
		description: "memory capacity check",
		status:      status,
		err: &pb.Error{
			Reason:     reason,
			Detail:     detail,
			FixMethods: fixmethod,
		},
	}
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, memoryItem)

	// check root disk volume
	op, err = system.NewCheckRootDiskOperation(nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return fmt.Errorf("failed to create check memory operation, error: %v", err)
	}

	stdErr, stdOut, err = op.Do()
	if err != nil {
		reason = "run check root disk command failed"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = "please check your scripts"
		return err
	}

	rootDiskVolume := string(stdOut[:])
	err = system.CheckRootDiskVolume(rootDiskVolume, desiredRootDiskVolume)
	if err != nil {
		reason = "root disk volume not enough"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = fmt.Sprintf("please optimize your root disk volume to %v", desiredRootDiskVolume)
		return err
	}

	status = nodeCheckItemSucessful

	rootDiskItem := &nodeCheckItem{
		name:        "root disk volume check",
		description: "root disk volume check",
		status:      status,
		err: &pb.Error{
			Reason:     reason,
			Detail:     detail,
			FixMethods: fixmethod,
		},
	}
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, rootDiskItem)

	// check system distribution
	op, err = system.NewCheckDistributionOperation(nodeCheckAction.nodeCheckConfig)
	if err != nil {
		return fmt.Errorf("failed to create check system distribution, error: %v", err)
	}

	stdErr, stdOut, err = op.Do()
	if err != nil {
		reason = "run check system distribution command failed"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = "please check your scripts"
		return err
	}

	disName := string(stdOut[:])
	err = system.CheckSystemDistribution(disName)
	if err != nil {
		reason = "system distribution not supported"
		detail = string(stdErr)
		status = nodeCheckItemFailed
		fixmethod = fmt.Sprintf("please change your system distribution, supported: %v", systemDistributions)
		return err
	}

	status = nodeCheckItemSucessful

	systemDistributionItem := &nodeCheckItem{
		name:        "system distribution check",
		description: "system distribution check",
		status:      status,
		err: &pb.Error{
			Reason:     reason,
			Detail:     detail,
			FixMethods: fixmethod,
		},
	}
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, systemDistributionItem)

	nodeCheckAction.status = ActionDone
	logger.Debug("Finish to execute node check action")
	return nil
}
