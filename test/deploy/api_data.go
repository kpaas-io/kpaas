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
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/deploy/task"
)

// Note: this file contains the env and test data to test the deploy gRPC API.
// Don't set _skip to true when commit for PR, otherwise it will break the UT.

var (
	_skip = false // set it to false if you want to run/debug the test locally.

	_launchLocalServer   = true           // launch a local deploy server and connect to it to test
	_remoteServerAddress = "0.0.0.0:8081" // connect to the remote deploy server
)

type ApiParams struct {
	request interface{}
	reply   interface{}
}

var nodes = []*pb.Node{
	&pb.Node{
		Name: "g1-node0",
		Ip:   "47.102.123.5",
		Ssh: &pb.SSH{
			Port: 22,
			Auth: &pb.Auth{
				Type:       "password",
				Username:   "root",
				Credential: "PH5S~*Fg8d7o", // replace it with the correct password when run/debug testing.
			},
		},
	},
	//&pb.Node{
	//	Name: "g1-node1",
	//	Ip:   "47.102.126.233",
	//	Ssh: &pb.SSH{
	//		Port: 22,
	//		Auth: &pb.Auth{
	//			Type:       "password",
	//			Username:   "root",
	//			Credential: "PH5S~*Fg8d7o", // replace it with the correct password when run/debug testing.
	//		},
	//	},
	//},
	//&pb.Node{
	//	Name: "g1-node2",
	//	Ip:   "106.14.187.43",
	//	Ssh: &pb.SSH{
	//		Port: 22,
	//		Auth: &pb.Auth{
	//			Type:       "password",
	//			Username:   "root",
	//			Credential: "PH5S~*Fg8d7o", // replace it with the correct password when run/debug testing.
	//		},
	//	},
	//},
	//&pb.Node{
	//	Name: "g1-node3",
	//	Ip:   "47.103.15.20",
	//	Ssh: &pb.SSH{
	//		Port: 22,
	//		Auth: &pb.Auth{
	//			Type:       "password",
	//			Username:   "root",
	//			Credential: "PH5S~*Fg8d7o", // replace it with the correct password when run/debug testing.
	//		},
	//	},
	//},
}

var testConnectionData = &ApiParams{
	request: &pb.TestConnectionRequest{
		Node: nodes[0],
	},
	reply: &pb.TestConnectionReply{
		Passed: true,
	},
}

var checkNodesData = &ApiParams{
	request: &pb.CheckNodesRequest{
		Configs: []*pb.NodeCheckConfig{
			&pb.NodeCheckConfig{
				Node:  nodes[0],
				Roles: []string{"etcd", "master"},
			},
			//&pb.NodeCheckConfig{
			//	Node:  nodes[1],
			//	Roles: []string{"etcd", "master"},
			//},
			//&pb.NodeCheckConfig{
			//	Node:  nodes[2],
			//	Roles: []string{"etcd", "master"},
			//},
			//&pb.NodeCheckConfig{
			//	Node:  nodes[3],
			//	Roles: []string{"worker"},
			//},
		},
		// TODO NetWork Check
		// NetworkOptions: &pb.NetworkOptions{
		// 	NetworkType: string(consts.NetworkTypeCalico),
		// },
	},
	reply: &pb.CheckNodesReply{
		Accepted: true,
	},
}

var checkItems = []*pb.CheckItem{
	&pb.CheckItem{
		Name:        "Docker check",
		Description: "Docker check",
	},
	&pb.CheckItem{
		Name:        "CPU check",
		Description: "CPU check",
	},
	&pb.CheckItem{
		Name:        "Kernel check",
		Description: "Kernel check",
	},
	&pb.CheckItem{
		Name:        "Memory check",
		Description: "Memory check",
	},
	&pb.CheckItem{
		Name:        "Disk check",
		Description: "Disk check",
	},
	&pb.CheckItem{
		Name:        "Distribution check",
		Description: "Distribution check",
	},
	&pb.CheckItem{
		Name:        "SystemPreference check",
		Description: "SystemPreference check",
	},
	&pb.CheckItem{
		Name:        "SystemComponent check",
		Description: "SystemComponent check",
	},
}

var itemsResult []*pb.ItemCheckResult

func init() {
	// Create check itemsResult
	for _, checkItem := range checkItems {
		result := &pb.ItemCheckResult{
			Item:   checkItem,
			Status: "done",
			Err:    nil,
		}
		itemsResult = append(itemsResult, result)
	}
}

var getCheckNodesResultData = &ApiParams{
	request: &pb.GetCheckNodesResultRequest{},
	reply: &pb.GetCheckNodesResultReply{
		Status: string(task.TaskDone),
		Err:    nil,
		Nodes: map[string]*pb.NodeCheckResult{
			nodes[0].Name: &pb.NodeCheckResult{
				NodeName: nodes[0].Name,
				Status:   "done",
				Err:      nil,
				Items:    itemsResult,
			},
		},
	},
}
