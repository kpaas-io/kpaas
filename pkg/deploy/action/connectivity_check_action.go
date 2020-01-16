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
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const ActionTypeConnectivityCheck Type = "ConnectivityCheck"

const (
	// define constants for messages and formats for future i18n
	// 发送数据包失败
	reasonFailedToSendPacket = "failed to send packet"
	// 节点%s 发送数据包到节点 %s(%s:%d)，错误%v
	detailFailedToSendPacketFormat = "%s failed to send packet to %s (%s:%d), error %v"
	// 查看命令%s在节点%s上是否存在并有权限运行
	fixFailedToSendPacketFormat = "check existance and permission to run command %s on node %s"
)

// ConnectivityCheckItem an item representing one check item of checking whether a node can connect to another by the protocol and port.
type ConnectivityCheckItem struct {
	Protocol    consts.Protocol
	Port        uint16
	Name        string
	Description string
	Status      ItemStatus
	Err         *pb.Error
}

// ConnectivityCheckActionConfig configuration of checking connectivity from soruce to destination.
type ConnectivityCheckActionConfig struct {
	SourceNode             *pb.Node
	DestinationNode        *pb.Node
	ConnectivityCheckItems []*ConnectivityCheckItem
	LogFileBasePath        string
}

type ConnectivityCheckAction struct {
	Base

	SourceNode      *pb.Node
	DestinationNode *pb.Node
	CheckItems      []*ConnectivityCheckItem
}

// NewConnectivityCheckAction creates an action to check connectivity from soruce to destination.
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
	actionName := GenActionName(ActionTypeConnectivityCheck)
	return &ConnectivityCheckAction{
		Base: Base{
			Name:              actionName,
			ActionType:        ActionTypeConnectivityCheck,
			Status:            ActionPending,
			LogFilePath:       GenActionLogFilePath(cfg.LogFileBasePath, actionName, cfg.SourceNode.GetName()),
			CreationTimestamp: time.Now(),
			Node:              cfg.SourceNode,
		},
		SourceNode:      cfg.SourceNode,
		DestinationNode: cfg.DestinationNode,
		CheckItems:      cfg.ConnectivityCheckItems,
	}, nil
}

func init() {
	RegisterExecutor(ActionTypeConnectivityCheck, new(connectivityCheckExecutor))
}

type connectivityCheckExecutor struct{}

func (e *connectivityCheckExecutor) Execute(act Action) *pb.Error {
	connectivityCheckAction, ok := act.(*ConnectivityCheckAction)
	if !ok {
		return errOfTypeMismatched(new(ConnectivityCheckAction), act)
	}

	logger := logrus.WithFields(logrus.Fields{
		consts.LogFieldAction: act.GetName(),
	})

	dstNode := connectivityCheckAction.DestinationNode
	srcNode := connectivityCheckAction.SourceNode
	logger.Infof("check network connectiviy from %s to %s", srcNode.Name, dstNode.Name)
	// make a executor client for destination node to capture packets
	dstMachine, err := machine.NewMachine(dstNode)
	if err != nil {
		return &pb.Error{
			Reason: "failed to start SSH client",
			Detail: fmt.Sprintf("Failed to create connetion to %s by connecting to %s, error %v",
				dstNode.Name, dstNode.Ip, err),
			FixMethods: "configure no-password ssh login from deploy node",
		}
	}
	defer dstMachine.Close()

	// make a executor client for source node to send packets
	srcMachine, err := machine.NewMachine(srcNode)
	if err != nil {
		return &pb.Error{
			Reason: "failed to start SSH client",
			Detail: fmt.Sprintf("Failed to create connetion to %s by connecting to %s, error %v",
				srcNode.Name, srcNode.Ip, err),
			FixMethods: "configure no-password ssh login from deploy node",
		}
	}
	defer srcMachine.Close()

	for _, checkItem := range connectivityCheckAction.CheckItems {
		captureCommand := []string{"timeout", "5",
			"tcpdump", "-nni", "any", "-c", "1",
			"src", srcNode.Ip, "and", "dst", dstNode.Ip,
		}
		sendCommand := []string{}
		switch checkItem.Protocol {
		case consts.ProtocolTCP:
			captureCommand = append(captureCommand, "and", "tcp",
				"dst", "port", fmt.Sprintf("%d", checkItem.Port))
			sendCommand = append(sendCommand, "echo", "'test'", ">",
				fmt.Sprintf("/dev/tcp/%s/%d", dstNode.Ip, checkItem.Port))
		case consts.ProtocolUDP:
			captureCommand = append(captureCommand, "and", "udp",
				"dst", "port", fmt.Sprintf("%d", checkItem.Port))
			sendCommand = append(sendCommand, "echo", "'test'", ">",
				fmt.Sprintf("/dev/udp/%s/%d", dstNode.Ip, checkItem.Port))
		default:
			return &pb.Error{
				Reason: "protocol not supported",
				Detail: fmt.Sprintf("protocol %s is not supported. supported protocols are: TCP, UDP",
					string(checkItem.Protocol)),
				FixMethods: "Use a supported protocol",
			}
		}

		checkItem.Status = ItemDoing

		executeLogBuf := act.GetExecuteLogBuffer()

		captureChan := make(chan error)
		dstExecuteLogBuf := &bytes.Buffer{}
		go func(errCh chan error) {
			var e error
			dstCommand := command.NewShellCommand(dstMachine,
				captureCommand[0], captureCommand[1:]...).
				WithDescription("capture test packet on " + dstNode.Name).
				WithExecuteLogWriter(dstExecuteLogBuf)
			_, _, e = dstCommand.Execute()
			errCh <- e
		}(captureChan)

		// sleep one second to make sure that the packet is sent after capturing started
		time.Sleep(time.Second)
		srcExecuteLogBuf := &bytes.Buffer{}
		srcCommand := command.NewShellCommand(srcMachine, sendCommand[0], sendCommand[1:]...).
			WithDescription("send test packet").
			WithExecuteLogWriter(srcExecuteLogBuf)
		_, srcStderr, srcErr := srcCommand.Execute()
		if executeLogBuf != nil {
			io.Copy(executeLogBuf, srcExecuteLogBuf)
		}
		if srcErr != nil {
			srcStderrString := string(srcStderr)
			if strings.Contains(srcStderrString, "Connection refused") {
				// does nothing if the stderr tells that connection refused.
			} else {
				checkItem.Err = &pb.Error{
					Reason: reasonFailedToSendPacket,
					Detail: fmt.Sprintf(detailFailedToSendPacketFormat,
						srcNode.Name, string(checkItem.Protocol), dstNode.Name,
						checkItem.Port, srcStderrString),
					FixMethods: fmt.Sprintf(fixFailedToSendPacketFormat,
						sendCommand[0], srcNode.Name),
				}
				checkItem.Status = ItemFailed
				continue
			}
		}

		// wait for capture command to terminate
		dstErr := <-captureChan
		if executeLogBuf != nil {
			io.Copy(act.GetExecuteLogBuffer(), dstExecuteLogBuf)
		}
		if dstErr != nil {
			checkItem.Err = &pb.Error{
				Reason: "check connectivity failed",
				Detail: fmt.Sprintf("%s cannot connect to %s %s:%d",
					srcNode.Name, string(checkItem.Protocol), dstNode.Name, checkItem.Port),
				FixMethods: "configure network or firewall to allow these packets",
			}
			checkItem.Status = ItemFailed

			// does not return here to continue to check other items
		} else {
			checkItem.Status = ItemDone
		}
	} // end of for in range chekItems
	return nil
}
