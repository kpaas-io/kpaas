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

type checkNetworkRequirementsTask struct {
	base
	nodes          []*pb.Node
	networkOptions *pb.NetworkOptions
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
	return &checkNetworkRequirementsTask{
		base: base{
			name:              name,
			taskType:          TaskTypeCheckNetworkRequirements,
			status:            TaskPending,
			logFilePath:       GenTaskLogFilePath(config.LogFileBasePath, name),
			creationTimestamp: time.Now(),
		},
		nodes:          config.Nodes,
		networkOptions: config.NetworkOptions,
	}, nil
}

type checkNetworkRequirementsProcessor struct{}

func (p *checkNetworkRequirementsProcessor) SplitTask(task Task) error {
	checkNetworkRequirementsTask, ok := task.(*checkNetworkRequirementsTask)
	if !ok {
		return fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, task)
	}
	// set a default network option if not specified.
	if checkNetworkRequirementsTask.networkOptions == nil {
		checkNetworkRequirementsTask.networkOptions = &pb.NetworkOptions{
			NetworkType: string(consts.NetworkTypeCalico),
		}
	}
	var actions []action.Action
	switch checkNetworkRequirementsTask.networkOptions.NetworkType {
	case string(consts.NetworkTypeCalico):
		var err error
		if checkNetworkRequirementsTask.networkOptions.CalicoOptions == nil {
			checkNetworkRequirementsTask.networkOptions.CalicoOptions = &pb.CalicoOptions{
				CheckConnectivityAll: false,
				EncapsulationMode:    "vxlan",
				VxlanPort:            4789,
			}
		}
		actions, err = p.splitActionsCalico(checkNetworkRequirementsTask)
		checkNetworkRequirementsTask.actions = actions
		if err != nil {
			return fmt.Errorf("failed to split actions, error %v", err)
		}
	default:
		return fmt.Errorf("unsupported network type: %s",
			checkNetworkRequirementsTask.networkOptions.NetworkType)
	}

	return nil
}

func (p *checkNetworkRequirementsProcessor) splitActionsCalico(
	task *checkNetworkRequirementsTask) ([]action.Action, error) {
	// randomly choose a "peer" for each node.
	numNodes := len(task.nodes)
	if numNodes == 1 {
		return []action.Action{}, nil
	}
	actions := []action.Action{}
	for i, node := range task.nodes {
		randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
		// choose the index of peer. If index of itself is chosen, use the last node instead.
		peerIndex := randGen.Intn(numNodes - 1)
		if peerIndex == i {
			peerIndex = numNodes - 1
		}
		cfg := &action.ConnectivityCheckActionConfig{
			SourceNode:      node,
			DestinationNode: task.nodes[peerIndex],
			ConnectivityCheckItems: []action.ConnectivityCheckItem{
				action.ConnectivityCheckItem{
					Protocol: consts.ProtocolTCP,
					Port:     uint16(179),
				},
				action.ConnectivityCheckItem{
					Protocol: consts.ProtocolTCP,
					Port:     uint16(6443),
				},
				action.ConnectivityCheckItem{
					Protocol: consts.ProtocolUDP,
					Port:     uint16(task.networkOptions.CalicoOptions.VxlanPort),
				},
			},
		}
		connectivityCheckAction, err := action.NewConnectivityCheckAction(cfg)
		if err != nil {
			return []action.Action{}, fmt.Errorf("failed to split task into actions, error %v", err)
		}
		actions = append(actions, connectivityCheckAction)
	}
	return actions, nil
}
