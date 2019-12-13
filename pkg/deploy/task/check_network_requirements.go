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

package task

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// CheckNetworkRequirementsTaskConfig configuration to create a CheckNetworkRequirements task.
type CheckNetworkRequirementsTaskConfig struct {
	Nodes           []*pb.Node
	NetworkOptions  *pb.NetworkOptions
	LogFileBasePath string
}

type CheckNetworkRequirementsTask struct {
	Base

	Nodes          []*pb.Node
	NetworkOptions *pb.NetworkOptions
}

// NewCheckNetworkRequirementsTask create a CheckNetworkRequirements task to check prerequisites of deploying network.
func NewCheckNetworkRequirementsTask(
	name string, config *CheckNetworkRequirementsTaskConfig) (Task, error) {
	if config == nil {
		return nil, fmt.Errorf("Invalid task config: empty config")
	}
	if len(config.Nodes) == 0 {
		return nil, fmt.Errorf("Invalid task config: node list empty")
	}
	return &CheckNetworkRequirementsTask{
		Base: Base{
			Name:              name,
			TaskType:          TaskTypeCheckNetworkRequirements,
			Status:            TaskPending,
			LogFilePath:       GenTaskLogFilePath(config.LogFileBasePath, name),
			CreationTimestamp: time.Now(),
		},
		Nodes:          config.Nodes,
		NetworkOptions: config.NetworkOptions,
	}, nil
}

type checkNetworkRequirementsProcessor struct{}

func (p *checkNetworkRequirementsProcessor) SplitTask(task Task) error {
	checkNetworkRequirementsTask, ok := task.(*CheckNetworkRequirementsTask)
	if !ok {
		return fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, task)
	}
	// set a default network option if not specified.
	if checkNetworkRequirementsTask.NetworkOptions == nil {
		checkNetworkRequirementsTask.NetworkOptions = &pb.NetworkOptions{
			NetworkType: string(consts.NetworkTypeCalico),
		}
	}
	var actions []action.Action
	switch checkNetworkRequirementsTask.NetworkOptions.NetworkType {
	case string(consts.NetworkTypeCalico):
		var err error
		if checkNetworkRequirementsTask.NetworkOptions.CalicoOptions == nil {
			checkNetworkRequirementsTask.NetworkOptions.CalicoOptions = &pb.CalicoOptions{
				CheckConnectivityAll: false,
				EncapsulationMode:    "vxlan",
				VxlanPort:            4789,
			}
		}
		actions, err = p.splitActionsCalico(checkNetworkRequirementsTask)
		checkNetworkRequirementsTask.Actions = actions
		if err != nil {
			return fmt.Errorf("failed to split actions, error %v", err)
		}
	default:
		return fmt.Errorf("unsupported network type: %s",
			checkNetworkRequirementsTask.NetworkOptions.NetworkType)
	}

	return nil
}

func (p *checkNetworkRequirementsProcessor) splitActionsCalico(
	task *CheckNetworkRequirementsTask) ([]action.Action, error) {
	// randomly choose a "peer" for each node.
	numNodes := len(task.Nodes)
	if numNodes == 1 {
		return []action.Action{}, nil
	}
	actions := []action.Action{}
	for i, node := range task.Nodes {
		randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
		// choose the index of peer. If index of itself is chosen, use the last node instead.
		peerIndex := randGen.Intn(numNodes - 1)
		if peerIndex == i {
			peerIndex = numNodes - 1
		}
		connectivityCheckAction, err := makeConnectivityCheckActionCalico(
			node, task.Nodes[peerIndex], task.NetworkOptions.CalicoOptions)
		if err != nil {
			return []action.Action{}, fmt.Errorf("failed to split task into actions, error %v", err)
		}
		actions = append(actions, connectivityCheckAction)
	}
	return actions, nil
}

func makeConnectivityCheckActionCalico(
	src *pb.Node, dst *pb.Node, calicoOptions *pb.CalicoOptions) (action.Action, error) {
	if src != nil {
		return nil, fmt.Errorf("source node empty")
	}
	if dst == nil {
		return nil, fmt.Errorf("destination node empty")
	}
	if calicoOptions == nil {
		return nil, fmt.Errorf("calico options empty")
	}
	cfg := &action.ConnectivityCheckActionConfig{
		SourceNode:      src,
		DestinationNode: dst,
		ConnectivityCheckItems: []action.ConnectivityCheckItem{
			action.ConnectivityCheckItem{
				Protocol: consts.ProtocolTCP,
				Port:     uint16(179),
				CheckResult: &pb.ItemCheckResult{
					Item: &pb.CheckItem{
						Name:        fmt.Sprintf("connectivity-check-BGP"),
						Description: "check connectivity to BGP port",
					},
					Status: action.ItemActionPending,
				},
			},
			action.ConnectivityCheckItem{
				Protocol: consts.ProtocolTCP,
				Port:     uint16(6443),
				CheckResult: &pb.ItemCheckResult{
					Item: &pb.CheckItem{
						Name:        fmt.Sprintf("connectivity-check-kube-API"),
						Description: "check connectivity to kubernetes API port",
					},
					Status: action.ItemActionPending,
				},
			},
			action.ConnectivityCheckItem{
				Protocol: consts.ProtocolUDP,
				Port:     uint16(calicoOptions.VxlanPort),
				CheckResult: &pb.ItemCheckResult{
					Item: &pb.CheckItem{
						Name:        fmt.Sprintf("connectivity-check-vxlan"),
						Description: "check connectivity for vxlan packets",
					},
					Status: action.ItemActionPending,
				},
			},
		},
	}
	return action.NewConnectivityCheckAction(cfg)
}
