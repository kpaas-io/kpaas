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

package docker

import (
	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	script    = "/scripts/check_docker_version.sh"
	remoteDir = "/tmp"
)

type CheckDockerOperation struct {
	operation.BaseOperation
}

func NewCheckDockerOperation(config *pb.NodeCheckConfig) (operation.Operation, error) {
	ops := &CheckDockerOperation{}
	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}

	scriptFile, err := assets.Assets.Open(script)
	if err != nil {
		return nil, err
	}

	if err := m.PutFile(scriptFile, remoteDir+script); err != nil {
		return nil, err
	}

	ops.AddCommands(command.NewShellCommand(m, "bash", remoteDir+script, nil))
	return ops, nil
}

// new docker version check to compare
func CheckDockerVersion(dockerVersion string, standardVersion string, comparedSymbol string) error {
	err := operation.CheckVersion(dockerVersion, standardVersion, comparedSymbol)
	if err != nil {
		return err
	}
	return nil
}
