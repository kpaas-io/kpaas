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

// getRelease inner function of calling helm actions to get detailed info of a release
func getRelease(c *gin.Context, cluster string, namespace string, releaseName string) (
	*api.HelmRelease, error) {
	logEntry := log.ReqEntry(c)

	logEntry.Debug("getting action config...")
	getReleaseConfig, err := generateHelmActionConfig(cluster, namespace, logEntry)
	if err != nil {
		logEntry.WithField("cluster", cluster).
			Warningf("failed to generate configuration for helm action")
		return nil, err
	}
	getReleaseAction := action.NewGet(getReleaseConfig)
	getResult, err := getReleaseAction.Run(releaseName)
	if err != nil {
		logEntry.WithField("cluster", cluster).WithField("namespace", namespace).
			WithField("releaseName", releaseName).WithField("error", err).
			Warningf("failed to run get release action")
		return nil, err
	}
	// TODO: include more information from helm in returned result.
	res := &api.HelmRelease{
		Cluster:      cluster,
		Namespace:    namespace,
		Name:         releaseName,
		Revision:     uint32(getResult.Version),
		Chart:        getResult.Chart.Metadata.Name,
		ChartVersion: getResult.Chart.Metadata.Version,
	}
	return res, nil
}
