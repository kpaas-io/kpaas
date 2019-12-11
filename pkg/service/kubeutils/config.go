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

package kubeutils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
	clientutils "github.com/kpaas-io/kpaas/pkg/service/grpcutils/client"
	"github.com/kpaas-io/kpaas/pkg/service/model/sshcertificate"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
)

const (
	// DefaultKubeConfigDirectory directory for storing kubeconfig file for each cluster
	DefaultKubeConfigDirectory = ".kpaas/kubeconfigs/"
)

// KubeConfigPathForCluster returns the local path of kubeconfig file
// for accessing kubernetes API server in specified cluster
func KubeConfigPathForCluster(clusterName string) (string, error) {
	logEntry := logrus.WithField("cluster", clusterName)

	// TODO: use a configurable directory/filename, maybe add an argument for it?
	homeDir, _ := os.UserHomeDir()
	kubeConfigDirectory := homeDir + "/" + DefaultKubeConfigDirectory
	filename := kubeConfigDirectory + "cluster-" + clusterName + ".conf"
	// returns if the file is already exist
	// TODO: add expiration for local kubeconfig file?
	if _, err := os.Stat(filename); err == nil {
		logEntry.WithField("filename", filename).
			Debug("kubeconfig file already found locally")
		return filename, nil
	}
	logEntry.WithField("filename", filename).
		Debug("kubeconfig file not found locally, try to fetch it from wizard...")
	err := os.MkdirAll(DefaultKubeConfigDirectory, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create directory, error %v", err)
	}
	// download kubeconfig and save it.
	w := wizard.GetCurrentWizard()
	if w == nil {
		return "", fmt.Errorf("failed to get wizard")
	}
	if w.KubeConfig == nil || (*w.KubeConfig) == "" {
		return "", fmt.Errorf("kubeconfig file is not ready yet, try it later")
	}
	// following code are from deploy.fetchKubeConfigContent
	// To prevent potential dependency loop, implement code again here
	// TODO: put the code for fetching kubeconfig into common library
	var masterNode *wizard.Node
	client := clientutils.GetDeployController()
	for _, node := range w.Nodes {
		if node.IsMatchMachineRole(constant.MachineRoleMaster) {
			masterNode = node
			break
		}
	}

	if masterNode == nil {
		return "", fmt.Errorf("no master node ready in cluster %s", clusterName)
	}
	logEntry.WithField("nodename", masterNode.Name).WithField("IP", masterNode.IP).
		Debug("fetch kubeconfig from master  node...")
	connectionData := masterNode.ConnectionData
	sshAuth := protos.Auth{
		Username: masterNode.Username,
		Type:     string(connectionData.AuthenticationType),
	}
	switch connectionData.AuthenticationType {
	case wizard.AuthenticationTypePassword:
		sshAuth.Credential = connectionData.Password
	case wizard.AuthenticationTypePrivateKey:
		sshAuth.Credential = sshcertificate.GetPrivateKey(connectionData.PrivateKeyName)
	}
	fetchResponse, err := client.FetchKubeConfig(context.Background(),
		&protos.FetchKubeConfigRequest{Node: &protos.Node{
			Name: masterNode.Name,
			Ip:   masterNode.IP,
			Ssh: &protos.SSH{
				Port: uint32(connectionData.Port),
				Auth: &sshAuth,
			},
		}})
	if err != nil {
		return "", fmt.Errorf("failed to get response of fetching kubeconfig")
	}
	kubeConfigContent := fetchResponse.GetKubeConfig()

	ioutil.WriteFile(filename, kubeConfigContent, 0644)
	return filename, nil
}
