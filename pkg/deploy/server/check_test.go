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

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/deploy/task"
)

func TestCheckItemToItemCheckResult(t *testing.T) {
	tests := []struct {
		input *action.NodeCheckItem
		want  *pb.ItemCheckResult
	}{
		{
			input: nil,
			want:  nil,
		},
		{
			input: &action.NodeCheckItem{
				Name:        "test checkitem name",
				Description: "test checkitem description",
				Status:      action.ItemActionDone,
				Err: &pb.Error{
					Reason: "checkitem reason",
					Detail: "checkitem detail",
				},
			},
			want: &pb.ItemCheckResult{
				Item: &pb.CheckItem{
					Name:        "test checkitem name",
					Description: "test checkitem description",
				},
				Status: action.ItemActionDone,
				Err: &pb.Error{
					Reason: "checkitem reason",
					Detail: "checkitem detail",
				},
			},
		},
	}

	for _, tt := range tests {
		result := checkItemToItemCheckResult(tt.input)
		assert.Equal(t, tt.want, result)
	}
}

func TestCheckActionToNodeCheckResult(t *testing.T) {
	memCheckItemError := &pb.Error{
		Reason:     "test reason",
		Detail:     "test detail",
		FixMethods: "test fixmethod",
	}
	tests := []struct {
		input *action.NodeCheckAction
		want  *pb.NodeCheckResult
	}{
		{
			input: nil,
			want:  nil,
		},
		{
			input: &action.NodeCheckAction{
				Base: action.Base{
					Node: nil,
				},
			},
			want: nil,
		},
		{
			input: &action.NodeCheckAction{
				Base: action.Base{
					Node: &pb.Node{
						Name: "node1",
					},
					Status: action.ActionDone,
				},
				CheckItems: []*action.NodeCheckItem{
					&action.NodeCheckItem{
						Name:        "check cpucore",
						Description: "check cpucore description",
						Status:      action.ItemActionDone,
					},
					&action.NodeCheckItem{
						Name:        "check memroy",
						Description: "check memroy description",
						Status:      action.ItemActionFailed,
						Err:         memCheckItemError,
					},
				},
			},
			want: &pb.NodeCheckResult{
				NodeName: "node1",
				Status:   string(action.ActionDone),
				Err:      nil,
				Items: []*pb.ItemCheckResult{
					&pb.ItemCheckResult{
						Item: &pb.CheckItem{
							Name:        "check cpucore",
							Description: "check cpucore description",
						},
						Status: action.ItemActionDone,
						Err:    nil,
					},
					&pb.ItemCheckResult{
						Item: &pb.CheckItem{
							Name:        "check memroy",
							Description: "check memroy description",
						},
						Status: action.ItemActionFailed,
						Err:    memCheckItemError,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		result := checkActionToNodeCheckResult(tt.input)
		assert.Equal(t, tt.want, result)
	}
}

func TestConnectivityCheckToNodeCheckResult(t *testing.T) {
	tests := []struct {
		input *action.ConnectivityCheckAction
		want  *pb.NodeCheckResult
	}{
		{
			input: nil,
			want:  nil,
		},
		{
			input: &action.ConnectivityCheckAction{
				Base: action.Base{
					Node: nil,
				},
			},
			want: nil,
		},
		{
			input: &action.ConnectivityCheckAction{
				Base: action.Base{
					Node: &pb.Node{
						Name: "node1",
					},
					Status: action.ActionDone,
				},
				CheckItems: []action.ConnectivityCheckItem{
					action.ConnectivityCheckItem{
						Protocol: consts.ProtocolTCP,
						Port:     uint16(1234),
						CheckResult: &pb.ItemCheckResult{
							Item:   &pb.CheckItem{Name: "tcp-1234", Description: "check TCP port 1234"},
							Status: action.ItemActionDone,
							Err:    nil,
						},
					},
				},
			},
			want: &pb.NodeCheckResult{
				NodeName: "node1",
				Status:   "Done",
				Err:      nil,
				Items: []*pb.ItemCheckResult{
					&pb.ItemCheckResult{
						Item:   &pb.CheckItem{Name: "tcp-1234", Description: "check TCP port 1234"},
						Status: "done",
						Err:    nil,
					},
				},
			},
		},
	}

	for _, testCase := range tests {
		assert.Equal(t, testCase.want, connectivityCheckToNodeCheckResult(testCase.input))
	}
}

func TestGetNodeCheckResult(t *testing.T) {
	memCheckItemError := &pb.Error{
		Reason:     "test reason",
		Detail:     "test detail",
		FixMethods: "test fixmethod",
	}

	nodeCheckAction := &action.NodeCheckAction{
		Base: action.Base{
			Node: &pb.Node{
				Name: "node1",
			},
			Status: action.ActionDone,
		},
		CheckItems: []*action.NodeCheckItem{
			&action.NodeCheckItem{
				Name:        "check cpucore",
				Description: "check cpucore description",
				Status:      action.ItemActionDone,
			},
			&action.NodeCheckItem{
				Name:        "check memroy",
				Description: "check memroy description",
				Status:      action.ItemActionFailed,
				Err:         memCheckItemError,
			},
		},
	}

	tests := []struct {
		input     task.Task
		wantReply *pb.GetCheckNodesResultReply
		wantError error
	}{
		{
			input: &task.NodeCheckTask{
				Base: task.Base{
					Actions: []action.Action{nodeCheckAction},
					Status:  task.TaskDone,
					Err:     nil,
				},
			},
			wantReply: &pb.GetCheckNodesResultReply{
				Status: string(task.TaskDone),
				Err:    nil,
				Nodes: []*pb.NodeCheckResult{
					&pb.NodeCheckResult{
						NodeName: "node1",
						Status:   string(action.ActionDone),
						Err:      nil,
						Items: []*pb.ItemCheckResult{
							&pb.ItemCheckResult{
								Item: &pb.CheckItem{
									Name:        "check cpucore",
									Description: "check cpucore description",
								},
								Status: action.ItemActionDone,
								Err:    nil,
							},
							&pb.ItemCheckResult{
								Item: &pb.CheckItem{
									Name:        "check memroy",
									Description: "check memroy description",
								},
								Status: action.ItemActionFailed,
								Err:    memCheckItemError,
							},
						},
					},
				},
			},
			wantError: nil,
		},
	}
	c := &controller{}
	for _, testCase := range tests {
		reply, err := c.getCheckNodeResult(testCase.input, false)
		assert.Equal(t, testCase.wantReply, reply)
		assert.Equal(t, testCase.wantError, err)
	}

}
