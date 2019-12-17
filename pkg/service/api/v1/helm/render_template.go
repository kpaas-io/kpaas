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

// renderTemplate is the inner function to render a template of a chart, returning the rendered YAML manifests.
func renderTemplate(c *gin.Context, r *api.HelmRelease) (string, error) {
	logEntry := log.ReqEntry(c)

	logEntry.Debug("getting action config...")
	installConfig, err := generateHelmActionConfig(r.Cluster, r.Namespace, logEntry)
	if err != nil {
		logEntry.WithField("cluster", r.Cluster).
			Warningf("failed to generate configuration for helm action")
		// generateHelmActionConfig returns h.AppErr, so we directly return err here
		return "", err
	}
	// rendering the template is just dry-running the install action and retieving the manifest in the result.
	installAction := action.NewInstall(installConfig)

	logEntry.WithField("chartPath", r.Chart).Debug("loading chart...")
	ch, err := loader.Load(r.Chart)
	if err != nil {
		logEntry.WithField("chart", r.Chart).
			WithField("error", err.Error()).Warningf("failed to load chart")
		appErr := h.ENotFound.WithPayload(fmt.Sprintf("chart '%s' not found", r.Chart))
		return "", appErr
	}

	installAction.Namespace = r.Namespace
	if r.Name != "" {
		installAction.ReleaseName = r.Name
	} else {
		installAction.ReleaseName = "RELEASE-NAME"
	}
	installAction.DryRun = true
	installAction.ClientOnly = false
	res, err := installAction.Run(ch, r.Values)
	if err != nil {
		logEntry.WithField("cluster", r.Cluster).WithField("namespace", r.Namespace).
			WithField("chart", r.Chart).WithField("error", err.Error()).
			Warning("failed to run install action")
		// TODO: analyze errors happened in running installAction.Run and return proper AppErr
		return "", fmt.Errorf("failed to run install action")
	}
	return res.Manifest, nil
}
