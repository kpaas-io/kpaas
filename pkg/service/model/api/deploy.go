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

type (
	GetDeploymentReportResponse struct {
		Roles               []DeploymentResponseData `json:"roles"`
		DeployClusterStatus DeployClusterStatus      `json:"deployClusterStatus" enums:"notRunning,running,successful,failed,workedButHaveError"` // The cluster deployment status
	}

	DeploymentResponseData struct {
		Role  MachineRole      `json:"role" enums:"master,worker,etcd"`
		Nodes []DeploymentNode `json:"nodes"`
	}

	DeploymentNode struct {
		Name   string       `json:"name"`                                                      // node name
		Status DeployStatus `json:"result" enums:"pending,deploying,completed,failed,aborted"` // Checking Result
		Error  *Error       `json:"error,omitempty"`
	}

	DeployStatus        string
	DeployClusterStatus string
)

const (
	DeployStatusPending   DeployStatus = "pending"
	DeployStatusDeploying DeployStatus = "deploying"
	DeployStatusCompleted DeployStatus = "completed"
	DeployStatusFailed    DeployStatus = "failed"
	DeployStatusAborted   DeployStatus = "aborted"

	DeployClusterStatusNotRunning         DeployClusterStatus = "notRunning"
	DeployClusterStatusRunning            DeployClusterStatus = "running"
	DeployClusterStatusSuccessful         DeployClusterStatus = "successful"
	DeployClusterStatusFailed             DeployClusterStatus = "failed"
	DeployClusterStatusWorkedButHaveError DeployClusterStatus = "workedButHaveError"
)
