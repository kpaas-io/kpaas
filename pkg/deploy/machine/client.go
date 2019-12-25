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
	"io/ioutil"
	"os"

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

	if m.SSHClient != nil {
		m.SSHClient.Close()
	}

	if m.SFTPClient != nil {
		m.SFTPClient.Close()
	}

	return
}

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it with permissions perm;
// otherwise WriteFile truncates it before writing.
// Like ioutil.WriteFile
func (m *ExecClient) WriteFile(filename string, data []byte, perm os.FileMode) (err error) {

	var file *sftp.File

	file, err = m.SFTPClient.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return
	}
	defer func() {
		closeError := file.Close()
		if closeError != nil {
			err = closeError
		}
	}()

	var n int
	n, err = file.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
		return
	}

	err = file.Chmod(perm)
	return
}

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
// Like ioutil.ReadFile
func (m *ExecClient) ReadFile(filename string) (content []byte, err error) {

	var file *sftp.File

	file, err = m.SFTPClient.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeError := file.Close()
		if closeError != nil {
			err = closeError
		}
	}()

	return ioutil.ReadAll(file)
}
