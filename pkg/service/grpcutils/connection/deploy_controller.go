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

package connection

import (
	"sync"

	"google.golang.org/grpc"

	"github.com/kpaas-io/kpaas/pkg/service/config"
)

var (
	deployControllerConnection *grpc.ClientConn
	deployControllerRWLock     *sync.RWMutex
)

func init() {
	deployControllerRWLock = new(sync.RWMutex)
}

func GetDeployControllerConnection() (*grpc.ClientConn, error) {

	clientConn := tryToGetConnection()
	if clientConn != nil {
		return clientConn, nil
	}

	deployControllerRWLock.Lock()
	defer deployControllerRWLock.Unlock()
	if deployControllerConnection != nil {
		return deployControllerConnection, nil
	}

	var err error
	deployControllerConnection, err = grpc.Dial(config.Config.DeployController.GetAddress(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return deployControllerConnection, nil
}

func tryToGetConnection() *grpc.ClientConn {

	deployControllerRWLock.RLock()
	defer deployControllerRWLock.RUnlock()
	return deployControllerConnection
}

func Close() error {

	if deployControllerConnection == nil {
		return nil
	}

	return deployControllerConnection.Close()
}
