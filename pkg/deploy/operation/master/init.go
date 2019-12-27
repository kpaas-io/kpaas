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
	CertKey       string
	Node          *pb.Node
	MasterNodes   []*pb.Node
	EtcdNodes     []*pb.Node
	ClusterConfig *pb.ClusterConfig
}

type initMasterOperation struct {
	operation.BaseOperation
	CertKey       string
	Logger        *logrus.Entry
	EtcdNodes     []*pb.Node
	MasterNodes   []*pb.Node
	machine       machine.IMachine
	ClusterConfig *pb.ClusterConfig
}

func NewInitMasterOperation(config *InitMasterOperationConfig) (*initMasterOperation, error) {
	ops := &initMasterOperation{
		Logger:        config.Logger,
		CertKey:       config.CertKey,
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
	etcdCACrt, etcCAKey, err := etcd.FetchEtcdCertAndKey(op.EtcdNodes[0], "ca")
	if err != nil {
		return err
	}

	// put peer cert and key to all cluster nodes
	config := etcd.GetAPIServerClientCrtConfig()
	encodedAPIServerKey, encodedAPIServerCert, err := etcd.CreateFromCA(config, etcdCACrt, etcCAKey)
	if err != nil {
		return fmt.Errorf("failed to generation etcd apiserver client key and cert for apiserver node:%v, error: %v", op.machine.GetName(), err)
	}

	_, encodedEtcdCACrt, err := etcd.ToByte(etcdCACrt, nil)

	if err := op.machine.PutFile(bytes.NewReader(encodedEtcdCACrt), etcd.DefaultEtcdCACertPath); err != nil {
		return fmt.Errorf("failed to put etcd ca cert to %v:%v, error: %v", op.machine.GetName(), etcd.DefaultEtcdCACertPath, err)
	}
	if err := op.machine.PutFile(bytes.NewReader(encodedAPIServerCert), defaultApiServerEtcdClientCertPath); err != nil {
		return fmt.Errorf("failed to put apiserver etcd client cert to %v:%v, error: %v", op.machine.GetName(), defaultApiServerEtcdClientCertPath, err)
	}
	if err := op.machine.PutFile(bytes.NewReader(encodedAPIServerKey), defaultApiServerEtcdClientKeyPath); err != nil {
		return fmt.Errorf("failed to put apiserver etcd client key to %v:%v, error: %v", op.machine.GetName(), defaultApiServerEtcdClientKeyPath, err)
	}

	kubeadmConfig, err := newInitConfig(op, op.CertKey)
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

	if err := masterUpAndRunning(op); err == nil {
		op.Logger.Infof("master1 already up and running, skipping init")
		return nil
	} else {
		op.Logger.Debugf("master not running, error:%v", err)
	}

	if err := op.PreDo(); err != nil {
		return err
	}

	op.Logger.Debug("predo of init master done")

	// init first master
	stdOut, stdErr, err := op.BaseOperation.Do()
	op.Logger.Debugf("init master result:%s\n%s\n%v", stdOut, stdErr, err)

	if err != nil {
		return fmt.Errorf("failed to initilize first master, error:%s", stdErr)
	}

	op.Logger.Debug("init master done, start post do")

	if err := op.PostDo(); err != nil {
		return err
	}

	op.Logger.Debug("post do done")

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
	op.Logger.Debugf("controlPlaneEndpoint: %v", controlPlaneEndpoint)

	if err != nil {
		return fmt.Errorf("failed to get control plane endpoint, error:%v", err)
	}

	healthCheckUrl := fmt.Sprintf("https://%v/healthz", controlPlaneEndpoint)

	op.Logger.Debugf("after healthCheckUrl:%v", healthCheckUrl)

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

	op.Logger.Debugf("after get(%v):%s, %v", healthCheckUrl, body, err)

	if err != nil {
		return err
	}

	if string(body) != "ok" {
		return fmt.Errorf("controlplane status: %v, not ok", string(body))
	}

	return nil
}
