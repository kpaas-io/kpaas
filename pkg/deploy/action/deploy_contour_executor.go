// Copyright 2020 Shanghai JingDuo Information Technology co., Ltd.
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

package action

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	deployMachine "github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/contour"
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const installContourFilePath = "/tmp/installContour.yaml"
const contourYAML = `apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  labels:
    app: contour
  name: contour
  namespace: kube-system
spec:
  selector:
    matchLabels:
      app: contour
  template:
    metadata:
      annotations:
        prometheus.io/format: prometheus
        prometheus.io/path: /stats
        prometheus.io/port: "8003"
        prometheus.io/scrape: "true"
      labels:
        app: contour
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: node-role.kubernetes.io/ingress
                operator: In
                values:
                - envoy
      containers:
      - args:
        - serve
        - --incluster
        - --envoy-http-port=80
        - --envoy-https-port=443
        command:
        - contour
        image: kpaas/contour:v0.8
        imagePullPolicy: IfNotPresent
        name: contour
        ports:
        - containerPort: 8000
          hostPort: 8000
          name: contour
          protocol: TCP
      - args:
        - -c
        - /config/contour.yaml
        - --service-cluster
        - cluster0
        - --service-node
        - node0
        - -l
        - info
        - --v2-config-only
        command:
        - envoy
        image: kpaas/envoy:v1.7.0
        imagePullPolicy: IfNotPresent
        name: envoy
        ports:
        - containerPort: 80
          hostPort: 80
          name: http
          protocol: TCP
        - containerPort: 443
          hostPort: 443
          name: https
          protocol: TCP
        volumeMounts:
        - mountPath: /config
          name: contour-config
      - image: kpaas/statsd-exporter:v0.1
        imagePullPolicy: IfNotPresent
        name: statsd-prom-bridge
        ports:
        - containerPort: 9102
          hostPort: 9102
          protocol: TCP
        - containerPort: 9125
          hostPort: 9125
          protocol: UDP
      dnsPolicy: ClusterFirst
      hostNetwork: true
      initContainers:
      - args:
        - bootstrap
        - /config/contour.yaml
        - --statsd-enabled
        - --statsd-address=127.0.0.1
        - --admin-address=0.0.0.0
        command:
        - contour
        image: kpaas/contour:v0.8
        imagePullPolicy: IfNotPresent
        name: envoy-initconfig
        volumeMounts:
        - mountPath: /config
          name: contour-config
      restartPolicy: Always
      schedulerName: default-scheduler
      serviceAccount: contour
      serviceAccountName: contour
      tolerations:
      - operator: Exists
      volumes:
      - emptyDir: {}
        name: contour-config
  updateStrategy:
    type: OnDelete
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: contour
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: contour
subjects:
- kind: ServiceAccount
  name: contour
  namespace: heptio-contour
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: contour
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - endpoints
  - nodes
  - pods
  - secrets
  verbs:
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - get
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - extensions
  resources:
  - ingresses
  verbs:
  - get
  - list
  - watch
- apiGroups: ["contour.heptio.com"]
  resources: ["ingressroutes"]
  verbs:
  - get
  - list
  - watch
  - put
  - post
  - patch
---
apiVersion: v1
kind: Namespace
metadata:
  name: heptio-contour
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: contour
  namespace: heptio-contour
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: ingressroutes.contour.heptio.com
  labels:
    component: ingressroute
spec:
  group: contour.heptio.com
  version: v1beta1
  scope: Namespaced
  names:
    plural: ingressroutes
    kind: IngressRoute
  additionalPrinterColumns:
    - name: FQDN
      type: string
      description: Fully qualified domain name
      JSONPath: .spec.virtualhost.fqdn
    - name: TLS Secret
      type: string
      description: Secret with TLS credentials
      JSONPath: .spec.virtualhost.tls.secretName
    - name: First route
      type: string
      description: First routes defined
      JSONPath: .spec.routes[0].match
    - name: Status
      type: string
      description: The current status of the IngressRoute
      JSONPath: .status.currentStatus
    - name: Status Description
      type: string
      description: Description of the current status
      JSONPath: .status.description
  validation:
    openAPIV3Schema:
      properties:
        spec:
          properties:
            virtualhost:
              properties:
                fqdn:
                  type: string
                  pattern: ^([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-z]{2,}$
                tls:
                  properties:
                    secretName:
                      type: string
                      pattern: ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ # DNS-1123 subdomain
                    minimumProtocolVersion:
                      type: string
                      enum:
                        - "1.3"
                        - "1.2"
                        - "1.1"
            strategy:
              type: string
              enum:
                - RoundRobin
                - WeightedLeastRequest
                - Random
                - RingHash
                - Maglev
            healthCheck:
              type: object
              required:
                - path
              properties:
                path:
                  type: string
                  pattern: ^\/.*$
                intervalSeconds:
                  type: integer
                timeoutSeconds:
                  type: integer
                unhealthyThresholdCount:
                  type: integer
                healthyThresholdCount:
                  type: integer
            routes:
              type: array
              items:
                required:
                  - match
                properties:
                  match:
                    type: string
                    pattern: ^\/.*$
                  delegate:
                    type: object
                    required:
                      - name
                    properties:
                      name:
                        type: string
                        pattern: ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ # DNS-1123 subdomain
                      namespace:
                        type: string
                        pattern: ^[a-z0-9]([-a-z0-9]*[a-z0-9])?$ # DNS-1123 label
                  services:
                    type: array
                    items:
                      type: object
                      required:
                        - name
                        - port
                      properties:
                        name:
                          type: string
                          pattern: ^[a-z]([-a-z0-9]*[a-z0-9])?$ # DNS-1035 label
                        port:
                          type: integer
                        weight:
                          type: integer
                        strategy:
                          type: string
                          enum:
                            - RoundRobin
                            - WeightedLeastRequest
                            - Random
                            - RingHash
                            - Maglev
                        healthCheck:
                          type: object
                          required:
                            - path
                          properties:
                            path:
                              type: string
                              pattern: ^\/.*$
                            intervalSeconds:
                              type: integer
                            timeoutSeconds:
                              type: integer
                            unhealthyThresholdCount:
                              type: integer
                            healthyThresholdCount:
                              type: integer
`

func init() {
	RegisterExecutor(ActionTypeDeployContour, new(deployContourExecutor))
}

type deployContourExecutor struct {
	logger           *logrus.Entry
	masterMachine    deployMachine.IMachine
	action           *DeployContourAction
	executeLogWriter io.Writer
}

func (executor *deployContourExecutor) Execute(act Action) *protos.Error {

	action, ok := act.(*DeployContourAction)
	if !ok {
		return errOfTypeMismatched(new(DeployContourAction), act)
	}

	executor.action = action

	executor.initLogger()
	executor.initExecuteLogWriter()

	executor.logger.Info("start to execute deploy contour executor")

	if err := executor.connectMasterNode(); err != nil {
		return err
	}
	defer executor.disconnectMasterNode()

	operations := []func() *protos.Error{
		executor.writeYAML,
		executor.applyYAML,
	}

	for _, operation := range operations {
		err := operation()
		if err != nil {
			return err
		}
	}

	executor.logger.Info("deploy contour finished")

	return nil
}

func (executor *deployContourExecutor) initLogger() {
	executor.logger = logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: executor.action.GetName(),
		"clusterName":         executor.action.config.ClusterConfig.GetClusterName(),
	})
}

func (executor *deployContourExecutor) initExecuteLogWriter() {

	executor.executeLogWriter = executor.action.GetExecuteLogBuffer()
}

func (executor *deployContourExecutor) connectMasterNode() *protos.Error {
	var err error
	executor.masterMachine, err = deployMachine.NewMachine(executor.action.config.MasterNodes[0])
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("failed to connect master node")
		return &protos.Error{
			Reason:     "connecting failed",
			Detail:     fmt.Sprintf("failed to connect master node, err: %s", err),
			FixMethods: "please check deploy node config to ensure master node can be connected successfully",
		}
	}
	return nil
}

func (executor *deployContourExecutor) disconnectMasterNode() {
	if executor.masterMachine != nil {
		executor.masterMachine.Close()
	}
}

func (executor *deployContourExecutor) writeYAML() *protos.Error {

	executor.logger.Debug("Start to write contour yaml")

	operation := contour.NewWriteFile(
		&contour.WriteFileConfig{
			Node:             executor.masterMachine,
			Logger:           executor.logger,
			ExecuteLogWriter: executor.executeLogWriter,
			FilePath:         installContourFilePath,
			FileContent:      contourYAML,
		},
	)

	if err := operation.Execute(); err != nil {
		executor.logger.WithField("error", err).Error("write contour yaml error")
		return err
	}

	executor.logger.Info("Finish to write contour yaml action")
	return nil
}

func (executor *deployContourExecutor) applyYAML() *protos.Error {

	executor.logger.Debug("Start to apply yaml")

	operation := contour.NewApplyYAML(
		&contour.ApplyYAMLConfig{
			Node:             executor.masterMachine,
			Logger:           executor.logger,
			ExecuteLogWriter: executor.executeLogWriter,
			FilePath:         installContourFilePath,
		},
	)

	if err := operation.Execute(); err != nil {
		executor.logger.WithField("error", err).Error("apply yaml error")
		return err
	}

	executor.logger.Info("Finish to apply yaml action")
	return nil
}
