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
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/master"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	RegisterExecutor(ActionTypeInitMaster, new(initMasterExecutor))
}

type initMasterExecutor struct {
}

func (a *initMasterExecutor) Execute(act Action) *pb.Error {
	action, ok := act.(*InitMasterAction)
	if !ok {
		return errOfTypeMismatched(new(InitMasterAction), act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	logger.Debug("Start to init first master action")

	var needUntaint bool
	rolesSet := sets.NewString(action.Roles...)
	if rolesSet.Has(string(constant.MachineRoleIngress)) || rolesSet.Has(string(constant.MachineRoleWorker)) {
		needUntaint = true
	}
	config := &master.InitMasterOperationConfig{
		Logger:        logger,
		CertKey:       action.CertKey,
		Node:          action.Node,
		NeedUntaint:   needUntaint,
		MasterNodes:   action.MasterNodes,
		EtcdNodes:     action.EtcdNodes,
		ClusterConfig: action.ClusterConfig,
	}

	op, err := master.NewInitMasterOperation(config)
	if err != nil {
		return &pb.Error{
			Reason: "failed to get init master operation",
			Detail: err.Error(),
		}
	}

	logger.Debugf("Start to init master on nodes: %s", action.Node.Name)

	if err := op.Do(); err != nil {
		return &pb.Error{
			Reason:     "failed to do init master operation",
			Detail:     err.Error(),
			FixMethods: "",
		}
	}

	logger.Debug("Finish to execute init master action")
	return nil
}
