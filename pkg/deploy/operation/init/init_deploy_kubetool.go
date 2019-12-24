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

package init

import (
	"fmt"
	"net"

	"k8s.io/kubernetes/pkg/registry/core/service/ipallocator"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/deploy/assets"
	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type InitKubeToolOperation struct {
	operation.BaseOperation
	InitOperations
	Machine        *machine.Machine
	NodeInitAction *operation.NodeInitAction
}

func (itOps *InitKubeToolOperation) getScript() string {
	itOps.Script = consts.DefaultKubeToolScript
	return itOps.Script
}

func (itOps *InitKubeToolOperation) getScriptPath() string {
	itOps.ScriptPath = RemoteScriptPath
	return itOps.ScriptPath
}

func (itOps *InitKubeToolOperation) GetOperations(node *pb.Node, initAction *operation.NodeInitAction) (operation.Operation, error) {

	var imageRepository string
	var clusterDNSIP string
	var nodeIP string

	pkgMirrorUrl := fmt.Sprintf("--pkg-mirror %v", constant.DefaultPkgMirror)
	kubernetesVersion := fmt.Sprintf("--version %v", constant.DefaultKubeVersion)

	// we would use initAction's service subnet in the future
	clusterDNSIP = fmt.Sprintf("--cluster-dns %v", getDNSIP(constant.DefaultServiceSubnet))

	if initAction.NodeInitConfig.Node.Ip == "" {
		return nil, fmt.Errorf("current node %v ip can not be empty", initAction.NodeInitConfig.Node.Name)
	}

	nodeIP = fmt.Sprintf("--node-ip %v", initAction.NodeInitConfig.Node.Ip)

	// we would use initAction's image repository in the future
	imageRepository = fmt.Sprintf("--image-repository %v", constant.DefaultImageRepository)

	ops := &InitKubeToolOperation{}
	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, err
	}
	itOps.Machine = m
	itOps.NodeInitAction = initAction

	scriptFile, err := assets.Assets.Open(itOps.getScript())
	if err != nil {
		return nil, err
	}

	if err := m.PutFile(scriptFile, itOps.getScriptPath()+itOps.getScript()); err != nil {
		return nil, err
	}

	// setup repos
	ops.AddCommands(command.NewShellCommand(m, "bash", fmt.Sprintf("%v setup repos %v", itOps.getScriptPath()+itOps.getScript(),
		pkgMirrorUrl)))

	// install kubelet, kubeadm, kubectl
	ops.AddCommands(command.NewShellCommand(m, "bash", fmt.Sprintf("%v setup kubelet %v %v %v %v", itOps.getScriptPath()+itOps.getScript(),
		kubernetesVersion, nodeIP, imageRepository, clusterDNSIP)))

	ops.AddCommands(command.NewShellCommand(m, "bash", itOps.getScriptPath()+itOps.getScript()))
	return ops, nil
}

func (itOps *InitKubeToolOperation) CloseSSH() {
	itOps.Machine.Close()
}

// get dns IP from subnet
func getDNSIP(serviceSubnet string) string {
	dnsIP, err := parseServiceSubnet(serviceSubnet)
	if err == nil {
		return fmt.Sprintf("--cluster-dns %v", dnsIP.String())
	}
	return ""
}

// parse dns ip from service subnet
func parseServiceSubnet(serviceSubnet string) (net.IP, error) {
	// Get the service subnet CIDR
	_, svcSubnetCIDR, err := net.ParseCIDR(serviceSubnet)
	if err != nil {
		return nil, fmt.Errorf("%v, couldn't parse service subnet CIDR %q", err, serviceSubnet)
	}

	// Selects the 10th IP in service subnet CIDR range as dnsIP
	dnsIP, err := ipallocator.GetIndexedIP(svcSubnetCIDR, 10)
	if err != nil {
		return nil, fmt.Errorf("%v, unable to get tenth IP address from service subnet CIDR %s", err, svcSubnetCIDR.String())
	}

	return dnsIP, nil
}
