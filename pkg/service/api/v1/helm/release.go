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

	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
)

const (
	ParamCluster   = "cluster"
	ParamNamespace = "namespace"
	ParamName      = "name"
)

// InstallRelease install a helm release from specified chart and values
// @ID InstallRelease
// @Summary Install helm release
// @Description Install a helm release in a cluster from specified chart and values
// @Tags helm
// @Produce application/json
// @Param cluster path string true "cluster to create release in"
// @Param namespace path string true "kubernetes namespace in that cluster to create release in"s
// @Success 201 {object} api.HelmRelease
// @Failure 409 {object} h.AppErr
// @Failure 400 {object} h.AppErr
// @Router  /api/v1/helm/clusters/{cluster}/namespaces/{namespace}/releases [post]
func InstallRelease(c *gin.Context) {
	release := api.HelmRelease{}
	err := parseRelease(c, &release)
	if err != nil {
		return
	}
	if release.Chart == "" {
		appErr := h.EParamsError.WithPayload("empty chart")
		h.E(c, appErr)
		return
	}
	res, err := RunInstallReleaseAction(c, &release)
	if err != nil {
		h.E(c, err)
		return
	}
	h.R(c, res)
}

// UpgradeRelease upgrade a helm release by its name.
// @ID UpgradeRelease
// @Summary Upgrade installed release
// @Description Upgrade an installed release with new chart and/or values
// @Tags helm
// @Produce application/json
// @Param cluster path string  true "cluster to upgrade release in"
// @Param namespace path string true "kubernetes namespace in that cluster to upgrade release in"
// @Param name path string true "name of release to upgrade"
// @Success 200 {object} api.HelmRelease
// @Failure 404 {object} h.AppErr
// @Failure 400 {object} h.AppErr
// @Router /api/v1/helm/clusters/{cluster}/namespaces/{namespace}/releases/{name} [put]
func UpgradeRelease(c *gin.Context) {
	release := api.HelmRelease{}
	err := parseRelease(c, &release)
	if err != nil {
		return
	}
	if release.Chart == "" {
		appErr := h.EParamsError.WithPayload("empty chart")
		h.E(c, appErr)
		return
	}
	res, err := upgradeRelease(c, &release)
	if err != nil {
		h.E(c, err)
		return
	}
	h.R(c, res)
}

// RollbackRelease rollbacks a release to a certain version.
// @ID RollbackRelease
// @Summary rollback a release
// @Description roll back a release to an earlier version
// @Tags helm
// @Produce application/json
// @Param cluster path string  true "cluster to rollback release in"
// @Param namespace path string true "kubernetes namespace in that cluster to rollback release in"
// @Param name path string true "name of release to rollback"
// @Success 200 {object} api.SuccessfulOption
// @Failure 404 {object} h.AppErr
// @Failure 400 {object} h.AppErr
// @Router /api/v1/helm/clusters/{cluster}/namespaces/{namespace}/releases/{name}/rollback [put]
func RollbackRelease(c *gin.Context) {
	release := api.HelmRelease{}
	err := parseRelease(c, &release)
	if err != nil {
		return
	}
	err = rollbackRelease(c, &release)
	if err != nil {
		h.E(c, err)
		return
	}
	h.R(c, api.SuccessfulOption{Success: true})
}

// GetRelease download all information for a named release
// @ID GetRelease
// @Summary get information of a named release
// @Description get manifest, chart, and values of a named release
// @Tags helm
// @Produce application/json
// @Param cluster path string  true "kubernetes cluster where the release is"
// @Param namespace path string true "kubernetes namespace where the release is"
// @Param name path string true "name of release"
// @Success 200 {object} api.HelmRelease
// @Failure 404 {object} h.AppErr
// @Router /api/v1/helm/clusters/{cluster}/namespaces/{namespace}/releases/{name} [get]
func GetRelease(c *gin.Context) {
	cluster := c.Param(ParamCluster)
	namespace := c.Param(ParamNamespace)
	releaseName := c.Param(ParamName)

	res, err := getRelease(c, cluster, namespace, releaseName)
	if err != nil {
		h.E(c, err)
		return
	}
	h.R(c, res)
}

// ListRelease list all releases in a namespace.
// @ID ListRelease
// @Summary list releases
// @Description list all releases in a namespace
// @Tags helm
// @Produce application/json
// @Param cluster path string  true "kubernetes cluster to list releases in"
// @Param namespace path string true "kubernetes namespace to list release in"
// @Success 200 {array} api.HelmRelease
// @Failure 404 {object} h.AppErr
// @Router /api/v1/helm/clusters/{cluster}/namespaces/{namespace}/releases [get]
func ListRelease(c *gin.Context) {
	cluster := c.Param(ParamCluster)
	namespace := c.Param(ParamNamespace)
	releases, err := listRelease(c, cluster, namespace)
	if err != nil {
		h.E(c, err)
		return
	}
	h.R(c, releases)
}

// UninstallRelease uninstalls a named release.
// @ID UninstallRelease
// @Summary uninstall a release
// @Description uninstall a named release and deleted all resources in kubernetes created for the release
// @Tags helm
// @Produce application/json
// @Param cluster path string  true "kubernetes cluster where the release is"
// @Param namespace path string true "kubernetes namespace where the release is"
// @Param name path string true "name of release"
// @Success 204
// @Failure 404 {object} h.AppErr
// @Router /api/v1/helm/clusters/{cluster}/namespaces/{namespace}/releases/{name} [delete]
func UninstallRelease(c *gin.Context) {
	cluster := c.Param(ParamCluster)
	namespace := c.Param(ParamNamespace)
	releaseName := c.Param(ParamName)
	err := uninstallRelease(c, cluster, namespace, releaseName)
	if err != nil {
		h.E(c, err)
		return
	}
	h.R(c, nil)
}

// ExportRelease exports manifests of a named release, in yaml/json format.
// @ID ExportRelease
// @Summary export manifests of a release
// @Description export manifests of a release in yaml/json format
// @Produce application/x-yaml
// @Produce application/json
// @Param cluster path string  true "kubernetes cluster where the release is"
// @Param namespace path string true "kubernetes namespace where the release is"
// @Param name path string true "name of release"
// @Success 200
// @Failure 404 {object} h.AppErr
// @Router /api/v1/helm/clusters/{cluster}/namespaces/{namespace}/releases/{name}/export [get]
func ExportRelease(c *gin.Context) {
	cluster := c.Param(ParamCluster)
	namespace := c.Param(ParamNamespace)
	releaseName := c.Param(ParamName)
	manifest, err := exportRelease(c, cluster, namespace, releaseName)
	if err != nil {
		h.E(c, err)
		return
	}
	c.Data(200, gin.MIMEYAML, []byte(manifest))
}

// RenderTemplate render chart templates locally and display the output.
// @ID RenderTemplate
// @Summary render templates in a chart
// @Description render chart templates locally and display the output
// @Tags helm
// @Produce application/x-yaml
// @Produce application/json
// @Success 201
// @Failure 404 {object} h.AppErr
// @Failure 400 {object} h.AppErr
// @Router /api/v1/helm/render [post]
func RenderTemplate(c *gin.Context) {
	release := api.HelmRelease{}
	err := c.ShouldBindJSON(&release)
	if err != nil {
		h.E(c, h.EBindBodyError.WithPayload(
			fmt.Sprintf("failed to parse request body for helm release, error %v", err)))
		return
	}
	if release.Chart == "" {
		appErr := h.EParamsError.WithPayload("empty chart")
		h.E(c, appErr)
		return
	}
	manifest, err := renderTemplate(c, &release)
	if err != nil {
		h.E(c, err)
		return
	}
	c.Data(201, gin.MIMEYAML, []byte(manifest))
}

func parseRelease(c *gin.Context, r *api.HelmRelease) error {
	if c == nil {
		return fmt.Errorf("gin context is nil")
	}
	if r == nil {
		return fmt.Errorf("target HelmRelease is nil")
	}
	cluster := c.Param(ParamCluster)
	namespace := c.Param(ParamNamespace)
	releaseName := c.Param(ParamName)
	err := c.ShouldBindJSON(r)
	if err != nil {
		h.E(c, h.EBindBodyError.WithPayload(
			fmt.Sprintf("failed to parse request body for helm release, error %v", err)))
		return fmt.Errorf("failed to parse JSON: %v", err)
	}
	// fill in cluster, namespace, and name if not presented in body.
	if r.Cluster == "" {
		r.Cluster = cluster
	}
	if r.Namespace == "" {
		r.Namespace = namespace
	}
	if r.Name == "" {
		r.Name = releaseName
	}
	if cluster != r.Cluster {
		h.E(c, h.EParamsError.WithPayload(
			fmt.Sprintf("invalid cluster name, %s in path, but %s in body", cluster, r.Cluster)))
		return fmt.Errorf("cluster name not consistent")
	}
	if namespace != r.Namespace {
		h.E(c, h.EParamsError.WithPayload(fmt.Sprintf(
			"invalid namespace, %s in path, but %s in body", namespace, r.Namespace)))
		return fmt.Errorf("namespace not consistent")
	}
	if releaseName != "" && releaseName != r.Name {
		h.E(c, h.EParamsError.WithPayload(
			fmt.Sprintf("invalid release name, %s in path, but %s in body",
				releaseName, r.Name)))
		return fmt.Errorf("release name not consistent")
	}
	return nil
}
