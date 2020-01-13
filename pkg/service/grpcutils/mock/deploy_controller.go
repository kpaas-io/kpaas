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

package mock

import (
	"context"

	"google.golang.org/grpc"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type DeployController struct {
}

func NewDeployController() protos.DeployContollerClient {
	return &DeployController{}
}

func (mock *DeployController) TestConnection(ctx context.Context, in *protos.TestConnectionRequest, opts ...grpc.CallOption) (*protos.TestConnectionReply, error) {

	return &protos.TestConnectionReply{
		Passed: true,
		Err:    nil,
	}, nil
}

func (mock *DeployController) CheckNodes(ctx context.Context, in *protos.CheckNodesRequest, opts ...grpc.CallOption) (*protos.CheckNodesReply, error) {

	return &protos.CheckNodesReply{
		Accepted: true,
		Err:      nil,
	}, nil
}

func (mock *DeployController) GetCheckNodesResult(ctx context.Context, in *protos.GetCheckNodesResultRequest, opts ...grpc.CallOption) (*protos.GetCheckNodesResultReply, error) {

	return &protos.GetCheckNodesResultReply{
		Status: "passed",
		Err:    nil,
		Nodes: map[string]*protos.NodeCheckResult{
			"master1": {
				NodeName: "master1",
				Status:   "passed",
				Err:      nil,
				Items: []*protos.ItemCheckResult{
					{
						Item: &protos.CheckItem{
							Name:        "Root disk size",
							Description: "Check root disk size > 50G",
						},
						Status: "passed",
						Err:    nil,
						Logs: `Filesystem   512-blocks     Used Available Capacity iused      ifree %iused  Mounted on
/dev/disk1s5  489620264 21361864  49491704    31%  483568 2447617752    0%   /`,
					},
					{
						Item: &protos.CheckItem{
							Name:        "kernel support",
							Description: "Check kernel version",
						},
						Status: "passed",
						Err:    nil,
						Logs:   `Darwin LuckyMac.local 19.0.0 Darwin Kernel Version 19.0.0: Thu Oct 17 16:17:15 PDT 2019; root:xnu-6153.41.3~29/RELEASE_X86_64 x86_64`,
					},
				},
			},
		},
	}, nil
}
func (mock *DeployController) Deploy(ctx context.Context, in *protos.DeployRequest, opts ...grpc.CallOption) (*protos.DeployReply, error) {

	return &protos.DeployReply{
		Accepted: true,
		Err:      nil,
	}, nil
}
func (mock *DeployController) GetDeployResult(ctx context.Context, in *protos.GetDeployResultRequest, opts ...grpc.CallOption) (*protos.GetDeployResultReply, error) {

	return &protos.GetDeployResultReply{
		Status: "successful",
		Err:    nil,
		Items: []*protos.DeployItemResult{
			{
				DeployItem: &protos.DeployItem{
					ItemName:            string(constant.MachineRoleMaster),
					NodeName:            "master1",
					FailureCanBeIgnored: false,
				},
				Status: "completed",
				Logs:   "",
			},
			{
				DeployItem: &protos.DeployItem{
					ItemName:            string(constant.MachineRoleEtcd),
					NodeName:            "master1",
					FailureCanBeIgnored: false,
				},
				Status: "completed",
				Logs:   "",
			},
		},
	}, nil
}

func (mock *DeployController) FetchKubeConfig(ctx context.Context, in *protos.FetchKubeConfigRequest, opts ...grpc.CallOption) (*protos.FetchKubeConfigReply, error) {
	return &protos.FetchKubeConfigReply{
		KubeConfig: []byte("kube config content")}, nil
}

func (mock *DeployController) CheckNetworkRequirements(
	ctx context.Context, in *protos.CheckNetworkRequirementRequest, opts ...grpc.CallOption) (
	*protos.CheckNetworkRequirementsReply, error) {
	return &protos.CheckNetworkRequirementsReply{
		Passed: true,
		Err:    nil,
		Nodes:  []*protos.NodeCheckResult{},
		Connectivities: []*protos.ConnectivityCheckResult{
			&protos.ConnectivityCheckResult{
				SourceNodeName:      "master1",
				DestinationNodeName: "master2",
				Status:              "successful",
				Err:                 nil,
				Items: []*protos.ItemCheckResult{
					&protos.ItemCheckResult{
						Item:   &protos.CheckItem{Name: "vxlan", Description: "connectivity of UDP port passing vxlan packets"},
						Status: "successful",
						Err:    nil,
						Logs:   "",
					},
				},
			},
		},
	}, nil
}

func (mock *DeployController) GetCheckNodesLog(ctx context.Context, in *protos.GetCheckNodesLogRequest,
	opts ...grpc.CallOption) (*protos.GetCheckNodesLogReply, error) {
	// To be implmented
	return nil, nil
}

func (mock *DeployController) GetDeployLog(ctx context.Context, in *protos.GetDeployLogRequest,
	opts ...grpc.CallOption) (*protos.GetDeployLogReply, error) {
	// To be implmented
	return nil, nil
}
