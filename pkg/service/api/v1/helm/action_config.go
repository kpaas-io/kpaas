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

package helm

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/kpaas-io/kpaas/pkg/service/kubeutils"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

const (
	HelmStorageSecrets = "secrets"
)

// generateHelmActionConfig generates the configuration to run helm actions.
func generateHelmActionConfig(cluster string, namespace string, logger *logrus.Entry) (
	*action.Configuration, error) {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	kubeConfigPath, err := kubeutils.KubeConfigPathForCluster(cluster)
	if err != nil {
		logger.WithField("cluster", cluster).WithField("error", err.Error()).
			Warning("failed to fetch kubeconfig for cluster")
		appErr := h.ENotFound.WithPayload(fmt.Sprintf(
			"cannot find kubeconfig for cluster %s", cluster))
		return nil, appErr
	}
	configFlag := genericclioptions.NewConfigFlags(false)
	configFlag.KubeConfig = &kubeConfigPath
	actionConfig := &action.Configuration{}
	actionConfig.Init(configFlag, namespace, HelmStorageSecrets, logger.Debugf)
	return actionConfig, nil
}
