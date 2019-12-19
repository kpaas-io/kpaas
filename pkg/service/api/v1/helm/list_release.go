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

func listRelease(c *gin.Context, cluster string, namespace string) ([]*api.HelmRelease, error) {
	logEntry := log.ReqEntry(c).
		WithField("cluster", cluster).WithField("namespace", namespace)

	logEntry.Debug("getting action config...")
	listReleaseConfig, err := generateHelmActionConfig(cluster, namespace, logEntry)
	if err != nil {
		logEntry.WithField("error", err).Warning("failed to generate configuration for helm action")
		return nil, err
	}
	// TODO: add more options for listing releases, such as filter, sorting method.
	listReleaseAction := action.NewList(listReleaseConfig)
	listResult, err := listReleaseAction.Run()
	if err != nil {
		logEntry.Warning("failed to run list action")
		return nil, err
	}
	ret := []*api.HelmRelease{}
	for _, releaseContent := range listResult {
		r := &api.HelmRelease{
			Cluster:      cluster,
			Namespace:    namespace,
			Name:         releaseContent.Name,
			Revision:     uint32(releaseContent.Version),
			Chart:        releaseContent.Chart.Metadata.Name,
			ChartVersion: releaseContent.Chart.Metadata.Version,
		}
		ret = append(ret, r)
	}

	return ret, nil
}
