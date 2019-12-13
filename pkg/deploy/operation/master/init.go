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
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/etcd"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	defaultControlPlaneReadyTimeout    = 5 * time.Minute
	kubeadmConfigFileName              = "kubeadm_config.yaml"
	kubeadmConfigPath                  = consts.DefaultK8sConfigDir + "/" + kubeadmConfigFileName
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
	machine       machine.IMachine
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
		return fmt.Errorf("failed to put etcd ca cert to %v:%v, error: %v", op.machine.GetName(), etcd.DefaultEtcdCACertPath, err)
	}
	if err := op.machine.PutFile(bytes.NewReader(apiServerEtcdClientCert), defaultApiServerEtcdClientCertPath); err != nil {
		return fmt.Errorf("failed to put apiserver etcd client cert to %v:%v, error: %v", op.machine.GetName(), defaultApiServerEtcdClientCertPath, err)
	}
	if err := op.machine.PutFile(bytes.NewReader(apiServerEtcdClientKey), defaultApiServerEtcdClientKeyPath); err != nil {
		return fmt.Errorf("failed to put apiserver etcd client key to %v:%v, error: %v", op.machine.GetName(), defaultApiServerEtcdClientKeyPath, err)
	}

	kubeadmConfig, err := newInitConfig(op)
	if err != nil {
		return fmt.Errorf("failed to generate %v, error: %v", kubeadmConfigPath, err)
	}

	if err := op.machine.PutFile(strings.NewReader(kubeadmConfig), kubeadmConfigPath); err != nil {
		return fmt.Errorf("failed to put kubeadm init config file to %v:%v, error: %v", op.machine.GetName(), defaultApiServerEtcdClientKeyPath, err)
	}

	op.AddCommands(
		command.NewShellCommand(op.machine, "systemctl", "start", "kubelet"),
		command.NewShellCommand(op.machine, "kubeadm", "init",
			"--config", kubeadmConfigPath,
			"--upload-certs"),
	)
	return nil
}

func (op *initMasterOperation) Do() error {
	defer op.machine.Close()

	if err := op.PreDo(); err != nil {
		return err
	}

	// init first master
	stdErr, _, err := op.BaseOperation.Do()
	if err != nil {
		return fmt.Errorf("failed to initilize first master, error:%s", stdErr)
	}

	if err := op.PostDo(); err != nil {
		return err
	}

	return nil
}

func (op *initMasterOperation) PostDo() error {
	// wait until master cluster ready

	deadline := time.Now().Add(defaultControlPlaneReadyTimeout)
	for retries := 0; time.Now().Before(deadline); retries++ {
		err := masterUpAndRunning(op)
		if err == nil {
			return nil
		}
		op.Logger.Warnf("controlplane not ready, error: %v, will retry", err)
		time.Sleep(time.Second << uint(retries))
	}

	return fmt.Errorf("wait for controlplane to be ready timeout after:%v", defaultControlPlaneReadyTimeout)
}

func masterUpAndRunning(op *initMasterOperation) error {
	controlPlaneEndpoint, err := deploy.GetControlPlaneEndpoint(op.ClusterConfig, op.MasterNodes)
	if err != nil {
		return fmt.Errorf("failed to get control plane endpoint, error:%v", err)
	}

	healthCheckUrl := fmt.Sprintf("https://%v/healthz", controlPlaneEndpoint)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpC := &http.Client{
		Timeout:   time.Minute,
		Transport: tr,
	}

	resp, err := httpC.Get(healthCheckUrl)
	if err != nil {
		return fmt.Errorf("get %v failed, error: %v", healthCheckUrl, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "ok" {
		return fmt.Errorf("controlplane status: %v, not ok", string(body))
	}

	return nil
}