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
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	defaultEtcdDialTimeout         = 5 * time.Second
	defaultEtcdClusterReadyTimeout = 5 * time.Minute

	defaultEtcdServerPort = 2379
	defaultEtcdPeerPort   = 2380
	defaultEtcdDataDir    = "/var/lib/etcd"
	etcdImage             = "reg.kpaas.io/kpaas/etcd:3.3.15-0"
	defaultPKIDir         = "/etc/kubernetes/pki/"
	defautEtcdPKIDir      = defaultPKIDir + "etcd"

	defaultEtcdCACertName     = "ca.crt"
	defaultEtcdCAKeyName      = "ca.key"
	defaultEtcdServerCertName = "server.crt"
	defaultEtcdServerKeyName  = "server.key"
	defaultEtcdPeerCertName   = "peer.crt"
	defaultEtcdPeerKeyName    = "peer.key"

	defaultEtcdCACertPath     = defautEtcdPKIDir + "/" + defaultEtcdCACertName
	defaultEtcdCAKeyPath      = defautEtcdPKIDir + "/" + defaultEtcdCAKeyName
	defaultEtcdServerCertPath = defautEtcdPKIDir + "/" + defaultEtcdServerCertName
	defaultEtcdServerKeyPath  = defautEtcdPKIDir + "/" + defaultEtcdServerKeyName
	defaultEtcdPeerCertPath   = defautEtcdPKIDir + "/" + defaultEtcdPeerCertName
	defaultEtcdPeerKeyPath    = defautEtcdPKIDir + "/" + defaultEtcdPeerKeyName
)

var (
	// reserve for later apiserver usage
	EtcdCAcrt, ApiServerClientCrt, ApiServerClientKey []byte
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
	machine                         *machine.Machine
	clusterNodes                    []*pb.Node
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

	m.SetDockerTunnel()
	if err := m.SetDockerClient(); err != nil {
		return nil, fmt.Errorf("failed to set dockerclient, error: %v", err)
	}

	ops.machine = m
	//ops.AddCommands(command.NewShellCommand(m, "bash", "/tmp/scripts/checkdocker.sh", nil))
	return ops, nil
}

// PreDo generate etcd certs and put it to etcd node
func (d *deployEtcdOperation) PreDo() error {
	if err := d.machine.DockerTunnel.Start(); err != nil {
		return fmt.Errorf("failed to start docker tunnel to remote node: %v, error: %v", d.machine.Name, err)
	}

	// put ca cert and key to all cluster nodes
	encodedCert, encodedKey, err := ToByte(d.caCrt, d.caKey)
	if err != nil {
		return fmt.Errorf("failed to convert key and cert to byte, error: %v", err)
	}

	// save for later use
	EtcdCAcrt = encodedCert

	if err := d.machine.PutFile(bytes.NewReader(encodedCert), defaultEtcdCACertPath); err != nil {
		return fmt.Errorf("failed to put ca cert to:%v, error: %v", d.machine.Name, err)
	}
	if err := d.machine.PutFile(bytes.NewReader(encodedKey), defaultEtcdCAKeyPath); err != nil {
		return fmt.Errorf("failed to put ca key to:%v, error: %v", d.machine.Name, err)
	}

	// put server cert and key to all cluster nodes
	config, err := GetServerCrtConfig(d.machine.Name, d.machine.Ip)
	if err != nil {
		return fmt.Errorf("failed to get etd server cert config for node:%v, error: %v", d.machine.Name, err)
	}
	encodedCert, encodedKey, err = CreateFromCA(config, d.caCrt, d.caKey)
	if err != nil {
		return fmt.Errorf("failed to generation etcd server key and cert for etcd node:%v, error: %v", d.machine.Name, err)
	}
	// set for later etcd client use
	d.encodedPeerCert = encodedCert
	d.encodedPeerKey = encodedKey

	if err := d.machine.PutFile(bytes.NewReader(encodedCert), defaultEtcdServerCertPath); err != nil {
		return fmt.Errorf("failed to put etcd server cert to:%v, error: %v", d.machine.Name, err)
	}
	if err := d.machine.PutFile(bytes.NewReader(encodedKey), defaultEtcdServerKeyPath); err != nil {
		return fmt.Errorf("failed to put etcd server key to:%v, error: %v", d.machine.Name, err)
	}

	// put peer cert and key to all cluster nodes
	config, err = GetPeerCrtConfig(d.machine.Name, d.machine.Ip)
	if err != nil {
		return fmt.Errorf("failed to get etd peer cert config for node:%v, error: %v", d.machine.Name, err)
	}
	encodedCert, encodedKey, err = CreateFromCA(config, d.caCrt, d.caKey)
	if err != nil {
		return fmt.Errorf("failed to generation etcd peer key and cert for etcd node:%v, error: %v", d.machine.Name, err)
	}
	if err := d.machine.PutFile(bytes.NewReader(encodedCert), defaultEtcdPeerCertPath); err != nil {
		return fmt.Errorf("failed to put etcd peer cert to:%v, error: %v", d.machine.Name, err)
	}
	if err := d.machine.PutFile(bytes.NewReader(encodedKey), defaultEtcdPeerKeyPath); err != nil {
		return fmt.Errorf("failed to put etcd peer key to:%v, error: %v", d.machine.Name, err)
	}

	// put peer cert and key to all cluster nodes
	config = GetAPIServerClientCrtConfig()
	encodedCert, encodedKey, err = CreateFromCA(config, d.caCrt, d.caKey)
	if err != nil {
		return fmt.Errorf("failed to generation etcd apiserver client key and cert for etcd node:%v, error: %v", d.machine.Name, err)
	}

	// save for later apiserver use
	ApiServerClientCrt, ApiServerClientKey = encodedCert, encodedKey

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

func composeEtcdDockerCmd(d *deployEtcdOperation) []string {

	cmd := []string{"etcd"}

	cmd = append(cmd, "--client-cert-auth=true")
	cmd = append(cmd, "--peer-client-cert-auth=true")

	cmd = append(cmd, fmt.Sprintf("--snapshot-count=%v)", 10000))

	cmd = append(cmd, fmt.Sprintf("--name=%v)", d.machine.Name))
	cmd = append(cmd, fmt.Sprintf("--data-dir=%v", defaultEtcdDataDir))
	cmd = append(cmd, fmt.Sprintf("--key-file=%v", defaultEtcdServerKeyPath))
	cmd = append(cmd, fmt.Sprintf("--cert-file=%v", defaultEtcdServerCertPath))
	cmd = append(cmd, fmt.Sprintf("--peer-cert-file=%v", defaultEtcdPeerCertPath))
	cmd = append(cmd, fmt.Sprintf("--peer-key-file=%v)", defaultEtcdPeerKeyPath))
	cmd = append(cmd, fmt.Sprintf("--trusted-ca-file=%v)", defaultEtcdCACertPath))
	cmd = append(cmd, fmt.Sprintf("--peer-trusted-ca-file=%v)", defaultEtcdCACertPath))

	cmd = append(cmd, fmt.Sprintf("--advertise-client-urls=https://%v:%v", d.machine.Ip, defaultEtcdServerPort))
	cmd = append(cmd, fmt.Sprintf("--initial-advertise-peer-urls=https://%v:%v", d.machine.Ip, defaultEtcdPeerPort))
	cmd = append(cmd, fmt.Sprintf("--listen-client-urls=https://127.0.0.1:2379,https://%v:%v", d.machine.Ip, defaultEtcdServerPort))
	cmd = append(cmd, fmt.Sprintf("--listen-peer-urls=https://%v:%v", d.machine.Ip, defaultEtcdPeerPort))

	//initial-cluster: infra0=https://10.0.0.6:2380,infra1=https://10.0.0.7:2380,infra2=https://10.0.0.8:2380
	cmd = append(cmd, composeInitialClusterUrl(d.clusterNodes))

	return cmd
}

func (d *deployEtcdOperation) Do() error {
	defer d.machine.DockerTunnel.Close()

	if err := d.PreDo(); err != nil {
		return err
	}

	config := &container.Config{
		Cmd:   composeEtcdDockerCmd(d),
		Image: etcdImage,
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/etc/kubernetes/pki/etcd",
				Target: "/etc/kubernetes/pki/etcd",
			},
		},
		NetworkMode: "host",
	}

	//create and start etcd containers
	body, err := d.machine.DockerClient.ContainerCreate(context.Background(), config, hostConfig, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create etcd container on etcd node:%v, error: %v", d.machine.Name, err)
	}

	if err := d.machine.DockerClient.ContainerStart(context.Background(), body.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start etcd container on etcd node:%v, error: %v", d.machine.Name, err)
	}

	// post do
	if err := d.PostDo(); err != nil {
		return err
	}

	return nil
}

func (d *deployEtcdOperation) PostDo() error {
	cli, err := newEtcdV3SecureClient(d)
	if err != nil {
		return fmt.Errorf("failed to get etcd client of etcd node:%v, error:%v", d.machine.Name, err)
	}
	defer cli.Close()

	deadline := time.Now().Add(defaultEtcdClusterReadyTimeout)
	for retries := 0; time.Now().Before(deadline); retries++ {
		resp, err := cli.MemberList(context.Background())
		if len(resp.Members) == len(d.clusterNodes) {
			return nil
		}

		d.logger.Warnf("etd cluster not ready, error: %v, will retry", err)
		time.Sleep(time.Second << uint(retries))
	}

	return fmt.Errorf("wait for etcd cluster ready timeout after:%v", defaultEtcdClusterReadyTimeout)
}
