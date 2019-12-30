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
	"io"
	"strings"

	dockerclient "github.com/docker/docker/client"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

var (
	errMachineErr = fmt.Errorf("this is an error machine")

	//IsTesting only set it in test package
	IsTesting = false
)

type MockMachine struct {
	*pb.Node
	DockerClient *dockerclient.Client
}

func newMockMachine(node *pb.Node) (IMachine, error) {
	if node.Name == "error" {
		return nil, errMachineErr
	}

	return &MockMachine{
		Node: node,
	}, nil
}

func (m *MockMachine) GetName() string {
	return m.Name
}

func (m *MockMachine) GetIp() string {
	return m.Ip
}

func (m *MockMachine) Close() {}

func (m *MockMachine) StartDockerTunnel() error {
	return nil
}

// Run return different response by node name
func (m *MockMachine) Run(cmd string) (stdout, stderr []byte, err error) {
	if m.Name == "error" {
		err = fmt.Errorf("this is an error machine")
	}

	switch {
	case m.Name == "error":
		return nil, nil, errMachineErr
	case strings.HasPrefix(cmd, "cat /proc/cpuinfo"):
		return []byte("8"), nil, nil
	case strings.HasPrefix(cmd, "docker version"):
		return []byte("18.09.0"), nil, nil
	case strings.HasPrefix(cmd, "uname -r"):
		return []byte("5.18.5-041805-generic"), nil, nil
	case strings.HasPrefix(cmd, "free -b"):
		return []byte("270455574528"), nil, nil
	case strings.HasPrefix(cmd, "df -B1"):
		return []byte("294605168640"), nil, nil
	case strings.HasPrefix(cmd, "ps -p 1"):
		return []byte("systemd"), nil, nil
	case strings.HasPrefix(cmd, "cat /etc/*-release"):
		return []byte("ubuntu"), nil, nil
	}

	return []byte(""), []byte(""), nil
}

func (m *MockMachine) FetchDir(localDir, remoteDir string, fileNeeded func(path string) bool) error {
	if m.Name == "error" {
		return errMachineErr
	}
	return nil
}

func (m *MockMachine) FetchFile(dst io.Writer, remotePath string) error {
	if m.Name == "error" {
		return errMachineErr
	}

	dst.Write([]byte("this is test data"))
	return nil
}

func (m *MockMachine) FetchFileToLocalPath(localPath, remotePath string) error {
	if m.Name == "error" {
		return errMachineErr
	}
	return nil
}

func (m *MockMachine) PutDir(localDir, remoteDir string, fileNeeded func(path string) bool) error {
	if m.Name == "error" {
		return errMachineErr
	}
	return nil
}

func (m *MockMachine) PutFile(content io.Reader, remotePath string) error {
	if m.Name == "error" {
		return errMachineErr
	}
	return nil
}

func (m *MockMachine) GetNode() *pb.Node {
	return m.Node
}
