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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/kpaas-io/kpaas/pkg/utils/validator"
)

type (
	Cluster struct {
		ShortName                string                   `json:"shortName" binding:"required" minLength:"1" maxLength:"20"`
		Name                     string                   `json:"name" binding:"required"`
		KubeAPIServerConnectType KubeAPIServerConnectType `json:"kubeAPIServerConnectType" binding:"required" enums:"firstMasterIP,keepalived,loadbalancer"` // kube-apiserver connect type
		VIP                      string                   `json:"vip,omitempty" maxLength:"15"`                                                              // keepalived listen virtual ip
		NetInterfaceName         string                   `json:"netInterfaceName,omitempty" maxLength:"30"`                                                 // keepalived listen net interface name
		LoadbalancerIP           string                   `json:"loadbalancerIP,omitempty" maxLength:"15"`                                                   // kube-apiserver loadbalancer ip when kubeAPIServerConnectType is loadbalancer required
		LoadbalancerPort         uint16                   `json:"loadbalancerPort,omitempty" minimum:"1" maximum:"65535"`                                    // kube-apiserver loadbalancer port when kubeAPIServerConnectType is loadbalancer required
		NodePortMinimum          uint16                   `json:"nodePortMinimum" minimum:"1" default:"30000"`
		NodePortMaximum          uint16                   `json:"nodePortMaximum" maximum:"65535" default:"32767"`
		Labels                   []Label                  `json:"labels"`
		Annotations              []Annotation             `json:"annotations"`
	}

	KubeAPIServerConnectType string

	Label struct {
		Key   string `json:"key" binding:"required" minimum:"1" maximum:"253"`
		Value string `json:"value" binding:"required" minimum:"1"`
	}

	Annotation struct {
		Key   string `json:"key" binding:"required" minimum:"1" maximum:"253"`
		Value string `json:"value" binding:"required" minimum:"1"`
	}
)

const (
	KubeAPIServerConnectTypeFirstMasterIP KubeAPIServerConnectType = "firstMasterIP"
	KubeAPIServerConnectTypeKeepalived    KubeAPIServerConnectType = "keepalived"
	KubeAPIServerConnectTypeLoadBalancer  KubeAPIServerConnectType = "loadbalancer"

	ClusterNameLengthLimit          = 30
	ClusterShortNameLengthLimit     = 20
	ClusterIPLengthLimit            = 15
	ClusterNetInterfaceLengthLimit  = 30
	ClusterLoadbalancerPortMinimum  = 1
	ClusterLoadbalancerPortMaximum  = 65535
	ClusterNodePortMinimum          = 1
	ClusterNodePortMaximum          = 65535
	LabelKeyLengthLimit             = 253
	LabelKeySegmentLengthLimit      = 63
	AnnotationKeyLengthLimit        = 253
	AnnotationKeySegmentLengthLimit = 63
)

func (cluster *Cluster) Validate() error {

	wrapper := validator.NewWrapper(
		validator.ValidateString(cluster.Name, "name", validator.ItemNotEmptyLimit, ClusterNameLengthLimit),
		validator.ValidateString(cluster.ShortName, "shortName", validator.ItemNotEmptyLimit, ClusterShortNameLengthLimit),
		validator.ValidateStringOptions(string(cluster.KubeAPIServerConnectType),
			"kubeAPIServerConnectType",
			[]string{string(KubeAPIServerConnectTypeFirstMasterIP), string(KubeAPIServerConnectTypeKeepalived), string(KubeAPIServerConnectTypeLoadBalancer)}),
		validator.ValidateIntRange(int(cluster.NodePortMinimum), "nodePortMinimum", ClusterNodePortMinimum, ClusterNodePortMaximum),
		validator.ValidateIntRange(int(cluster.NodePortMaximum), "nodePortMaximum", ClusterNodePortMinimum, ClusterNodePortMaximum),
		func() error {
			if cluster.NodePortMinimum > cluster.NodePortMaximum {
				return fmt.Errorf("nodePortMinimum must be not larger than nodePortMaximum")
			}
			return nil
		},
	)

	switch cluster.KubeAPIServerConnectType {
	case KubeAPIServerConnectTypeKeepalived:
		wrapper.AddValidateFunc(
			validator.ValidateString(cluster.VIP, "vip", validator.ItemNotEmptyLimit, ClusterIPLengthLimit),
			validator.ValidateIP(cluster.VIP, "vip"),
			validator.ValidateString(cluster.NetInterfaceName, "netInterfaceName", validator.ItemNotEmptyLimit, ClusterNetInterfaceLengthLimit),
		)
	case KubeAPIServerConnectTypeLoadBalancer:
		wrapper.AddValidateFunc(
			validator.ValidateString(cluster.LoadbalancerIP, "loadbalancerIP", validator.ItemNotEmptyLimit, ClusterIPLengthLimit),
			validator.ValidateIP(cluster.LoadbalancerIP, "loadbalancerIP"),
			validator.ValidateIntRange(int(cluster.LoadbalancerPort), "loadbalancerPort", ClusterLoadbalancerPortMinimum, ClusterLoadbalancerPortMaximum),
		)
	}

	for _, label := range cluster.Labels {

		wrapper.AddValidateFunc(
			func() error {
				return label.Validate()
			},
		)
	}

	for _, annotation := range cluster.Annotations {

		wrapper.AddValidateFunc(
			func() error {
				return annotation.Validate()
			},
		)
	}

	return wrapper.Validate()
}

func (label *Label) Validate() error {

	return validator.NewWrapper(
		validator.ValidateString(label.Key, "key", validator.ItemNotEmptyLimit, LabelKeyLengthLimit),
		ValidateStringFunctionReturnErrorMessages(validation.IsQualifiedName, label.Key, "label.key"),
		validator.ValidateString(label.Value, "value", validator.ItemNotEmptyLimit, validator.ItemNoLimit),
		ValidateStringFunctionReturnErrorMessages(validation.IsValidLabelValue, label.Value, "label.value"),
	).Validate()
}

func (annotation *Annotation) Validate() error {

	return validator.NewWrapper(
		validator.ValidateString(annotation.Key, "key", validator.ItemNotEmptyLimit, AnnotationKeyLengthLimit),
		ValidateStringFunctionReturnErrorMessages(validation.IsQualifiedName, annotation.Key, "annotation.key"),
		validator.ValidateString(annotation.Value, "value", validator.ItemNotEmptyLimit, validator.ItemNoLimit),
	).Validate()
}

func ValidateStringFunctionReturnErrorMessages(function func(string) []string, value, keyName string) validator.ValidateFunc {

	return func() error {
		errorMessages := validation.IsQualifiedName(value)

		if len(errorMessages) > 0 {
			return fmt.Errorf("%s %s", keyName, strings.Join(errorMessages, ", "))
		}

		return nil
	}
}
