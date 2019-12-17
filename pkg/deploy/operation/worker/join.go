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

package worker

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	RemoteScriptPathSetup = "/usr/local/bin/setup.sh"
)

type JoinClusterConfig struct {
	Machine     *deployMachine.Machine
	Logger      *logrus.Entry
	Node        *pb.NodeDeployConfig
	Cluster     *pb.ClusterConfig
	MasterNodes []*pb.Node
}

type JoinCluster struct {
	operation.BaseOperation
	config *JoinClusterConfig
}

func NewJoinCluster(config *JoinClusterConfig) *JoinCluster {
	return &JoinCluster{
		config: config,
	}
}

func (operation *JoinCluster) JoinKubernetes() *pb.Error {

	controlPlaneEndpoint, err := deploy.GetControlPlaneEndpoint(operation.config.Cluster, operation.config.MasterNodes)
	if err != nil {
		return &pb.Error{
			Reason:     "Get control panel endpoint error",
			Detail:     "When deploying worker, get the control plane endpoint error",
			FixMethods: "Please create issues for us.",
		}
	}

	return RunCommand(
		command.NewShellCommand(
			operation.config.Machine,
			fmt.Sprintf("/bin/bash %s", RemoteScriptPathSetup),
			"join",
			"--token="+consts.KubernetesToken,
			"--master="+controlPlaneEndpoint,
		),
		"Join node to cluster failed",     // 添加节点到集群失败
		"join node to kubernetes cluster", // 添加节点到Kubernetes集群
	)
}

func (operation *JoinCluster) Execute() *pb.Error {

	return operation.JoinKubernetes()
}
