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

package master

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type JoinMasterOperationConfig struct {
	Logger        *logrus.Entry
	Node          *pb.Node
	MasterNodes   []*pb.Node
	ClusterConfig *pb.ClusterConfig
}

type joinMasterOperation struct {
	operation.BaseOperation
	Logger        *logrus.Entry
	MasterNodes   []*pb.Node
	machine       machine.IMachine
	ClusterConfig *pb.ClusterConfig
}

func NewJoinMasterOperation(config *JoinMasterOperationConfig) (*joinMasterOperation, error) {
	ops := &joinMasterOperation{
		Logger:        config.Logger,
		MasterNodes:   config.MasterNodes,
		ClusterConfig: config.ClusterConfig,
	}

	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}

	ops.machine = m

	return ops, nil
}

func (op *joinMasterOperation) PreDo() error {
	// compose join command
	//kubeadm join 192.168.0.200:6443 --token 9vr73a.a8uxyaju799qwdjv --control-plane --discovery-token-unsafe-skip-ca-verification
	endpoint, err := deploy.GetControlPlaneEndpoint(op.ClusterConfig, op.MasterNodes)
	if err != nil {
		return fmt.Errorf("failed to get control plane endpoint addr, error: %v", err)
	}

	op.AddCommands(
		command.NewShellCommand(op.machine, "systemctl", "start", "kubelet"),
		command.NewShellCommand(op.machine, "kubeadm", "join", endpoint,
			"--token", Token,
			"--control-plane",
			"--discovery-token-unsafe-skip-ca-verification"),
	)
	return nil
}

func (op *joinMasterOperation) Do() error {
	defer op.machine.Close()

	if err := op.PreDo(); err != nil {
		return err
	}

	// join master
	stdErr, _, err := op.BaseOperation.Do()
	if err != nil {
		return fmt.Errorf("failed to join master:%v to cluster, error:%s", op.machine.GetName(), stdErr)
	}
	return nil
}

//func (op *joinMasterOperation) PostDo() error {
//	// TODO
//	// wait until master cluster ready
//	return nil
//}
