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

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type nodeCheckExecutor struct {
	config *pb.NodeCheckConfig
}

func (a *nodeCheckExecutor) Execute(act Action) error {
	logrus.Debugf("Start to execute action: %+v", act)

	nodeCheckAction, ok := act.(*nodeCheckAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be node check action, but is %T", act)
	}
	// TODO: implemented the node check logic

	// TODO: update action status
	nodeCheckAction.err = &pb.Error{
		Reason:     "todo",
		Detail:     "todo",
		FixMethods: "todo",
	}

	logrus.Debugf("End to execute action: %+v", act)
	return nil
}
