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
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/copycerts"

	"github.com/kpaas-io/kpaas/pkg/constant"
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
		consts.LogFieldTask: t.GetName(),
	})

	logger.Debug("Start to split deploy task")

	// split task into subtask: init, deploy etcd, deploy master, deploy worker, deploy ingress
	var subTasks []Task

	// first collect all roles and their related nodes
	roles := p.groupByRole(deployTask.NodeConfigs)

	// create the init sub tasks with priority = 10
	initTask, err := p.createInitSubTask(deployTask, roles)
	if err != nil {
		err = fmt.Errorf("failed to create common init sub tasks: %s", err)
		logger.Error(err)
		return err
	}
	subTasks = append(subTasks, initTask)

	// create the deploy etcd sub tasks with priority = 20
	if _, ok := roles[constant.MachineRoleEtcd]; ok {
		etcdTask, err := p.createDeploySubTask(constant.MachineRoleEtcd, deployTask, roles)
		if err != nil {
			err = fmt.Errorf("failed to create deploy etcd sub tasks: %s", err)
			logger.Error(err)
			return err
		}
		subTasks = append(subTasks, etcdTask)
	}

	// create the deploy master sub tasks with priority = 30
	if _, ok := roles[constant.MachineRoleMaster]; ok {
		masterTask, err := p.createDeploySubTask(constant.MachineRoleMaster, deployTask, roles)
		if err != nil {
			err = fmt.Errorf("failed to create deploy master sub tasks: %s", err)
			logger.Error(err)
			return err
		}
		subTasks = append(subTasks, masterTask)
	}

	// create the deploy worker sub tasks with priority = 40
	if _, ok := roles[constant.MachineRoleWorker]; ok {
		workerTask, err := p.createDeploySubTask(constant.MachineRoleWorker, deployTask, roles)
		if err != nil {
			err = fmt.Errorf("failed to create deploy worker sub tasks: %s", err)
			logger.Error(err)
			return err
		}
		subTasks = append(subTasks, workerTask)
	}

	// create the deploy ingress sub tasks with priority = 50
	if _, ok := roles[constant.MachineRoleIngress]; ok {
		ingressTask, err := p.createDeploySubTask(constant.MachineRoleIngress, deployTask, roles)
		if err != nil {
			err = fmt.Errorf("failed to create deploy ingress sub tasks: %s", err)
			logger.Error(err)
			return err
		}
		subTasks = append(subTasks, ingressTask)
	}

	// create the deploy config sub task with priority = 60
	configTask, err := p.createConfigSubTask(deployTask, roles)
	if err != nil {
		err = fmt.Errorf("failed to create deploy config sub tasks: %s", err)
		logger.Error(err)
		return err
	}
	subTasks = append(subTasks, configTask)

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

func (p *deployProcessor) groupByRole(cfgs []*pb.NodeDeployConfig) map[constant.MachineRole][]*pb.NodeDeployConfig {
	roles := make(map[constant.MachineRole][]*pb.NodeDeployConfig)
	for _, nodeCfg := range cfgs {
		nodeRoles := nodeCfg.GetRoles()
		for _, role := range nodeRoles {
			roleName := constant.MachineRole(role)
			roles[roleName] = append(roles[roleName], nodeCfg)
		}
	}
	return roles
}

func (p *deployProcessor) createInitSubTask(parent *DeployTask, rn map[constant.MachineRole][]*pb.NodeDeployConfig) (task Task, err error) {

	config := &NodeInitTaskConfig{
		NodeConfigs:     parent.NodeConfigs,
		LogFileBasePath: parent.GetLogFileDir(),
		Priority:        int(initPriority),
		Parent:          parent.GetName(),
		ClusterConfig:   parent.ClusterConfig,
	}
	task, err = NewNodeInitTask(fmt.Sprintf("init"), config)
	return
}

func (p *deployProcessor) createConfigSubTask(parent *DeployTask, rn map[constant.MachineRole][]*pb.NodeDeployConfig) (task Task, err error) {

	config := &DeployConfigTaskConfig{
		NodeConfigs:     parent.NodeConfigs,
		LogFileBasePath: parent.GetLogFileDir(),
		Priority:        int(ConfigPriority),
		Parent:          parent.GetName(),
		ClusterConfig:   parent.ClusterConfig,
		MasterNodes:     p.unwrapNodes(rn[constant.MachineRoleMaster]),
	}
	task, err = NewDeployConfigTask(fmt.Sprintf("deploy-config"), config)
	return
}

func (p *deployProcessor) createDeploySubTask(role constant.MachineRole, parent *DeployTask, rn map[constant.MachineRole][]*pb.NodeDeployConfig) (task Task, err error) {

	switch role {
	case constant.MachineRoleEtcd:
		config := &DeployEtcdTaskConfig{
			Nodes:           p.unwrapNodes(rn[role]),
			LogFileBasePath: parent.GetLogFileDir(),
			Priority:        int(Priorities[role]),
			Parent:          parent.GetName(),
		}
		// Use the role name as the task name for now.
		task, err = NewDeployEtcdTask(fmt.Sprintf("deploy-%s", role), config)

	case constant.MachineRoleMaster:
		certificateKey, err := copycerts.CreateCertificateKey()
		if err != nil {
			return nil, err
		}

		config := &DeployMasterTaskConfig{
			CertKey:         certificateKey,
			NodeConfigs:     parent.NodeConfigs,
			EtcdNodes:       p.unwrapNodes(rn[constant.MachineRoleEtcd]),
			Nodes:           p.unwrapNodes(rn[role]),
			ClusterConfig:   parent.ClusterConfig,
			LogFileBasePath: parent.GetLogFileDir(),
			Priority:        int(Priorities[role]),
			Parent:          parent.GetName(),
		}
		// Use the role name as the task name for now.
		task, err = NewDeployMasterTask(fmt.Sprintf("deploy-%s", role), config)

	case constant.MachineRoleWorker:

		// Use the role name as the task name for now.
		return NewDeployWorkerTask(fmt.Sprintf("deploy-%s", role),
			&DeployWorkerTaskConfig{
				BaseTaskConfig: BaseTaskConfig{
					LogFileBasePath: parent.GetLogFileDir(), // /app/deploy/logs/unknown
					Priority:        int(Priorities[role]),
					Parent:          parent.GetName(),
				},
				Nodes:         rn[constant.MachineRoleWorker],
				ClusterConfig: parent.ClusterConfig,
				MasterNodes:   p.unwrapNodes(rn[constant.MachineRoleMaster]),
			},
		)

	case constant.MachineRoleIngress:

		// Use the role name as the task name for now.
		return NewDeployIngressTask(fmt.Sprintf("deploy-%s", role),
			&DeployIngressTaskConfig{
				BaseTaskConfig: BaseTaskConfig{
					LogFileBasePath: parent.GetLogFileDir(), // /app/deploy/logs/unknown
					Priority:        int(Priorities[role]),
					Parent:          parent.GetName(),
				},
				Nodes:         rn[constant.MachineRoleIngress],
				ClusterConfig: parent.ClusterConfig,
				MasterNodes:   p.unwrapNodes(rn[constant.MachineRoleMaster]),
			},
		)

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
