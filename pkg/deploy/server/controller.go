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

package server

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/deploy/task"
	"github.com/kpaas-io/kpaas/pkg/utils/idcreator"
)

type controller struct {
	store      task.Store
	logFileLoc string
}

func (c *controller) TestConnection(ctx context.Context, req *pb.TestConnectionRequest) (*pb.TestConnectionReply, error) {
	logrus.Info("Begins TestConnection request")

	if req == nil {
		return nil, fmt.Errorf("invalid request: request paramter is nil")
	}
	if req.Node == nil {
		return nil, fmt.Errorf("invalid request: 'Node' in request paramter is nil")
	}

	taskName := getTestConnectionTaskName(req.Node.Name)
	taskConfig := &task.TestConnectionTaskConfig{
		Node:            req.Node,
		LogFileBasePath: c.logFileLoc,
	}

	testConnTask, err := task.NewTestConnectionTask(taskName, taskConfig)
	if err != nil {
		logrus.Errorf("request failed: %s", err)
		return nil, err
	}

	if err = c.storeAndExecuteTask(testConnTask); err != nil {
		logrus.Errorf("request failed: %s", err)
		return nil, err
	}

	var reply *pb.TestConnectionReply
	taskErr := testConnTask.GetErr()
	if taskErr != nil {
		reply = &pb.TestConnectionReply{
			Passed: false,
			Err:    taskErr,
		}
	} else {
		reply = &pb.TestConnectionReply{
			Passed: true,
			Err:    nil,
		}
	}

	logrus.Infof("Ends TestConnection request, test result: %v", reply.Passed)
	return reply, nil
}

func (c *controller) CheckNodes(ctx context.Context, req *pb.CheckNodesRequest) (*pb.CheckNodesReply, error) {
	logrus.Info("Begins CheckNodes request")

	taskName := getCheckNodeTaskName()
	taskConfig := &task.NodeCheckTaskConfig{
		NodeConfigs:     req.GetConfigs(),
		NetworkOptions:  req.GetNetworkOptions(),
		LogFileBasePath: c.logFileLoc,
	}

	nodeCheckTask, err := task.NewNodeCheckTask(taskName, taskConfig)
	if err == nil {
		// store and launch the task
		err = c.storeAndLanuchTask(nodeCheckTask)
	}
	if err != nil {
		logrus.Errorf("CheckNodes request failed: %s", err)
		return &pb.CheckNodesReply{
			Accepted: false,
			Err: &pb.Error{
				Reason: consts.MsgRequestFailed,
				Detail: err.Error(),
			},
		}, err
	}

	logrus.Info("CheckNodes request succeeded")
	return &pb.CheckNodesReply{
		Accepted: true,
		Err:      nil,
	}, nil
}

func (c *controller) GetCheckNodesResult(ctx context.Context, req *pb.GetCheckNodesResultRequest) (*pb.GetCheckNodesResultReply, error) {
	logrus.Info("Begins GetCheckNodesResult request")

	var err error
	defer func() {
		if err != nil {
			logrus.Errorf("Failed to reply GetCheckNodesResult request, error: %v", err)
		} else {
			logrus.Info("Succeeded to reply GetCheckNodesResult request.")
		}
	}()

	tsk, err := c.getTask(getCheckNodeTaskName())
	if err != nil {
		return nil, err
	}

	return c.getCheckNodeResult(tsk, req.GetWithLogs())
}

func (c *controller) Deploy(ctx context.Context, req *pb.DeployRequest) (*pb.DeployReply, error) {
	logrus.Info("Begins Deploy request")

	taskName := getDeployTaskName()
	taskConfig := &task.DeployTaskConfig{
		NodeConfigs:     req.NodeConfigs,
		ClusterConfig:   req.ClusterConfig,
		LogFileBasePath: c.logFileLoc, // /app/deploy/logs
	}

	deployTask, err := task.NewDeployTask(taskName, taskConfig)
	if err == nil {
		// store and launch the task
		err = c.storeAndLanuchTask(deployTask)
	}
	if err != nil {
		logrus.Errorf("Deploy request failed: %s", err)
		return &pb.DeployReply{
			Accepted: false,
			Err: &pb.Error{
				Reason: consts.MsgRequestFailed,
				Detail: err.Error(),
			},
		}, err
	}

	logrus.Info("Deploy request succeeded")
	return &pb.DeployReply{
		Accepted: true,
		Err:      nil,
	}, nil
}

func (c *controller) GetDeployResult(ctx context.Context, req *pb.GetDeployResultRequest) (*pb.GetDeployResultReply, error) {
	logrus.Info("Begins GetDeployResult request")

	var err error
	defer func() {
		if err != nil {
			logrus.Errorf("Failed to reply GetDeployResult request, error: %v", err)
		} else {
			logrus.Info("Succeeded to reply GetDeployResult request.")
		}
	}()

	tsk, err := c.getTask(getDeployTaskName())
	if err != nil {
		return nil, err
	}

	return c.getDeployResult(tsk, req.GetWithLogs())
}

func (c *controller) FetchKubeConfig(ctx context.Context, req *pb.FetchKubeConfigRequest) (*pb.FetchKubeConfigReply, error) {
	logrus.Info("Begins FetchKubeConfig request")

	var err error
	defer func() {
		if err != nil {
			logrus.Errorf("request failed: %s", err)
		}
	}()

	taskName := getFetchKubeConfigTaskName(req)
	taskConfig := &task.FetchKubeConfigTaskConfig{
		Node:            req.Node,
		LogFileBasePath: c.logFileLoc,
	}

	kubeConfigTask, err := task.NewFetchKubeConfigTask(taskName, taskConfig)
	if err != nil {
		return nil, err
	}

	if err = c.storeAndExecuteTask(kubeConfigTask); err != nil {
		return nil, err
	}

	taskErr := kubeConfigTask.GetErr()
	if taskErr != nil {
		err = fmt.Errorf(taskErr.String())
		return &pb.FetchKubeConfigReply{
			Err: taskErr,
		}, err
	}

	logrus.Info("Ends FetchKubeConfig request: succeeded")
	return &pb.FetchKubeConfigReply{
		KubeConfig: kubeConfigTask.(*task.FetchKubeConfigTask).KubeConfig,
	}, nil
}

func (c *controller) CheckNetworkRequirements(
	context context.Context, req *pb.CheckNetworkRequirementRequest) (
	*pb.CheckNetworkRequirementsReply, error) {
	logrus.Info("Begins CheckNetworkRequirements request")
	taskConfig := &task.CheckNetworkRequirementsTaskConfig{
		Nodes:           req.GetNodes(),
		NetworkOptions:  req.GetOptions(),
		LogFileBasePath: c.logFileLoc,
	}

	taskName := "check-network-requirements"
	if taskConfig.NetworkOptions != nil {
		taskName = taskName + "-" + taskConfig.NetworkOptions.GetNetworkType()
	}

	checkTask, err := task.NewCheckNetworkRequirementsTask(taskName, taskConfig)
	if err == nil {
		err = c.storeAndExecuteTask(checkTask)
	}
	if err != nil {
		logrus.Errorf("failed to create task for CheckNetworkRequirements, error %v", err)
		return &pb.CheckNetworkRequirementsReply{
			Passed: false,
			Err: &pb.Error{
				Reason: consts.MsgRequestFailed,
				Detail: err.Error(),
			},
		}, err
	}
	logrus.Info("CheckNetworkRequirements request succeeded")
	return &pb.CheckNetworkRequirementsReply{
		Passed: true,
	}, nil
}

func (c *controller) storeTask(task task.Task) error {
	if c.store == nil {
		return fmt.Errorf("no task store")
	}

	// Currently, we have to use the same task name for some kinds of repeated requests,
	// for example, multiple check nodes and deploy requests. That means we only keep the latest
	// result for the same kind request.
	// TODO: review this design when support multiple clusters.
	return c.store.UpdateOrAddTask(task)
}

func (c *controller) getTask(name string) (task.Task, error) {
	if c.store == nil {
		return nil, fmt.Errorf("no task store")
	}

	tsk := c.store.GetTask(name)
	if tsk == nil {
		return nil, fmt.Errorf("could't find task: %s", name)
	}
	return tsk, nil
}

// Store the task and start the task, will not wait task to finish execution.
func (c *controller) storeAndLanuchTask(aTask task.Task) error {
	// store the task
	if err := c.storeTask(aTask); err != nil {
		return err
	}

	// launch the task
	return task.StartTask(aTask)
}

// Store the task and wait the task to finish execution.
func (c *controller) storeAndExecuteTask(aTask task.Task) error {
	// store the task
	if err := c.storeTask(aTask); err != nil {
		return err
	}

	// execute the task
	return task.ExecuteTask(aTask)
}

func getCheckNodeTaskName() string {
	// use a fixed name for checknode task, it may be changed in the future
	return "node-check"
}

func getDeployTaskName() string {
	// use "<cluster name>-deploy" as the deploy task name
	clusterName := "unknown"
	// TODO: review this later
	// if req.ClusterConfig != nil && reg.ClusterConfig.ClusterName != "" {
	// 	clusterName = reg.ClusterConfig.ClusterName
	// }

	return fmt.Sprintf("%s-%s", clusterName, "deploy")
}

func getFetchKubeConfigTaskName(req *pb.FetchKubeConfigRequest) string {
	// use a fixed name for now, it may be changed in the future
	return "fetch-kube-config"
}

func getTestConnectionTaskName(nodeName string) string {
	// User may test a node's connection repeatly, so create a unique task name
	// for each request
	return fmt.Sprintf("testconnection-%v-%v", nodeName, idcreator.NextString())
}

func taskStatusToOperationStatus(status task.Status) constant.OperationStatus {
	switch status {
	case task.TaskPending:
		return constant.OperationStatusPending
	case task.TaskDoing, task.TaskSplitting, task.TaskSplitted:
		return constant.OperationStatusRunning
	case task.TaskDone:
		return constant.OperationStatusSuccessful
	case task.TaskFailed:
		return constant.OperationStatusFailed
	default:
		return constant.OperationStatusUnknown
	}
}

func actionStatusToOperationStatus(status action.Status) constant.OperationStatus {
	switch status {
	case action.ActionPending:
		return constant.OperationStatusPending
	case action.ActionDoing:
		return constant.OperationStatusRunning
	case action.ActionDone:
		return constant.OperationStatusSuccessful
	case action.ActionFailed:
		return constant.OperationStatusFailed
	default:
		return constant.OperationStatusUnknown
	}
}

func itemStatusToOperationStatus(status action.ItemStatus) constant.OperationStatus {
	switch status {
	case action.ItemPending:
		return constant.OperationStatusPending
	case action.ItemDoing:
		return constant.OperationStatusRunning
	case action.ItemDone:
		return constant.OperationStatusSuccessful
	case action.ItemFailed:
		return constant.OperationStatusFailed
	default:
		return constant.OperationStatusUnknown
	}
}
