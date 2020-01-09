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

package init

import (
	"fmt"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type InitHostNameOperation struct {
	operation.BaseOperation
	Machine        machine.IMachine
	NodeInitAction *operation.NodeInitAction
}

func (itOps *InitHostNameOperation) RunCommands(node *pb.Node, initAction *operation.NodeInitAction) (stdOut, stdErr []byte, err error) {

	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, nil, err
	}

	itOps.Machine = m
	itOps.NodeInitAction = initAction

	// close ssh client if machine is not nil
	if itOps.Machine != nil {
		defer itOps.Machine.Close()
	}

	currentName := node.Name
	if currentName == "" {
		return nil, nil, fmt.Errorf("node name can not be empty")
	}

	itOps.AddCommands(command.NewShellCommand(m, "hostnamectl", fmt.Sprintf("set-hostname %v", currentName)))

	if len(itOps.Commands) == 0 {
		return nil, nil, fmt.Errorf("init host alias command is empty")
	}

	// run commands
	stdOut, stdErr, err = itOps.Do()

	return
}
