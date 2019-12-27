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

package etcd

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	localEtcdCADir = "/tmp"

	defaultEtcdDialTimeout         = 5 * time.Second
	defaultEtcdClusterReadyTimeout = 5 * time.Minute

	defaultEtcdServerPort = 2379
	defaultEtcdPeerPort   = 2380
	defaultEtcdDataDir    = "/var/lib/etcd"
	// TODO: registry should be obtained from cluster config
	defaultRegistry      = "docker.io"
	defaultEtcdImageRepo = "kpaas"
	defaultEtcdImageTag  = "3.3.15-0"
	defaultEtcdImageName = "etcd"
	defaultEtcdImageUrl  = defaultRegistry + "/" + defaultEtcdImageRepo + "/" + defaultEtcdImageName + ":" + defaultEtcdImageTag

	DefaultPKIDir    = "/etc/kubernetes/pki/"
	defautEtcdPKIDir = DefaultPKIDir + "etcd"

	defaultEtcdCACertName     = "ca.crt"
	defaultEtcdCAKeyName      = "ca.key"
	defaultEtcdServerCertName = "server.crt"
	defaultEtcdServerKeyName  = "server.key"
	defaultEtcdPeerCertName   = "peer.crt"
	defaultEtcdPeerKeyName    = "peer.key"

	DefaultEtcdCACertPath     = defautEtcdPKIDir + "/" + defaultEtcdCACertName
	defaultEtcdCAKeyPath      = defautEtcdPKIDir + "/" + defaultEtcdCAKeyName
	defaultEtcdServerCertPath = defautEtcdPKIDir + "/" + defaultEtcdServerCertName
	defaultEtcdServerKeyPath  = defautEtcdPKIDir + "/" + defaultEtcdServerKeyName
	defaultEtcdPeerCertPath   = defautEtcdPKIDir + "/" + defaultEtcdPeerCertName
	defaultEtcdPeerKeyPath    = defautEtcdPKIDir + "/" + defaultEtcdPeerKeyName
)

type DeployEtcdOperationConfig struct {
	Logger       *logrus.Entry
	CACrt        *x509.Certificate
	CAKey        crypto.Signer
	Node         *pb.Node
	ClusterNodes []*pb.Node
}

type deployEtcdOperation struct {
	operation.BaseOperation
	logger                          *logrus.Entry
	caCrt                           *x509.Certificate
	caKey                           crypto.Signer
	encodedPeerCert, encodedPeerKey []byte
	machine                         machine.IMachine
	clusterNodes                    []*pb.Node
	containerName                   string
}

func NewDeployEtcdOperation(config *DeployEtcdOperationConfig) (*deployEtcdOperation, error) {
	ops := &deployEtcdOperation{
		logger:       config.Logger,
		caCrt:        config.CACrt,
		caKey:        config.CAKey,
		clusterNodes: config.ClusterNodes,
	}
	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}

	ops.machine = m
	//ops.AddCommands(command.NewShellCommand(m, "bash", "/tmp/scripts/checkdocker.sh", nil))
	return ops, nil
}

func (d *deployEtcdOperation) composeContainerName() {
	d.containerName = fmt.Sprintf("etcd-kpaas-%v", d.machine.GetName())
}

func (d *deployEtcdOperation) removeExistEtcdContainer() error {
	d.logger.Debug("start removeExistEtcdContainer")

	filterArg := fmt.Sprintf("name=%v", d.containerName)

	d.logger.Debugf("filterArg: %v", filterArg)

	d.AddCommands(
		command.NewShellCommand(d.machine,
			"docker",
			"ps",
			"-q",
			"--filter",
			filterArg,
		),
	)

	stdOut, stdErr, err := d.BaseOperation.Do()
	// reset d.Commands
	d.ResetCommands()

	if err != nil {
		return fmt.Errorf("failed to get existing docker container, error:%s", stdErr)
	}

	if len(stdOut) == 0 {
		return nil
	}

	containerID := string(stdOut)

	d.logger.Debugf("remove existing ectd container: %s", stdOut)

	// found existing container, removing
	d.AddCommands(
		command.NewShellCommand(d.machine,
			"docker",
			"rm",
			"-f",
			containerID,
		),
	)

	stdOut, stdErr, err = d.BaseOperation.Do()
	// reset d.Commands
	d.ResetCommands()

	if err != nil {
		return fmt.Errorf("failed to remove existing docker container, error:%s", stdErr)
	}

	return nil
}

// PreDo generate etcd certs and put it to etcd node
func (d *deployEtcdOperation) PreDo() (err error) {
	d.composeContainerName()

	if err = d.removeExistEtcdContainer(); err != nil {
		return err
	}

	// put ca cert and key to all cluster nodes
	encodedCAKey, encodedCACert, err := ToByte(d.caCrt, d.caKey)
	if err != nil {
		return fmt.Errorf("failed to convert key and cert to byte, error: %v", err)
	}

	if err := d.machine.PutFile(bytes.NewReader(encodedCACert), DefaultEtcdCACertPath); err != nil {
		return fmt.Errorf("failed to put ca cert to:%v, error: %v", d.machine.GetName(), err)
	}
	if err := d.machine.PutFile(bytes.NewReader(encodedCAKey), defaultEtcdCAKeyPath); err != nil {
		return fmt.Errorf("failed to put ca key to:%v, error: %v", d.machine.GetName(), err)
	}

	// put server cert and key to all cluster nodes
	config, err := GetServerCrtConfig(d.machine.GetName(), d.machine.GetIp())
	if err != nil {
		return fmt.Errorf("failed to get etd server cert config for node:%v, error: %v", d.machine.GetName(), err)
	}
	encodedServerKey, encodedServerCert, err := CreateFromCA(config, d.caCrt, d.caKey)
	if err != nil {
		return fmt.Errorf("failed to generation etcd server key and cert for etcd node:%v, error: %v", d.machine.GetName(), err)
	}

	if err := d.machine.PutFile(bytes.NewReader(encodedServerCert), defaultEtcdServerCertPath); err != nil {
		return fmt.Errorf("failed to put etcd server cert to:%v, error: %v", d.machine.GetName(), err)
	}
	if err := d.machine.PutFile(bytes.NewReader(encodedServerKey), defaultEtcdServerKeyPath); err != nil {
		return fmt.Errorf("failed to put etcd server key to:%v, error: %v", d.machine.GetName(), err)
	}

	// put peer cert and key to all cluster nodes
	config, err = GetPeerCrtConfig(d.machine.GetName(), d.machine.GetIp())
	if err != nil {
		return fmt.Errorf("failed to get etd peer cert config for node:%v, error: %v", d.machine.GetName(), err)
	}
	encodedPeerKey, encodedPeerCert, err := CreateFromCA(config, d.caCrt, d.caKey)
	if err != nil {
		return fmt.Errorf("failed to generation etcd peer key and cert for etcd node:%v, error: %v", d.machine.GetName(), err)
	}

	// set for later etcd client use
	d.encodedPeerCert = encodedPeerCert
	d.encodedPeerKey = encodedPeerKey

	if err := d.machine.PutFile(bytes.NewReader(encodedPeerCert), defaultEtcdPeerCertPath); err != nil {
		return fmt.Errorf("failed to put etcd peer cert to:%v, error: %v", d.machine.GetName(), err)
	}
	if err := d.machine.PutFile(bytes.NewReader(encodedPeerKey), defaultEtcdPeerKeyPath); err != nil {
		return fmt.Errorf("failed to put etcd peer key to:%v, error: %v", d.machine.GetName(), err)
	}

	return nil
}

func composeInitialClusterUrl(nodes []*pb.Node) (clusterUrl string) {
	for i := range nodes {
		if clusterUrl == "" {
			clusterUrl = fmt.Sprintf("%v=https://%v:%v", nodes[i].Name, nodes[i].Ip, defaultEtcdPeerPort)
			continue
		}

		clusterUrl = fmt.Sprintf("%v=https://%v:%v,%v", nodes[i].Name, nodes[i].Ip, defaultEtcdPeerPort, clusterUrl)
	}

	return
}

func (d *deployEtcdOperation) composeEtcdDockerCmd() {

	cmd := []string{"etcd"}

	cmd = append(cmd, "--client-cert-auth=true")
	cmd = append(cmd, "--peer-client-cert-auth=true")

	cmd = append(cmd, fmt.Sprintf("--snapshot-count=%v", 10000))

	cmd = append(cmd, fmt.Sprintf("--name=%v", d.machine.GetName()))
	cmd = append(cmd, fmt.Sprintf("--data-dir=%v", defaultEtcdDataDir))
	cmd = append(cmd, fmt.Sprintf("--key-file=%v", defaultEtcdServerKeyPath))
	cmd = append(cmd, fmt.Sprintf("--cert-file=%v", defaultEtcdServerCertPath))
	cmd = append(cmd, fmt.Sprintf("--peer-cert-file=%v", defaultEtcdPeerCertPath))
	cmd = append(cmd, fmt.Sprintf("--peer-key-file=%v", defaultEtcdPeerKeyPath))
	cmd = append(cmd, fmt.Sprintf("--trusted-ca-file=%v", DefaultEtcdCACertPath))
	cmd = append(cmd, fmt.Sprintf("--peer-trusted-ca-file=%v", DefaultEtcdCACertPath))

	cmd = append(cmd, fmt.Sprintf("--advertise-client-urls=https://%v:%v", d.machine.GetIp(), defaultEtcdServerPort))
	cmd = append(cmd, fmt.Sprintf("--initial-advertise-peer-urls=https://%v:%v", d.machine.GetIp(), defaultEtcdPeerPort))
	cmd = append(cmd, fmt.Sprintf("--listen-client-urls=https://0.0.0.0:%v", defaultEtcdServerPort))
	cmd = append(cmd, fmt.Sprintf("--listen-peer-urls=https://0.0.0.0:%v", defaultEtcdPeerPort))

	//initial-cluster: infra0=https://10.0.0.6:2380,infra1=https://10.0.0.7:2380,infra2=https://10.0.0.8:2380
	cmd = append(cmd, fmt.Sprintf("--initial-cluster=%v", composeInitialClusterUrl(d.clusterNodes)))

	nameArg := fmt.Sprintf("--name=%v", d.containerName)

	d.AddCommands(
		command.NewShellCommand(d.machine, "docker",
			"run",
			"-d",
			"--restart=always",
			"--net=host",
			"-v",
			"/etc/kubernetes/pki/etcd:/etc/kubernetes/pki/etcd",
			"-v",
			"/var/lib/etcd:/var/lib/etcd",
			nameArg,
			defaultEtcdImageUrl,
			strings.Join(cmd, " "),
		),
	)
}

func (d *deployEtcdOperation) Do() error {
	defer d.machine.Close()
	// save
	originCACrt, originCAKey, originEncodedPeerCert, originEncodedPeerKey := d.caCrt, d.caKey, d.encodedPeerCert, d.encodedPeerKey

	etcdCACrt, etcdCAKey, caErr := FetchEtcdCertAndKey(d.machine.GetNode(), "ca")
	peerCert, peerKey, peerErr := FetchEtcdCertAndKey(d.machine.GetNode(), "peer")
	encodedPeerKey, encodedPeerCert, toByteErr := ToByte(peerCert, peerKey)

	if caErr == nil && peerErr == nil && toByteErr == nil {
		d.caCrt, d.caKey, d.encodedPeerCert, d.encodedPeerKey = etcdCACrt, etcdCAKey, encodedPeerCert, encodedPeerKey

		if err := etcdUpAndRunning(d); err == nil {
			d.logger.Info("etcd cluster already up and running, skipping deploy")
			return nil
		} else {
			d.logger.Debugf("etcd cluster not up and running, error:%v", err)
		}
	}

	// restore and clear etcd if any error occurred
	d.caCrt, d.caKey, d.encodedPeerCert, d.encodedPeerKey = originCACrt, originCAKey, originEncodedPeerCert, originEncodedPeerKey

	if err := d.PreDo(); err != nil {
		return err
	}

	d.composeEtcdDockerCmd()

	d.logger.Debugf("start command: %v", d.Commands)

	stdOut, stdErr, err := d.BaseOperation.Do()
	if err != nil {
		return fmt.Errorf("failed to exec command: %q on machine:%v", d.Commands, d.machine.GetName())
	}

	d.logger.Debugf("exec command: %#v done, %s, %s, %v", d.Commands, stdOut, stdErr, err)

	// post do
	if err := d.PostDo(); err != nil {
		d.logger.Errorf("post do error:%v", err)
		return err
	}

	d.logger.Debug("deploy etcd done")

	return nil
}

func (d *deployEtcdOperation) PostDo() error {
	deadline := time.Now().Add(defaultEtcdClusterReadyTimeout)
	for retries := 0; time.Now().Before(deadline); retries++ {
		err := etcdUpAndRunning(d)
		if err == nil {
			return nil
		}

		d.logger.Warnf("etd cluster not ready, error: %v, will retry", err)
		time.Sleep(time.Second << uint(retries))
	}

	return fmt.Errorf("wait for etcd cluster ready timeout after:%v", defaultEtcdClusterReadyTimeout)
}

func etcdUpAndRunning(d *deployEtcdOperation) error {
	d.logger.Debug("start check etcd cluster status")
	cli, err := newEtcdV3SecureClient(d)
	if err != nil {
		return fmt.Errorf("failed to get etcd client of etcd node:%v, error:%v", d.machine.GetName(), err)
	}
	defer cli.Close()

	d.logger.Debugf("etcd client:%#v", cli)

	resp, err := cli.MemberList(context.Background())

	d.logger.Debugf("member list done, result:%#v, error: %v", resp, err)

	if len(resp.Members) == len(d.clusterNodes) {
		d.logger.Infof("%v members detected", len(resp.Members))
		return nil
	}

	return err
}
