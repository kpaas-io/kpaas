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

package deploy

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// This file along with the following local config file contains the env and test data to test the deploy gRPC API.

// localConfigFilePath stores the local config for deploy system testing.
// The following is an example, make sure there is at least 4 nodes.
/*
{
	"skip": false,
	"launchLocalServer": true,
	"remoteServerAddress": "0.0,0.0:8081",
	"nodes": [
		{
			"name": "g1-node0",
			"ip":   "47.102.123.5",
			"ssh": {
				"port": 22,
				"auth": {
					"type":       "password",
					"username":   "root",
					"credential": "the real password or private key"
				}
			}
		},
		{
			"name": "g1-node1",
			"ip":   "47.102.126.233",
			"ssh": {
				"port": 22,
				"auth": {
					"type":       "password",
					"username":   "root",
					"credential": "the real password or private key"
				}
			}
		},
		{
			"name": "g1-node2",
			"ip":   "106.14.187.43",
			"ssh": {
				"port": 22,
				"auth": {
					"type":       "password",
					"username":   "root",
					"credential": "the real password or private key"
				}
			}
		},
		{
			"name": "g1-node3",
			"ip":   "47.103.15.20",
			"ssh": {
				"port": 22,
				"auth": {
					"type":       "password",
					"username":   "root",
					"credential": "the real password or private key"
				}
			}
		}
	]
}
*/
const localConfigFilePath = "./tmp/config/deploy.json"

// Don't set Skip to false here, otherwise it will break the testing. It should be
// reset through the local config file
var _testConfig = DeployTestConfig{
	Skip:              true,
	LaunchLocalServer: true,
}

type DeployTestConfig struct {
	// Skip defines whether to skip the deploy system testing when run "go test"
	Skip bool `json:"skip,omitempty"`
	// LaunchLocalServer defines whether to launch a local deploy server for deploy system testing
	LaunchLocalServer bool `json:"launchLocalServer,omitempty"`
	// RemoteServerAddress defines an already running deploy server for deploy system testing, it will
	// only take effect if LaunchLocalServer is set to false.
	RemoteServerAddress string `json:"remoteServerAddress,omitempty"`
	// Nodes defines the nodes (machines) information for testing.
	Nodes []*pb.Node `json:"nodes,omitempty"`
}

func init() {
	if err := loadLocalConfig(); err != nil {
		_testConfig.Skip = true
	}
}

func loadLocalConfig() error {
	file, err := os.Open(localConfigFilePath)
	if err != nil {
		logrus.Infof("Failed to open local config file: %v", localConfigFilePath)
		return err
	}
	defer file.Close()

	if configBytes, err := ioutil.ReadAll(file); err != nil {
		logrus.Infof("Failed to read local config file: %v", localConfigFilePath)
		return err
	} else if err = json.Unmarshal(configBytes, &_testConfig); err != nil {
		logrus.Infof("Failed to unmarshal local config content: %v", string(localConfigFilePath))
		return err
	}
	return nil
}

func getTestConnectionData() (request *pb.TestConnectionRequest, reply *pb.TestConnectionReply) {
	request = &pb.TestConnectionRequest{
		Node: _testConfig.Nodes[0],
	}
	reply = &pb.TestConnectionReply{
		Passed: true,
	}
	return
}

func getCheckNodesData() (request *pb.CheckNodesRequest, reply *pb.CheckNodesReply) {
	request = &pb.CheckNodesRequest{
		Configs: []*pb.NodeCheckConfig{
			&pb.NodeCheckConfig{
				Node:  _testConfig.Nodes[0],
				Roles: []string{"etcd", "master"},
			},
			&pb.NodeCheckConfig{
				Node:  _testConfig.Nodes[1],
				Roles: []string{"etcd", "master"},
			},
			&pb.NodeCheckConfig{
				Node:  _testConfig.Nodes[2],
				Roles: []string{"etcd", "master"},
			},
			&pb.NodeCheckConfig{
				Node:  _testConfig.Nodes[3],
				Roles: []string{"worker"},
			},
		},
		// TODO NetWork Check
		NetworkOptions: &pb.NetworkOptions{
			NetworkType: string(consts.NetworkTypeCalico),
			CalicoOptions: &pb.CalicoOptions{
				CheckConnectivityAll: true,
				EncapsulationMode:    "vxlan",
				VxlanPort:            4789,
			},
		},
	}
	reply = &pb.CheckNodesReply{
		Accepted: true,
	}
	return
}

func getGetCheckNodesResultData() (request *pb.GetCheckNodesResultRequest, reply *pb.GetCheckNodesResultReply) {
	request = &pb.GetCheckNodesResultRequest{}
	reply = &pb.GetCheckNodesResultReply{
		Status: string(constant.OperationStatusSuccessful),
		Err:    nil,
		Nodes: map[string]*pb.NodeCheckResult{
			_testConfig.Nodes[0].Name: &pb.NodeCheckResult{
				NodeName: _testConfig.Nodes[0].Name,
				Status:   string(constant.OperationStatusSuccessful),
				Err:      nil,
				Items:    nil,
			},
			_testConfig.Nodes[1].Name: &pb.NodeCheckResult{
				NodeName: _testConfig.Nodes[1].Name,
				Status:   string(constant.OperationStatusSuccessful),
				Err:      nil,
				Items:    nil,
			},
			_testConfig.Nodes[2].Name: &pb.NodeCheckResult{
				NodeName: _testConfig.Nodes[2].Name,
				Status:   string(constant.OperationStatusSuccessful),
				Err:      nil,
				Items:    nil,
			},
			_testConfig.Nodes[3].Name: &pb.NodeCheckResult{
				NodeName: _testConfig.Nodes[3].Name,
				Status:   string(constant.OperationStatusSuccessful),
				Err:      nil,
				Items:    nil,
			},
		},
	}

	var checkItems = []*pb.CheckItem{
		&pb.CheckItem{
			Name:        "check docker",
			Description: "检查 docker 环境",
		},
		&pb.CheckItem{
			Name:        "check cpu",
			Description: "检查 cpu 环境",
		},
		&pb.CheckItem{
			Name:        "check kernel",
			Description: "检查 kernel 环境",
		},
		&pb.CheckItem{
			Name:        "check memory",
			Description: "检查 memory 环境",
		},
		&pb.CheckItem{
			Name:        "check disk",
			Description: "检查 disk 环境",
		},
		&pb.CheckItem{
			Name:        "check distribution",
			Description: "检查 distribution 环境",
		},
		&pb.CheckItem{
			Name:        "check system-preference",
			Description: "检查 system-preference 环境",
		},
		&pb.CheckItem{
			Name:        "check system-manager",
			Description: "检查 system-manager 环境",
		},
		&pb.CheckItem{
			Name:        "check port-occupied",
			Description: "检查 port-occupied 环境",
		},
		&pb.CheckItem{
			Name:        "connectivity-check-BGP",
			Description: "check connectivity to BGP port",
		},
		&pb.CheckItem{
			Name:        "connectivity-check-kube-API",
			Description: "check connectivity to kubernetes API port",
		},
		&pb.CheckItem{
			Name:        "connectivity-check-vxlan",
			Description: "check connectivity for vxlan packets",
		},
	}
	var itemsResult []*pb.ItemCheckResult
	// Create check itemsResult
	for _, checkItem := range checkItems {
		result := &pb.ItemCheckResult{
			Item:   checkItem,
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		}
		itemsResult = append(itemsResult, result)
	}
	for _, checkResult := range reply.Nodes {
		checkResult.Items = itemsResult
	}
	return
}

func getGetCheckNodesLogData() (request *pb.GetCheckNodesLogRequest, reply *pb.GetCheckNodesLogReply) {
	request = &pb.GetCheckNodesLogRequest{
		NodeName: _testConfig.Nodes[0].Name,
	}
	reply = &pb.GetCheckNodesLogReply{}
	return
}

var clusterConfig = &pb.ClusterConfig{
	ClusterName: "TestCluster",
	KubeAPIServerConnect: &pb.KubeAPIServerConnect{
		Type: "firstMasterIP",
	},
	NodePortRange: &pb.NodePortRange{
		From: 30000,
		To:   32000,
	},
	NodeLabels:      map[string]string{"nodelabelkey": "nodelabelvalue"},
	NodeAnnotations: map[string]string{"nodeannokey": "nodeannovalue"},
	// TODO: the following paramters are using default value for now.
	// ImageRepository: "docker.io/kpaas",
	// PodSubnet: "",
	// ServiceSubnet: "",
	// KubernetesVersion: "",
}

func getDeployMultipleNodesData() (request *pb.DeployRequest, reply *pb.DeployReply) {
	request = &pb.DeployRequest{
		NodeConfigs: []*pb.NodeDeployConfig{
			&pb.NodeDeployConfig{
				Node:  _testConfig.Nodes[0],
				Roles: []string{"etcd", "master"},
				Labels: map[string]string{
					"kpaas-io/role": "master",
					"testkey":       "testvalue",
				},
				Taints: []*pb.Taint{
					&pb.Taint{
						Key:    "kpaas-io/role",
						Value:  "master",
						Effect: "NoSchedule",
					},
				},
			},
			&pb.NodeDeployConfig{
				Node:  _testConfig.Nodes[1],
				Roles: []string{"etcd", "master", "ingress"},
				Labels: map[string]string{
					"kpaas-io/role": "master",
					"testkey":       "testvalue",
				},
				Taints: []*pb.Taint{
					&pb.Taint{
						Key:    "kpaas-io/role",
						Value:  "master",
						Effect: "NoSchedule",
					},
				},
			},
			&pb.NodeDeployConfig{
				Node:  _testConfig.Nodes[2],
				Roles: []string{"etcd", "master", "ingress"},
				Labels: map[string]string{
					"kpaas-io/role": "master",
					"testkey":       "testvalue",
				},
				Taints: []*pb.Taint{
					&pb.Taint{
						Key:    "kpaas-io/role",
						Value:  "master",
						Effect: "NoSchedule",
					},
				},
			},
			&pb.NodeDeployConfig{
				Node:  _testConfig.Nodes[3],
				Roles: []string{"worker"},
				Labels: map[string]string{
					"kpaas-io/role": "worker",
					"testkey":       "testvalue",
				},
			},
		},
		ClusterConfig: clusterConfig,
	}
	reply = &pb.DeployReply{
		Accepted: true,
	}
	return
}

func getDeployResultData() (request *pb.GetDeployResultRequest, reply *pb.GetDeployResultReply) {
	var deployItemResults = []*pb.DeployItemResult{
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "etcd",
				NodeName: _testConfig.Nodes[0].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "etcd",
				NodeName: _testConfig.Nodes[1].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "etcd",
				NodeName: _testConfig.Nodes[2].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "master",
				NodeName: _testConfig.Nodes[0].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "master",
				NodeName: _testConfig.Nodes[1].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "master",
				NodeName: _testConfig.Nodes[2].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "ingress",
				NodeName: _testConfig.Nodes[1].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "ingress",
				NodeName: _testConfig.Nodes[2].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "worker",
				NodeName: _testConfig.Nodes[3].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
	}

	request = &pb.GetDeployResultRequest{}
	reply = &pb.GetDeployResultReply{
		Status: string(constant.OperationStatusSuccessful),
		Err:    nil,
		Items:  deployItemResults,
	}
	return
}

func getDeployAllInOneData() (request *pb.DeployRequest, reply *pb.DeployReply) {
	request = &pb.DeployRequest{
		NodeConfigs: []*pb.NodeDeployConfig{
			&pb.NodeDeployConfig{
				Node:  _testConfig.Nodes[0],
				Roles: []string{"etcd", "master", "worker", "ingress"},
				Labels: map[string]string{
					"kpaas-io/test": "abcdef",
					"testkey":       "testvalue",
				},
				Taints: []*pb.Taint{
					&pb.Taint{
						Key:    "kpaas-io/test",
						Value:  "testtaints",
						Effect: "NoSchedule",
					},
				},
			},
		},
		ClusterConfig: clusterConfig,
	}
	reply = &pb.DeployReply{
		Accepted: true,
	}
	return
}

func getDeployAllInOneResultData() (request *pb.GetDeployResultRequest, reply *pb.GetDeployResultReply) {
	var deployItemResults = []*pb.DeployItemResult{
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "etcd",
				NodeName: _testConfig.Nodes[0].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "master",
				NodeName: _testConfig.Nodes[0].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "worker",
				NodeName: _testConfig.Nodes[0].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
		&pb.DeployItemResult{
			DeployItem: &pb.DeployItem{
				Role:     "ingress",
				NodeName: _testConfig.Nodes[0].Name,
			},
			Status: string(constant.OperationStatusSuccessful),
			Err:    nil,
		},
	}

	request = &pb.GetDeployResultRequest{}
	reply = &pb.GetDeployResultReply{
		Status: string(constant.OperationStatusSuccessful),
		Err:    nil,
		Items:  deployItemResults,
	}
	return
}

func getGetDeployLogData() (request *pb.GetDeployLogRequest, reply *pb.GetDeployLogReply) {
	request = &pb.GetDeployLogRequest{
		Role:     "master",
		NodeName: _testConfig.Nodes[0].Name,
	}
	reply = &pb.GetDeployLogReply{}
	return
}

func getFetchKubeConfigData() (request *pb.FetchKubeConfigRequest, reply *pb.FetchKubeConfigReply) {
	request = &pb.FetchKubeConfigRequest{
		Node: _testConfig.Nodes[0],
	}
	reply = &pb.FetchKubeConfigReply{}
	return
}
