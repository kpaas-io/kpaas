// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package machine

import (
	"fmt"

	dockerclient "github.com/docker/docker/client"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine/docker"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type Machine struct {
	*ExecClient
	*pb.Node
	DockerTunnel *docker.Tunnel
	DockerClient *dockerclient.Client
}

func NewMachine(node *pb.Node) (*Machine, error) {
	client, err := NewExecClient(node)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution client for machine: %v(%v), error: %v", node.Name, node.Ip, err)
	}

	return &Machine{
		ExecClient: client,
		Node:       node,
	}, nil
}

// SetDockerTunnel create a Docker Tunnel to remote Node when needed
func (m *Machine) SetDockerTunnel() {
	m.DockerTunnel = docker.NewTunnel(m.SSHClient, m.Name)
}

func (m *Machine) SetDockerClient() error {
	client, err := docker.NewTunneledClient(m.Name)
	if err != nil {
		return fmt.Errorf("failed to create tunneled docker client to machine: %v, error: %v", m.Name, err)
	}

	m.DockerClient = client

	return nil
}

func (m *Machine) Close() {

	if m.ExecClient != nil {
		m.ExecClient.Close()
	}

	if m.DockerClient != nil {
		m.DockerClient.Close()
	}

	if m.DockerTunnel != nil {
		m.DockerTunnel.Close()
	}
}
