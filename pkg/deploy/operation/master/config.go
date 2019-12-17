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
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/etcd"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
	"sigs.k8s.io/yaml"

	"github.com/kpaas-io/kpaas/pkg/deploy"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	Token = "4a996f.8f1da0db96f8e50e"
)

func newInitConfig(op *initMasterOperation) (string, error) {
	var (
		err           error
		initYaml      bytes.Buffer
		initConfig    v1beta2.InitConfiguration
		clusterConfig v1beta2.ClusterConfiguration
	)

	initYaml.Write([]byte("---\n"))

	initConfig.TypeMeta = metav1.TypeMeta{
		Kind:       "InitConfiguration",
		APIVersion: "kubeadm.k8s.io/v1beta2",
	}

	initConfig.BootstrapTokens = make([]v1beta2.BootstrapToken, 1)
	initConfig.BootstrapTokens[0].Token = new(v1beta2.BootstrapTokenString)

	initConfig.BootstrapTokens[0].Token.ID = "4a996f"
	initConfig.BootstrapTokens[0].Token.Secret = "8f1da0db96f8e50e"
	initConfig.BootstrapTokens[0].TTL = &metav1.Duration{
		Duration: time.Duration(0),
	}

	clusterConfig.TypeMeta = metav1.TypeMeta{
		Kind:       "ClusterConfiguration",
		APIVersion: "kubeadm.k8s.io/v1beta2",
	}

	clusterConfig.KubernetesVersion = op.ClusterConfig.KubernetesVersion
	clusterConfig.ImageRepository = op.ClusterConfig.ImageRepository

	clusterConfig.ControlPlaneEndpoint, err = deploy.GetControlPlaneEndpoint(op.ClusterConfig, op.MasterNodes)
	if err != nil {
		return "", fmt.Errorf("failed to get control plane endpoint addr, error: %v", err)
	}

	clusterConfig.Networking = v1beta2.Networking{
		ServiceSubnet: op.ClusterConfig.ServiceSubnet,
		PodSubnet:     op.ClusterConfig.PodSubnet,
	}

	clusterConfig.Etcd.External = getExternalEtcd(op.EtcdNodes)

	initConfigData, err := yaml.Marshal(initConfig)
	if err != nil {
		return "", err
	}

	initYaml.Write(initConfigData)

	initYaml.Write([]byte("\n---\n"))
	clusterConfigData, err := yaml.Marshal(clusterConfig)
	if err != nil {
		return "", err
	}
	initYaml.Write(clusterConfigData)

	return initYaml.String(), nil
}

func getExternalEtcd(etcdNodes []*pb.Node) (externalEtcd *v1beta2.ExternalEtcd) {
	for i := range etcdNodes {
		// TODO: replace to use etcd const when pr merged
		ep := fmt.Sprintf("https:%v:%v", etcdNodes[i].Ip, 2379)
		externalEtcd.Endpoints = append(externalEtcd.Endpoints, ep)
	}

	externalEtcd.CAFile = etcd.DefaultEtcdCACertPath
	externalEtcd.CertFile = defaultApiServerEtcdClientCertPath
	externalEtcd.KeyFile = defaultApiServerEtcdClientKeyPath

	return
}
