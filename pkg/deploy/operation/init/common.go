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

package init

import (
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type ItemEnum int

type OperationsGenerator struct{}

type InitOperations struct {
	Script     string
	ScriptPath string
}

type InitAction interface {
	GetOperations(config *pb.Node) (operation.Operation, error)
	getScript() string
	getScriptPath() string
}

const (
	FireWall ItemEnum = iota
	HostAlias
	HostName
	Network
	Route
	Swap
	TimeZone
	Haproxy
	Keepalived
)

func NewInitOperations() *OperationsGenerator {
	return &OperationsGenerator{}
}

func (og *OperationsGenerator) CreateOperations(item ItemEnum) InitAction {
	switch item {
	case FireWall:
		return &InitFireWallOperation{}
	//case HostAlias:
	//	return &InitHostAliasOperation{}
	//case HostName:
	//	return &InitHostNameOperation{}
	//case Network:
	//	return &InitNetworkOperation{}
	//case Route:
	//	return &InitRoutOperation{}
	//case Swap:
	//	return &InitSwapOperation{}
	//case TimeZone:
	//	return &InitTimeZoneOperation{}
	//case Haproxy:
	//	return &InitHaproxyOperation{}
	//case Keepalived:
	//	return &InitKeepalivedOperation{}
	// TODO setup.sh for init kubeadm kubectl kubelet @yangruiray
	default:
		return nil
	}
}
