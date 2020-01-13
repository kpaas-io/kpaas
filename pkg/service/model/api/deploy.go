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

package api

import (
	"github.com/kpaas-io/kpaas/pkg/constant"
)

type (
	GetDeploymentReportResponse struct {
		DeployItems         []DeploymentResponseData `json:"deployItems"`
		DeployClusterStatus DeployClusterStatus      `json:"deployClusterStatus" enums:"pending,running,successful,failed,workedButHaveError"` // The cluster deployment status
		DeployClusterError  *Error                   `json:"deployClusterError,omitempty"`                                                     // Deploy cluster error message
	}

	DeploymentResponseData struct {
		DeployItem constant.DeployItem `json:"DeployItem" enums:"master,worker,etcd,ingress,network"`
		Nodes      []DeploymentNode    `json:"nodes"`
	}

	DeploymentNode struct {
		Name   string       `json:"name"`                                                     // node name
		Status DeployStatus `json:"result" enums:"pending,running,successful,failed,aborted"` // Checking Result
		Error  *Error       `json:"error,omitempty"`
	}

	DeployStatus        string
	DeployClusterStatus string
)

const (
	DeployStatusPending    DeployStatus = "pending"
	DeployStatusRunning    DeployStatus = "running"
	DeployStatusSuccessful DeployStatus = "successful"
	DeployStatusFailed     DeployStatus = "failed"
	DeployStatusAborted    DeployStatus = "aborted"

	DeployClusterStatusPending            DeployClusterStatus = "pending"
	DeployClusterStatusRunning            DeployClusterStatus = "running"
	DeployClusterStatusSuccessful         DeployClusterStatus = "successful"
	DeployClusterStatusFailed             DeployClusterStatus = "failed"
	DeployClusterStatusWorkedButHaveError DeployClusterStatus = "workedButHaveError"
)
