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

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"

	"github.com/kpaas-io/kpaas/pkg/utils/log"
)

func uninstallRelease(
	c *gin.Context, cluster string, namespace string, releaseName string) error {
	logEntry := log.ReqEntry(c).
		WithField("cluster", cluster).WithField("namespace", namespace).WithField("releaseName", releaseName)

	logEntry.Debugf("getting helm action config...")
	uninstallConfig, err := generateHelmActionConfig(cluster, namespace, logEntry)
	if err != nil {
		logEntry.WithField("error", err).
			Warningf("failed to generate configuration for helm action")
		return err
	}

	uninstallAction := action.NewUninstall(uninstallConfig)
	// TODO: analyze response from helm action and give more info in response
	_, err = uninstallAction.Run(releaseName)
	if err != nil {
		logEntry.WithField("error", err).Warning("failed to run uninstall action")
		// TODO: analyze error and return proper appError
		return fmt.Errorf("failed to run uninstall action")
	}
	return nil
}
