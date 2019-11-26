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
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/deploy/task"
)

type Interface interface {
	Run(stopCh <-chan struct{}) error
}

type ServerOptions struct {
	Port       uint16
	LogFileLoc string
}

type server struct {
	port       uint16
	logFileLoc string
}

func New(options ServerOptions) Interface {
	return &server{
		port:       options.Port,
		logFileLoc: options.LogFileLoc,
	}
}

func (s *server) Run(stopCh <-chan struct{}) error {
	gRpcSvr := grpc.NewServer()

	// use the map cache store
	store := task.GetGlobalCacheStore()
	protos.RegisterDeployContollerServer(gRpcSvr, &controller{
		store:      store,
		logFileLoc: s.logFileLoc,
	})
	reflection.Register(gRpcSvr)

	listenAddr := fmt.Sprintf("0.0.0.0:%d", s.port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %s", listenAddr, err)
	}

	logrus.Info("Begin to serve.")
	go gRpcSvr.Serve(listener)

	<-stopCh

	return nil
}
