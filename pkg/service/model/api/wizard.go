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
	GetWizardResponse struct {
		ClusterData         Cluster                      `json:"cluster"`                                                                          // Cluster Information
		NetworkOptions      NetworkOptions               `json:"networkOptions"`                                                                   // Network options
		NodesData           []NodeData                   `json:"nodes"`                                                                            // Nodes Information
		CheckingData        []CheckingResultResponseData `json:"checkingData"`                                                                     // Check result
		DeploymentData      []DeploymentResponseData     `json:"deploymentData"`                                                                   // Deployment result
		CheckResult         constant.CheckResult         `json:"checkResult" enums:"pending,running,successful,failed"`                            // Nodes check result
		DeployClusterStatus DeployClusterStatus          `json:"deployClusterStatus" enums:"pending,running,successful,failed,workedButHaveError"` // Cluster deployment status
	}
)
