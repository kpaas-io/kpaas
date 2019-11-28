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

package ssh

import (
	"fmt"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	defaultTimeout = 60 * time.Second
)

func newConfig(user string, auth *protos.Auth) (*ssh.ClientConfig, error) {
	var authMethod ssh.AuthMethod

	switch {
	case auth.Type == "password":
		authMethod = ssh.Password(auth.Credential)
	case auth.Type == "privatekey":
		privateKey, err := ssh.ParsePrivateKey([]byte(auth.Credential))
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v, error: %v", auth.Credential, err)
		}

		authMethod = ssh.PublicKeys(privateKey)
	default:
		return nil, fmt.Errorf("unrecognized auth type: %v", auth.Type)
	}

	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         defaultTimeout,
	}, nil
}

func NewClient(user string, host string, sshConfig *protos.SSH) (*ssh.Client, error) {
	config, err := newConfig(user, sshConfig.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to get ssh client config: %v, error: %v", err)
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", host, sshConfig.Port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v, error: %v", host, err)
	}

	return client, nil
}

func NewSession(client *ssh.Client) (*ssh.Session, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create new ssh, error: %v", err)
	}

	return session, nil
}
