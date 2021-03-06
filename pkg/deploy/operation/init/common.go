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
	"bytes"

	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type ItemEnum string

type OperationsGenerator struct{}

type InitOperation interface {
	RunCommands(config *pb.Node, initAction *operation.NodeInitAction, logChan chan<- *bytes.Buffer) ([]byte, []byte, error)
}

const (
	FireWall   ItemEnum = "firewall"
	HostAlias  ItemEnum = "hostalias"
	HostName   ItemEnum = "hostname"
	Network    ItemEnum = "network"
	Route      ItemEnum = "route"
	Swap       ItemEnum = "swap"
	TimeZone   ItemEnum = "timezone"
	Haproxy    ItemEnum = "haproxy"
	Keepalived ItemEnum = "keepalived"
	KubeTool   ItemEnum = "kubetool"
)

const (
	DefaultCommonLibPath = "/scripts/lib.sh"
	HaDockerFilePath     = "/scripts/init_deploy_haproxy_keepalived/docker.sh"
	HaLibFilePath        = "/scripts/init_deploy_haproxy_keepalived/lib.sh"
	HaSystemdFilePath    = "/scripts/init_deploy_haproxy_keepalived/systemd.sh"
)

func NewInitOperations() *OperationsGenerator {
	return &OperationsGenerator{}
}

func (og *OperationsGenerator) CreateOperations(item ItemEnum, action *operation.NodeInitAction) InitOperation {
	switch item {
	case FireWall:
		return &InitFireWallOperation{}
	case HostAlias:
		return &InitHostaliasOperation{}
	case HostName:
		return &InitHostNameOperation{}
	case Network:
		return &InitNetworkOperation{}
	case Route:
		return &InitRouteOperation{}
	case Swap:
		return &InitSwapOperation{}
	case TimeZone:
		return &InitTimeZoneOperation{}
	case Haproxy:
		return &InitHaproxyOperation{}
	case Keepalived:
		return &InitKeepalivedOperation{}
	case KubeTool:
		return &InitKubeToolOperation{}
	default:
		return nil
	}
}

// group by master role
func groupByRole(arr []string, match string) bool {
	for _, item := range arr {
		if item == match {
			return true
		}
	}
	return false
}
