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
	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"

	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
)

// rollbackRelease inner function of running actions to rollback a helm release.
func rollbackRelease(c *gin.Context, r *api.HelmRelease) error {
	logEntry := log.ReqEntry(c)

	logEntry.Debugf("getting helm action config...")
	rollbackConfig, err := generateHelmActionConfig(r.Cluster, r.Namespace, logEntry)
	if err != nil {
		logEntry.WithField("cluster", r.Cluster).
			Warningf("failed to generate configuration for helm action")
		return err
	}
	rollbackAction := action.NewRollback(rollbackConfig)
	rollbackAction.Version = int(r.Revision)
	err = rollbackAction.Run(r.Name)
	if err != nil {
		logEntry.WithField("cluster", r.Cluster).WithField("namespace", r.Namespace).
			WithField("releaseName", r.Name).WithField("error", err.Error()).
			Warningf("failed to run rollback action")
	}
	return nil
}
