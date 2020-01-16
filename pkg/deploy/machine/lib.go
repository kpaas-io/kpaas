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

package machine

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/kpaas-io/kpaas/pkg/deploy"
	mssh "github.com/kpaas-io/kpaas/pkg/deploy/machine/ssh"
)

// Run will run command on remote machine
func (m *Machine) Run(cmd string) (stdout, stderr []byte, err error) {
	session, err := mssh.NewSession(m.SSHClient)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get session of machine(%v), error: %v", m.Name, err)
	}

	defer session.Close()

	errReader, err := session.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pipe stderr for cmd(%v) on machine(%v), error: %v", cmd, m.Name, err)
	}

	outReader, err := session.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pipe stdout for cmd(%v) on machine(%v), error: %v", cmd, m.Name, err)
	}

	if err = session.Start(cmd); err != nil {
		return nil, nil, fmt.Errorf("unable to run cmd(%v) on machine(%v), error: %v", cmd, m.Name, err)
	}

	stderr, err = ioutil.ReadAll(errReader)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read stderr message for cmd(%v) returned from machine(%v), error: %v", cmd, m.Name, err)
	}

	stdout, err = ioutil.ReadAll(outReader)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read stderr message for cmd(%v) returned from machine(%v), error: %v", cmd, m.Name, err)
	}

	err = session.Wait()
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			return stdout, stderr, fmt.Errorf("command exited with error: %v", exitErr)
		} else if exitMissingErr, ok := err.(*ssh.ExitMissingError); ok {
			return stdout, stderr, fmt.Errorf("command exit status missing, error %v",
				exitMissingErr)
		} else {
			return stdout, stderr, err
		}
	}

	return
}

func (m *Machine) PutFile(content io.Reader, remotePath string) error {
	// create parent dir if not exists
	remoteDir := path.Dir(remotePath)
	if err := m.SFTPClient.MkdirAll(remoteDir); err != nil {
		return fmt.Errorf("mkdirall %v failed, error: %v", remoteDir, err)
	}

	remoteFile, err := m.SFTPClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("create file %v failed: %v", remotePath, err)
	}
	defer remoteFile.Close()

	if _, err = io.Copy(remoteFile, content); err != nil {
		return fmt.Errorf("copy content to remote file %v failed: %v", remotePath, err)
	}

	logrus.Debugf("put file to: %v", remotePath)

	return nil
}

func (m *Machine) FetchFileToLocalPath(localPath, remotePath string) error {
	logrus.Debugf("Begin to fetch file from %s on %s to %s", remotePath, m.Name, localPath)

	// create parent dir if not exists
	localDir := path.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("mkdirall %v failed, error: %v", localDir, err)
	}
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("create local file %v failed, error: %v", localPath, err)
	}
	defer localFile.Close()

	return m.FetchFile(localFile, remotePath)
}

func (m *Machine) FetchFile(dst io.Writer, remotePath string) error {
	if dst == nil {
		return fmt.Errorf("the destination is nil")
	}
	remoteFile, err := m.SFTPClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("open remote file %v failed, error: %v", remotePath, err)
	}
	defer remoteFile.Close()

	if _, err = io.Copy(dst, remoteFile); err != nil {
		return fmt.Errorf("copy from remote file %v failed, error: %v", remotePath, err)
	}

	logrus.Debugf("fetch file from %s on %s", remotePath, m.Name)

	return nil
}

func (m *Machine) FetchDir(localDir, remoteDir string, fileNeeded func(path string) bool) error {
	logrus.Debugf("fetch %v:%v to %v", m.Name, remoteDir, localDir)

	remoteDir = strings.TrimSuffix(remoteDir, "/")
	localDir = strings.TrimSuffix(localDir, "/") + "/" + filepath.Base(remoteDir)

	if _, err := m.SFTPClient.Stat(remoteDir); os.IsNotExist(err) {
		return fmt.Errorf("%v:%v does not exist", m.Name, remoteDir)
	}

	walker := m.SFTPClient.Walk(remoteDir)

	for walker.Step() {
		if err := walker.Err(); err != nil {
			return err
		}

		remotePath := walker.Path()
		info, err := m.SFTPClient.Stat(remotePath)
		if err != nil {
			return fmt.Errorf("stat %v:%v failed, error: %v", m.Name, remotePath, err)
		}

		localPath := localDir + strings.TrimPrefix(remotePath, remoteDir)

		if info.IsDir() {
			if !deploy.FileExist(localPath) {
				logrus.Debugf("make dir: %v", localPath)
				if err := os.MkdirAll(localPath, 0755); err != nil {
					return fmt.Errorf("failed to mkdir %v, error: %v", localPath, err)
				}
			}

			continue
		}

		if !fileNeeded(remotePath) {
			continue
		}

		logrus.Debugf("fetch %v:%v to %v", m.Name, remotePath, localPath)

		if err := m.FetchFileToLocalPath(localPath, remotePath); err != nil {
			return fmt.Errorf("failed to fetch file from: %v:%v to %v, error: %v", m.Name, remotePath, localPath, err)
		}
	}

	return nil
}

func (m *Machine) PutDir(localDir, remoteDir string, fileNeeded func(path string) bool) error {
	logrus.Debugf("copy %v to %v:%v", localDir, m.Name, remoteDir)

	localDir = strings.TrimPrefix(strings.TrimSuffix(localDir, "/"), "./")
	remoteDir = strings.TrimSuffix(remoteDir, "/") + "/" + filepath.Base(localDir)

	if !deploy.FileExist(localDir) {
		return fmt.Errorf("local directory:%v doesn't exist", localDir)
	}

	if err := filepath.Walk(localDir, func(localPath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk %v:%v failed, error: %v", m.Name, localPath, err)
		}

		remotePath := remoteDir + strings.TrimPrefix(localPath, localDir)
		logrus.Debugf("copy %v to %v:%v", localPath, m.Name, remotePath)

		// create directory
		if info.IsDir() {
			if _, err := m.SFTPClient.Stat(remotePath); os.IsNotExist(err) {
				if err := m.SFTPClient.MkdirAll(remotePath); err != nil {
					return fmt.Errorf("creating %v:%v failed. error: %v", m.Name, remotePath, err)
				}
			}
		} else {
			if fileNeeded(localPath) {
				// copy file
				localFile, err := os.Open(localPath)
				if err != nil {
					return fmt.Errorf("open %v failed, error: %v", localPath, err)
				}

				if err := m.PutFile(localFile, remotePath); err != nil {
					return fmt.Errorf("failed to copy file:%v to %v:%v, error: %v", localPath, m.Name, remotePath, err)
				}
			}
		}

		return nil

	}); err != nil {
		return fmt.Errorf("walk %v failed. error: %v", localDir, err)
	}

	return nil
}
