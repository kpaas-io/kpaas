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
)

// Note: this file contains the env and test data to test the deploy gRPC API.
// Don't set _skip to true when commit for PR, otherwise it will break the UT.

var (
	_skip = true // set it to true if you want to run/debug the test locally.

	_launchLocalServer   = true           // launch a local deploy server and connect to it to test
	_remoteServerAddress = "0.0.0.0:8081" // connect to the remote deploy server
)

type ApiParams struct {
	request interface{}
	reply   interface{}
}

var _testdata = map[string]*ApiParams{
	"TestConnection": &ApiParams{
		request: &pb.TestConnectionRequest{
			Node: &pb.Node{
				Name: "g1-node0",
				Ip:   "47.102.123.5",
				Ssh: &pb.SSH{
					Port: 22,
					Auth: &pb.Auth{
						Type:       "password",
						Username:   "root",
						Credential: "*", // replace it with the correct password when run/debug testing.
					},
				},
			},
		},
		reply: &pb.TestConnectionReply{
			Passed: true,
		},
	},
}
