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
	"bytes"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	RegisterExecutor(ActionTypeFetchKubeConfig, new(fetchKubeConfigExecutor))
}

type fetchKubeConfigExecutor struct {
}

func (a *fetchKubeConfigExecutor) Execute(act Action) *pb.Error {
	kubeCfgAction, ok := act.(*FetchKubeConfigAction)
	if !ok {
		return errOfTypeMismatched(new(FetchKubeConfigAction), act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	var pbErr *pb.Error

	defer func() {
		deploy.PBErrLogger(pbErr, logger).Debug()

		// TODO: write some information to action log file.
	}()

	logger.Debug("Start to execute action")

	m, err := machine.NewMachine(kubeCfgAction.Node)
	if err != nil {
		pbErr = &pb.Error{
			Reason: "failed to connect to target node",
			Detail: err.Error(),
		}
		return pbErr
	}
	defer m.Close()

	var buf bytes.Buffer
	if err = m.FetchFile(&buf, consts.KubeConfigPath); err != nil {
		pbErr = &pb.Error{
			Reason: "failed to fetch kube config",
			Detail: err.Error(),
		}
		return pbErr
	}

	// Update action
	kubeCfgAction.KubeConfig = buf.Bytes()

	logger.Debug("Finsih to execute action")
	return nil
}
