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
	"bytes"
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
	shellCmd       *command.ShellCommand
	NodeInitAction *operation.NodeInitAction
}

func (itOps *InitKubeToolOperation) RunCommands(node *pb.Node, initAction *operation.NodeInitAction, logChan chan<- *bytes.Buffer) (stdOut, stdErr []byte, err error) {

	var imageRepository string
	var clusterDNSIP string
	var nodeIp string

	pkgMirrorUrl := fmt.Sprintf("--pkg-mirror %v", constant.DefaultPkgMirror)
	kubernetesVersion := fmt.Sprintf("--version %v", constant.DefaultKubeVersion)

	// we would use initAction's service subnet in the future
	clusterDNSIP = fmt.Sprintf("--cluster-dns %v", getDNSIP(constant.DefaultServiceSubnet))

	// we would use initAction's image repository in the future
	imageRepository = fmt.Sprintf("--image-repository %v", constant.DefaultImageRepository)

	m, err := machine.NewMachine(node)
	if err != nil {
		return nil, nil, err
	}

	logBuffer := &bytes.Buffer{}

	itOps.NodeInitAction = initAction

	// close ssh client if machine is not nil
	if m != nil {
		defer m.Close()
	}

	// copy init_deploy_kubetool.sh to target machine
	scriptFile, err := assets.Assets.Open(consts.DefaultKubeToolScript)
	if err != nil {
		return nil, nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, operation.InitRemoteScriptPath+consts.DefaultKubeToolScript); err != nil {
		return nil, nil, err
	}

	// copy commmon lib.sh to target machine
	scriptFile, err = assets.Assets.Open(DefaultCommonLibPath)
	if err != nil {
		return nil, nil, err
	}
	defer scriptFile.Close()

	if err := m.PutFile(scriptFile, operation.InitRemoteScriptPath+DefaultCommonLibPath); err != nil {
		return nil, nil, err
	}

	// setup repos
	itOps.shellCmd = command.NewShellCommand(m, "bash", fmt.Sprintf("%v setup repos %v", operation.InitRemoteScriptPath+consts.DefaultKubeToolScript,
		pkgMirrorUrl)).
		WithDescription("初始化 kubernetes repos 环境").
		WithExecuteLogWriter(logBuffer)

	// install kubelet, kubeadm, kubectl
	itOps.shellCmd = command.NewShellCommand(m, "bash", fmt.Sprintf("%v setup kubelet %v %v %v %v", operation.InitRemoteScriptPath+consts.DefaultKubeToolScript,
		kubernetesVersion, imageRepository, clusterDNSIP, nodeIp)).
		WithDescription("初始化安装 kubernetes 工具").
		WithExecuteLogWriter(logBuffer)

	// run commands
	stdOut, stdErr, err = itOps.shellCmd.Execute()

	// write to log channel
	logChan <- logBuffer

	return
}

// get dns IP from subnet
func getDNSIP(serviceSubnet string) string {
	dnsIP, err := parseServiceSubnet(serviceSubnet)
	if err == nil {
		return fmt.Sprintf("%v", dnsIP.String())
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
