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

package command

import (
	"fmt"

	clientset "k8s.io/client-go/kubernetes"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
)

// KubectlCommand is a command execute by kubectl
type KubectlCommand struct {
	kubeClient clientset.Interface
	kubeConfig string
	namespace  string

	*ShellCommand
}

func NewKubectlCommand(machine *machine.Machine, kubeConfigPath string, ns string, subCommands ...string) Command {

	if ns != "" {
		subCommands = append(subCommands, fmt.Sprintf("-n %s", ns))
	}
	if kubeConfigPath != "" {
		subCommands = append(subCommands, fmt.Sprintf("--kubeconfig %s", kubeConfigPath))
	}

	c := &KubectlCommand{
		kubeConfig: kubeConfigPath,
		namespace:  ns,
		ShellCommand: NewShellCommand(
			machine,
			"/usr/bin/kubectl",
			subCommands...,
		),
	}

	return c
}
