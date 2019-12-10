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

package action

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type deployEtcdStatus string

// DeployEtcdActionConfig represents the config for a ectd deploy in a node
type DeployEtcdActionConfig struct {
	CaCrt           *x509.Certificate
	CaKey           crypto.Signer
	Node            *pb.Node
	ClusterNodes    []*pb.Node
	LogFileBasePath string
}

type deployEtcdAction struct {
	base
	caCrt        *x509.Certificate
	caKey        crypto.Signer
	node         *pb.Node
	clusterNodes []*pb.Node
}

// NewDeployEtcdAction returns a deploy etcd action based on the config.
// User should use this function to create a deploy etcd action.
func NewDeployEtcdAction(cfg *DeployEtcdActionConfig) (Action, error) {
	var err error
	if cfg == nil {
		err = fmt.Errorf("action config is nil")
	} else if cfg.Node == nil {
		err = fmt.Errorf("invalid node check config: node is nil")
	}

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	actionName := getDeployEtcdActionName(cfg)
	return &deployEtcdAction{
		base: base{
			name:              actionName,
			actionType:        ActionTypeDeployEtcd,
			status:            ActionPending,
			logFilePath:       GenActionLogFilePath(cfg.LogFileBasePath, actionName),
			creationTimestamp: time.Now(),
		},
		caCrt:        cfg.CaCrt,
		caKey:        cfg.CaKey,
		node:         cfg.Node,
		clusterNodes: cfg.ClusterNodes,
	}, nil
}

func getDeployEtcdActionName(cfg *DeployEtcdActionConfig) string {
	// used the node name as the the action name for now, this may be changed in the future.
	return cfg.Node.GetName()
}
