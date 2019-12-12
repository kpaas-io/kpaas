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
	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type JoinMasterOperationConfig struct {
	Logger        *logrus.Entry
	Node          *pb.Node
	MasterNodes   []*pb.Node
	ClusterConfig *pb.ClusterConfig
}

type joinMasterOperation struct {
	Logger        *logrus.Entry
	MasterNodes   []*pb.Node
	machine       *machine.Machine
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

func (d *joinMasterOperation) PreDo() error {
	// TODO
	// put apiserver etcd client cert and key to first master node

	return nil
}

func (d *joinMasterOperation) Do() error {
	// TODO
	// init first master
	return nil
}

func (d *joinMasterOperation) PostDo() error {
	// TODO
	// wait until master cluster ready
	return nil
}
