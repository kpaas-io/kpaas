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

package task

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	RegisterProcessor(TaskTypeDeploy, new(deployProcessor))
}

// deployProcessor implements the specific logic for the deploy task.
type deployProcessor struct {
}

// Spilt the task into one or more sub tasks
func (p *deployProcessor) SplitTask(t Task) error {
	deployTask, err := p.verifyTask(t)
	if err != nil {
		logrus.Errorf("Invalid task: %s", err)
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: t.GetName(),
	})

	logger.Debug("Start to split deploy task")

	// split task into subtask: init, deploy etcd, deploy master, deploy worker, deploy ingress
	var subTasks []Task

	// first collect all roles and their related nodes
	roles := p.groupByRole(deployTask.NodeConfigs)

	// create the init sub tasks with priority = 10
	initTask, err := p.createInitSubTask(deployTask, roles)
	if err != nil {
		err = fmt.Errorf("failed to create init sub tasks: %s", err)
		logger.Error(err)
		return err
	}
	subTasks = append(subTasks, initTask)

	// create the deploy etcd sub tasks with priority = 20
	if _, ok := roles[consts.NodeRoleEtcd]; ok {
		etcdTask, err := p.createDeploySubTask(consts.NodeRoleEtcd, deployTask, roles)
		if err != nil {
			err = fmt.Errorf("failed to create deploy etcd sub tasks: %s", err)
			logger.Error(err)
			return err
		}
		subTasks = append(subTasks, etcdTask)
	}

	// create the deploy master sub tasks with priority = 30
	if _, ok := roles[consts.NodeRoleMaster]; ok {
		masterTask, err := p.createDeploySubTask(consts.NodeRoleMaster, deployTask, roles)
		if err != nil {
			err = fmt.Errorf("failed to create deploy master sub tasks: %s", err)
			logger.Error(err)
			return err
		}
		subTasks = append(subTasks, masterTask)
	}

	// create the deploy worker sub tasks with priority = 40
	if _, ok := roles[consts.NodeRoleWorker]; ok {
		workerTask, err := p.createDeploySubTask(consts.NodeRoleWorker, deployTask, roles)
		if err != nil {
			err = fmt.Errorf("failed to create deploy worker sub tasks: %s", err)
			logger.Error(err)
			return err
		}
		subTasks = append(subTasks, workerTask)
	}

	// create the deploy ingress sub tasks with priority = 50
	if _, ok := roles[consts.NodeRoleIngress]; ok {
		ingressTask, err := p.createDeploySubTask(consts.NodeRoleIngress, deployTask, roles)
		if err != nil {
			err = fmt.Errorf("failed to create deploy ingress sub tasks: %s", err)
			logger.Error(err)
			return err
		}
		subTasks = append(subTasks, ingressTask)
	}

	deployTask.SubTasks = subTasks
	logger.Debugf("Finish to split deploy task: %d sub tasks", len(subTasks))

	return nil
}

// Verify if the task is valid.
func (p *deployProcessor) verifyTask(t Task) (*DeployTask, error) {
	if t == nil {
		return nil, consts.ErrEmptyTask
	}

	deployTask, ok := t.(*DeployTask)
	if !ok {
		return nil, fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, t)
	}

	if len(deployTask.NodeConfigs) == 0 {
		return nil, fmt.Errorf("nodeConfigs is empty")
	}

	return deployTask, nil
}

func (p *deployProcessor) groupByRole(cfgs []*pb.NodeDeployConfig) map[consts.NodeRole][]*pb.NodeDeployConfig {
	roles := make(map[consts.NodeRole][]*pb.NodeDeployConfig)
	for _, nodeCfg := range cfgs {
		nodeRoles := nodeCfg.GetRoles()
		for _, role := range nodeRoles {
			roleName := consts.NodeRole(role)
			roles[roleName] = append(roles[roleName], nodeCfg)
		}
	}
	return roles
}

func (p *deployProcessor) createInitSubTask(t *DeployTask, rn map[consts.NodeRole][]*pb.NodeDeployConfig) (Task, error) {
	// TODO
	return nil, nil
}

func (p *deployProcessor) createDeploySubTask(role consts.NodeRole, parent *DeployTask, rn map[consts.NodeRole][]*pb.NodeDeployConfig) (task Task, err error) {

	switch role {
	case consts.NodeRoleEtcd:
		config := &DeployEtcdTaskConfig{
			Nodes:           p.unwrapNodes(rn[role]),
			LogFileBasePath: parent.GetLogFilePath(),
			Priority:        int(Priorities[role]),
			Parent:          parent.GetName(),
		}
		// Use the role name as the task name for now.
		taskName := string(role)
		task, err = NewDeployEtcdTask(taskName, config)

	case consts.NodeRoleMaster:
		config := &DeployMasterTaskConfig{
			etcdNodes:       p.unwrapNodes(rn[consts.NodeRoleEtcd]),
			Nodes:           p.unwrapNodes(rn[role]),
			ClusterConfig:   parent.ClusterConfig,
			LogFileBasePath: parent.GetLogFilePath(),
			Priority:        int(Priorities[role]),
			Parent:          parent.GetName(),
		}
		// Use the role name as the task name for now.
		taskName := string(role)
		task, err = NewDeployMasterTask(taskName, config)

	case consts.NodeRoleWorker:

		config := &DeployWorkerTaskConfig{
			Nodes:           rn[consts.NodeRoleWorker],
			ClusterConfig:   parent.ClusterConfig,
			LogFileBasePath: parent.GetLogFilePath(),
			Priority:        int(Priorities[role]),
			Parent:          parent.GetName(),
		}

		// Use the role name as the task name for now.
		taskName := string(role)
		return NewDeployWorkerTask(taskName, config)
	default:
		err = fmt.Errorf("unrecognized role:%v", role)
	}

	return
}

func (p deployProcessor) unwrapNode(config *pb.NodeDeployConfig) *pb.Node {
	return config.GetNode()
}

func (p deployProcessor) unwrapNodes(nodeConfigs []*pb.NodeDeployConfig) []*pb.Node {

	nodes := make([]*pb.Node, 0, len(nodeConfigs))
	for _, nodeConfig := range nodeConfigs {
		nodes = append(nodes, p.unwrapNode(nodeConfig))
	}
	return nodes
}
