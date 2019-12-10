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
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/etcd"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type deployEtcdExecutor struct {
}

func (a *deployEtcdExecutor) Execute(act Action) error {
	etcdAction, ok := act.(*deployEtcdAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be deploy etcd action, but is %T", act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	logger.Debug("Start to execute deploy etcd action")

	config := &etcd.DeployEtcdOperationConfig{
		Logger:       logger,
		Node:         etcdAction.node,
		CACrt:        etcdAction.caCrt,
		CAKey:        etcdAction.caKey,
		ClusterNodes: etcdAction.clusterNodes,
	}
	op, err := etcd.NewDeployEtcdOperation(config)
	if err != nil {
		return fmt.Errorf("failed to get etcd operation, error: %v", err)
	}

	logger.Debugf("Start to deploy etcd on nodes: %s", etcdAction.node.Name)

	etcdAction.status = ActionDone
	if err := op.Do(); err != nil {
		etcdAction.status = ActionFailed
		etcdAction.err = &pb.Error{
			Reason:     err.Error(),
			Detail:     err.Error(),
			FixMethods: "",
		}
	}

	logger.Debug("Finish to execute deploy etcd action")
	return nil
}
