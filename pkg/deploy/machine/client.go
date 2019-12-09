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

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	mssh "github.com/kpaas-io/kpaas/pkg/deploy/machine/ssh"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type ExecClient struct {
	SSHClient  *ssh.Client
	SFTPClient *sftp.Client
}

// NewExecClient create a new execution client
func NewExecClient(node *pb.Node) (*ExecClient, error) {
	// use IP as host to create ssh client
	sshClient, err := mssh.NewClient(node.Ssh.Auth.Username, node.Ip, node.Ssh)
	if err != nil {
		return nil, fmt.Errorf("failed to create new ssh client to machine: %v(%v), error: %v", node.Name, node.Ip, err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get sftep client to machine: %v(%v), error: %v", node.Name, node.Ip, err)
	}

	return &ExecClient{
		SSHClient:  sshClient,
		SFTPClient: sftpClient,
	}, nil
}

// Close will only close ssh client
// no need to close sftp client since it's based on ssh client
func (m *ExecClient) Close() {
	m.SSHClient.Close()
}
