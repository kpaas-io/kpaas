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
	"helm.sh/helm/v3/pkg/chart/loader"

	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
)

func upgradeRelease(c *gin.Context, r *api.HelmRelease) (*api.HelmRelease, error) {
	logEntry := log.ReqEntry(c).
		WithField("cluster", r.Cluster).WithField("namespace", r.Namespace).WithField("releaseName", r.Name)

	// get helm action config for cluster
	logEntry.Debug("getting helm action config...")
	upgradeConfig, err := generateHelmActionConfig(r.Cluster, r.Namespace, logEntry)
	if err != nil {
		logEntry.WithField("error", err).Warningf("failed to generate configuration for helm action")
		return nil, err
	}
	upgradeAction := action.NewUpgrade(upgradeConfig)
	upgradeAction.Install = false
	upgradeAction.Namespace = r.Namespace

	// load chart
	// TODO: allow empty chart in request to use the chart used in current version of release
	logEntry = logEntry.WithField("chart", r.Chart)
	logEntry.Debug("loading chart..")
	ch, err := loader.Load(r.Chart)
	if err != nil {
		logEntry.WithField("error", err.Error()).Warningf("failed to load chart")
		appErr := h.ENotFound.WithPayload(fmt.Sprintf("chart '%s' not found", r.Chart))
		return nil, appErr
	}

	upgradeResult, err := upgradeAction.Run(r.Name, ch, r.Values)
	if err != nil {
		// TODO: analyze errors happened in running upgradeAction.Run and return proper AppErr
		logEntry.WithField("error", err).Warning("failed to run upgrade action")
		return nil, fmt.Errorf("failed to run upgrade action")
	}
	res := &api.HelmRelease{
		Cluster:      r.Cluster,
		Namespace:    r.Namespace,
		Name:         upgradeResult.Name,
		Chart:        upgradeResult.Chart.Metadata.Name,
		ChartVersion: upgradeResult.Chart.Metadata.Version,
		Revision:     uint32(upgradeResult.Version),
	}
	return res, nil
}
