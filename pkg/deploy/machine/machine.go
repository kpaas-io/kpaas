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

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type IMachine interface {
	GetName() string
	GetIp() string
	GetNode() *pb.Node
	Close()

	Run(cmd string) (stdout, stderr []byte, err error)
	FetchDir(localDir, remoteDir string, fileNeeded func(path string) bool) error
	FetchFile(dst io.Writer, remotePath string) error
	FetchFileToLocalPath(localPath, remotePath string) error
	PutDir(localDir, remoteDir string, fileNeeded func(path string) bool) error
	PutFile(content io.Reader, remotePath string) error
}

type Machine struct {
	*ExecClient
	*pb.Node
}

func NewMachine(node *pb.Node) (IMachine, error) {
	client, err := NewExecClient(node)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution client for machine: %v(%v), error: %v", node.Name, node.Ip, err)
	}

	return &Machine{
		ExecClient: client,
		Node:       node,
	}, nil
}

func (m *Machine) Close() {

	if m.ExecClient != nil {
		m.ExecClient.Close()
	}

}

func (m *Machine) GetName() string {
	return m.Name
}

func (m *Machine) GetIp() string {
	return m.Ip
}

func (m *Machine) GetNode() *pb.Node {
	return m.Node
}
