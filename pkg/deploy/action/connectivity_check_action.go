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

package action

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type ConnectivityCheckItem struct {
	Protocol consts.Protocol
	Port     uint16
}

type ConnectivityCheckActionConfig struct {
	SourceNode             *pb.Node
	DestinationNode        *pb.Node
	ConnectivityCheckItems []ConnectivityCheckItem
	LogFileBasePath        string
}

type connectivityCheckAction struct {
	base
	sourceNode      *pb.Node
	destinationNode *pb.Node
	checkItems      []ConnectivityCheckItem
}

func NewConnectivityCheckAction(cfg *ConnectivityCheckActionConfig) (Action, error) {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	if cfg == nil {
		err = fmt.Errorf("action config is nil")
		return nil, err
	}
	if cfg.SourceNode == nil {
		err = fmt.Errorf("source node in config is nil")
		return nil, err
	}
	if cfg.DestinationNode == nil {
		err = fmt.Errorf("destination node in config is nil")
		return nil, err
	}
	actionName := "connectivity-" + cfg.SourceNode.Name + "-" + cfg.DestinationNode.Name
	return &connectivityCheckAction{
		base: base{
			name:              actionName,
			actionType:        ActionTypeConnectivityCheck,
			status:            ActionPending,
			logFilePath:       GenActionLogFilePath(cfg.LogFileBasePath, actionName),
			creationTimestamp: time.Now(),
		},
		sourceNode:      cfg.SourceNode,
		destinationNode: cfg.DestinationNode,
		checkItems:      cfg.ConnectivityCheckItems,
	}, nil
}

type connectivityCheckExecutor struct{}

func (e *connectivityCheckExecutor) Execute(act Action) error {
	connectivityCheckAction, ok := act.(*connectivityCheckAction)
	if !ok {
		return fmt.Errorf("action type not match: should be connectivity check action, but is %T", act)
	}
	connectivityCheckAction.status = ActionDoing

	dstNode := connectivityCheckAction.destinationNode
	srcNode := connectivityCheckAction.sourceNode
	// start SSH connection to destination node to dump packets
	sshClientDst, err := ssh.Dial("tcp",
		fmt.Sprintf("%s:%d", dstNode.Ip, dstNode.Ssh.Port), nil)
	if err != nil {
		connectivityCheckAction.status = ActionFailed
		connectivityCheckAction.err = &pb.Error{
			Reason: "failed to start SSH client",
			Detail: fmt.Sprintf("Failed to create SSH connetion to %s by connecting to %s:%d, error %v",
				dstNode.Name, dstNode.Ip, dstNode.Ssh.Port, err),
			FixMethods: "configure no-password ssh login from deploy node",
		}
		return fmt.Errorf("SSH: failed to connect to %s, error %v", dstNode.Name, err)
	}

	// start SSH connection to source node to send packets
	sshClientSrc, err := ssh.Dial("tcp",
		fmt.Sprintf("%s:%d", srcNode.Ip, srcNode.Ssh.Port), nil)
	if err != nil {
		connectivityCheckAction.status = ActionFailed
		connectivityCheckAction.err = &pb.Error{
			Reason: "failed to start SSH client",
			Detail: fmt.Sprintf("Failed to create SSH connetion to %s by connecting to %s:%d, error %v",
				srcNode.Name, srcNode.Ip, srcNode.Ssh.Port, err),
			FixMethods: "configure no-password ssh login from deploy node",
		}
		return fmt.Errorf("SSH: failed to connect to %s, error %v", srcNode.Name, err)
	}

	for _, checkItem := range connectivityCheckAction.checkItems {
		randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
		srcPort := (randGen.Uint32() % 16384) + 45000
		sshSessionDst, _ := sshClientDst.NewSession()
		sshSessionSrc, _ := sshClientSrc.NewSession()

		captureCommand := []string{"timeout", "5",
			"tcpdump", "-nni", "any", "-c", "1",
			"src", srcNode.Ip, "and", "dst", dstNode.Ip,
		}
		sendCommand := []string{"nc", "-p", fmt.Sprintf("%d", srcPort),
			"-s", srcNode.Ip}
		switch checkItem.Protocol {
		case consts.ProtocolTCP:
			captureCommand = append(captureCommand, "and", "tcp",
				"dst", "port", "dst", "port", fmt.Sprintf("%d", checkItem.Port),
				"and", "src", "port", fmt.Sprintf("%d", srcPort))
			sendCommand = append(sendCommand, "-zv",
				dstNode.Ip, fmt.Sprintf("%d", checkItem.Port))
		case consts.ProtocolUDP:
			captureCommand = append(captureCommand, "and", "udp",
				"dst", "port", "dst", "port", fmt.Sprintf("%d", checkItem.Port),
				"and", "src", "port", fmt.Sprintf("%d", srcPort))
			sendCommand = append(sendCommand, "-zuv",
				dstNode.Ip, fmt.Sprintf("%d", checkItem.Port))
		default:
			connectivityCheckAction.status = ActionFailed
			connectivityCheckAction.err = &pb.Error{
				Reason: "protocol not supported",
				Detail: fmt.Sprintf("protocol %s is not supported. supported protocols are: TCP, UDP",
					string(checkItem.Protocol)),
				FixMethods: "Use a supported protocol",
			}
			return fmt.Errorf("unssported protocol: %s", string(checkItem.Protocol))
		}

		// first, start the capturing on destination node.
		sshSessionDst.Start(strings.Join(captureCommand, " "))
		captureChan := make(chan error, 1)
		go func(errCh chan error) {
			errCh <- sshSessionDst.Wait()
		}(captureChan)
		time.Sleep(time.Second)
		sshSessionSrc.Start(strings.Join(sendCommand, " "))
		err := <-captureChan
		if err != nil {
			connectivityCheckAction.status = ActionFailed
			connectivityCheckAction.err = &pb.Error{
				Reason: "check connectivity failed",
				Detail: fmt.Sprintf("%s cannot connect to %s %s:%d",
					srcNode.Name, string(checkItem.Protocol), dstNode.Name, checkItem.Port),
				FixMethods: "configure network or firewall to allow these packets",
			}
			return fmt.Errorf("check connectivity failed: %s -> %s", srcNode.Name, dstNode.Name)
		}
	}
	connectivityCheckAction.status = ActionDone
	return nil
}
