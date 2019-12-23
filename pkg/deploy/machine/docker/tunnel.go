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

package docker

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/kpaas-io/kpaas/pkg/deploy"
)

const (
	localSocketDir     = "/tmp/"
	dockerSocketSuffix = ".docker.sock"
	remoteDockerSocket = "/var/run/docker.sock"
)

type Tunnel struct {
	done                chan struct{}
	remoteHostName      string
	localUnixSocketFile string
	sshClient           *ssh.Client
}

func NewTunnel(sshClient *ssh.Client, hostName string) *Tunnel {
	return &Tunnel{
		sshClient:           sshClient,
		done:                make(chan struct{}),
		localUnixSocketFile: composeLocalDockerSocketFile(hostName),
	}
}

func (t *Tunnel) Start() (err error) {
	listener, err := net.Listen("unix", t.localUnixSocketFile)
	if err != nil {
		return fmt.Errorf("failed to listen local docker, error: %v", err)
	}

	for {
		select {
		case <-t.done:
			return
		default:

		}

		localConn, ea := listener.Accept()
		if err != nil {
			logrus.Errorf("failed to accept local connection, error: %v", ea)
			continue
		}

		dstConn, ed := t.sshClient.Dial("unix", remoteDockerSocket)
		if err != nil {
			logrus.Errorf("failed to dial %v:%v, error: %v", t.remoteHostName, remoteDockerSocket, ed)
			continue
		}

		go t.forward(dstConn, localConn)
	}
}

func (t *Tunnel) forward(dst, src net.Conn) {
	defer func() {
		dst.Close()
		src.Close()
	}()

	connCopy := func(dst io.Writer, src io.Reader) error {
		if _, err := io.Copy(dst, src); err != nil {
			logrus.Errorf("connection copy failed: %v", err)
			return err
		}

		return nil
	}

	go connCopy(src, dst)
	go connCopy(dst, src)

	<-t.done
}

func (t *Tunnel) Close() (err error) {
	close(t.done)

	if deploy.FileExist(t.localUnixSocketFile) {
		if err = os.Remove(t.localUnixSocketFile); err != nil {
			err = fmt.Errorf("failed to remove socket file: %v, error: %v", t.localUnixSocketFile, err)
		}
	}

	return
}

func composeLocalDockerSocketFile(name string) string {
	return strings.Join([]string{localSocketDir, name, dockerSocketSuffix}, "")
}
