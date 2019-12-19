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

	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	RegisterExecutor(ActionTypeTestConnection, new(testConnectionExecutor))
}

type testConnectionExecutor struct {
}

func (a *testConnectionExecutor) Execute(act Action) *pb.Error {
	testConnTask, ok := act.(*TestConnectionAction)
	if !ok {
		return errOfTypeMismatched(new(TestConnectionAction), act)
	}
	if act.GetNode() == nil {
		return consts.PBErrActionNodeEmpty
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
		consts.LogFieldNode:   act.GetNode().GetName(),
	})

	logger.Debug("Start to execute action")

	// machine.NewMachine() will test if the machine can be connected via ssh.
	m, err := machine.NewMachine(testConnTask.Node)
	if err != nil {
		pbErr := &pb.Error{
			Reason: "failed to test connection",
			Detail: err.Error(),
		}
		deploy.PBErrLogger(pbErr, logger).Debug()
		return pbErr
	}
	m.Close()

	logger.Debug("Finsih to execute action")
	return nil
}
