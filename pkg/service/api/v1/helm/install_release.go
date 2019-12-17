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
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"

	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
)

// installRelease inner function of calling helm actions to install a release
func installRelease(c *gin.Context, r *api.HelmRelease) (*api.HelmRelease, error) {
	logEntry := log.ReqEntry(c).
		WithField("cluster", r.Cluster).WithField("namspace", r.Namespace).WithField("releaseName", r.Name)

	// fetch kubeconfig for cluster
	logEntry.Debug("getting action config...")
	installConfig, err := generateHelmActionConfig(r.Cluster, r.Namespace, logEntry)
	if err != nil {
		logEntry.Warningf("failed to generate configuration for helm action")
		// generateHelmActionConfig returns h.AppErr, so we directly return err here
		return nil, err
	}
	installAction := action.NewInstall(installConfig)

	logEntry = logEntry.WithField("chart", r.Chart)
	logEntry.Debug("loading chart...")
	ch, err := loader.Load(r.Chart)
	if err != nil {
		logEntry.Warningf("failed to load chart")
		appErr := h.ENotFound.WithPayload(fmt.Sprintf("chart '%s' not found", r.Chart))
		return nil, appErr
	}

	installAction.Namespace = r.Namespace
	installAction.ReleaseName = r.Name
	if r.Name == "" {
		installAction.GenerateName = true
		installAction.ReleaseName = generateReleaseName(ch)
	}
	logEntry.Debug("running installation...")
	installResult, err := installAction.Run(ch, r.Values)
	if err != nil {
		logEntry.WithField("error", err).Warning("failed to run install action")
		// TODO: analyze errors happened in running installAction.Run and return proper AppErr
		return nil, fmt.Errorf("failed to run install action")
	}
	res := &api.HelmRelease{
		Cluster:      r.Cluster,
		Namespace:    r.Namespace,
		Name:         installResult.Name,
		Chart:        installResult.Chart.Metadata.Name,
		ChartVersion: installResult.Chart.Metadata.Version,
		Revision:     uint32(installResult.Version),
	}
	return res, nil
}

func generateReleaseName(ch *chart.Chart) string {
	if ch == nil || ch.Metadata == nil {
		return ""
	}
	return ch.Metadata.Name + "-" + strings.Replace(ch.Metadata.Version, ".", "-", -1) +
		"-" + strconv.Itoa(int(time.Now().Unix()))
}
