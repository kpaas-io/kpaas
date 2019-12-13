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

package master

import (
	"bytes"
	"fmt"
	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/etcd"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	kubeadmConfigFileName              = "kubeadm_config.yaml"
	kubeadmConfigPath                  = consts.DefaultK8sConfigDir + kubeadmConfigFileName
	defaultApiServerEtcdClientCertName = "apiserver-etcd-client.crt"
	defaultApiServerEtcdClientKeyName  = "apiserver-etcd-client.key"
	defaultApiServerEtcdClientCertPath = etcd.DefaultPKIDir + defaultApiServerEtcdClientCertName
	defaultApiServerEtcdClientKeyPath  = etcd.DefaultPKIDir + defaultApiServerEtcdClientKeyName
)

type InitMasterOperationConfig struct {
	Logger        *logrus.Entry
	Node          *pb.Node
	MasterNodes   []*pb.Node
	EtcdNodes     []*pb.Node
	ClusterConfig *pb.ClusterConfig
}

type initMasterOperation struct {
	operation.BaseOperation
	Logger        *logrus.Entry
	EtcdNodes     []*pb.Node
	MasterNodes   []*pb.Node
	machine       *machine.Machine
	ClusterConfig *pb.ClusterConfig
}

func NewInitMasterOperation(config *InitMasterOperationConfig) (*initMasterOperation, error) {
	ops := &initMasterOperation{
		Logger:        config.Logger,
		EtcdNodes:     config.EtcdNodes,
		MasterNodes:   config.MasterNodes,
		ClusterConfig: config.ClusterConfig,
	}

	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, err
	}

	ops.machine = m

	return ops, nil
}

func (op *initMasterOperation) PreDo() error {
	// TODO
	// put etcd ca cert, apiserver etcd client cert and key to first master node
	etcdCACrt := etcd.EtcdCAcrt
	if len(etcdCACrt) == 0 {
		return fmt.Errorf("failed go obtain etcd ca cert, it's empty")
	}
	apiServerEtcdClientCert := etcd.ApiServerClientCrt
	if len(apiServerEtcdClientCert) == 0 {
		return fmt.Errorf("failed go obtain apiserver client etcd cert, it's empty")
	}
	apiServerEtcdClientKey := etcd.ApiServerClientKey
	if len(apiServerEtcdClientKey) == 0 {
		return fmt.Errorf("failed go obtain apiserver client cert key, it's empty")
	}

	if err := op.machine.PutFile(bytes.NewReader(etcdCACrt), etcd.DefaultEtcdCACertPath); err != nil {
		return fmt.Errorf("failed to put etcd ca cert to %v:%v, error: %c", op.machine.Name, etcd.DefaultEtcdCACertPath)
	}
	if err := op.machine.PutFile(bytes.NewReader(apiServerEtcdClientCert), defaultApiServerEtcdClientCertPath); err != nil {
		return fmt.Errorf("failed to put apiserver etcd client cert to %v:%v, error: %c", op.machine.Name, defaultApiServerEtcdClientCertPath)
	}
	if err := op.machine.PutFile(bytes.NewReader(apiServerEtcdClientKey), defaultApiServerEtcdClientKeyPath); err != nil {
		return fmt.Errorf("failed to put apiserver etcd client key to %v:%v, error: %c", op.machine.Name, defaultApiServerEtcdClientKeyPath)
	}

	kubeadmConfig, err := newInitConfig(op)
	if err != nil {
		return fmt.Errorf("failed to generate %v, error: %v", kubeadmConfigPath, err)
	}

	if err := op.machine.PutFile(strings.NewReader(kubeadmConfig), defaultApiServerEtcdClientKeyPath); err != nil {
		return fmt.Errorf("failed to put kubeadm init config file to %v:%v, error: %v", op.machine.Name, defaultApiServerEtcdClientKeyPath, err)
	}

	// prepare commands to this operation
	endpoint, err := deploy.GetControlPlaneEndpoint(op.ClusterConfig, op.MasterNodes)
	if err != nil {
		return fmt.Errorf("failed to get control plane endpoint, error:%v", err)
	}

	initOptions := map[string]string{
		"--config":                 kubeadmConfigPath,
		"--upload-certs":           "",
		"--control-plane-endpoint": endpoint,
	}

	op.AddCommands(
		command.NewShellCommand(op.machine, "systemctl", "start kubelet", nil),
		command.NewShellCommand(op.machine, "kubeadm", "init", initOptions),
	)
	return nil
}

func (op *initMasterOperation) Do() error {
	defer op.machine.Close()

	// init first master
	stdErr, _, err := op.BaseOperation.Do()
	if err != nil {
		return fmt.Errorf("failed to initilize first master, error:%s", stdErr)
	}

	return nil
}

func (op *initMasterOperation) PostDo() error {
	// TODO
	// wait until master cluster ready
	return nil
}
