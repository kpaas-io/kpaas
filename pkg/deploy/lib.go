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
	"fmt"
	"io"
	"os"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const (
	defaultApiServerPort = 6443
)

var (
	AllFilesNeeded = func(path string) bool {
		return true
	}
)

func FileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func MustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		panic(err)
	}
}

// GetControlPlaneEndpoint return control plane endpoint address
func GetControlPlaneEndpoint(clusterConfig *pb.ClusterConfig, masterNodes []*pb.Node) (addr string, err error) {
	conn := clusterConfig.KubeAPIServerConnect
	if conn == nil {
		return "", fmt.Errorf("nil %T encountered", conn)
	}

	// type could be ["firstMasterIP", "keepalived", "loadbalancer"]
	switch conn.Type {
	case "firstMasterIP":
		ip := masterNodes[0].Ip
		if ip == "" {
			err = fmt.Errorf("failed to get first master ip")
			return
		}
		addr = fmt.Sprintf("%v:%v", ip, defaultApiServerPort)
	case "keepalived":
		addr = conn.Keepalived.Vip
	case "loadbalancer":
		addr = fmt.Sprintf("%v:%v", conn.Loadbalancer.Ip, conn.Loadbalancer.Port)
	default:
		err = fmt.Errorf("unrecognized apiserver connect type: %v", conn.Type)
	}

	return
}
