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

package deploy

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/constant"
	"github.com/kpaas-io/kpaas/pkg/service/model/api"
	"github.com/kpaas-io/kpaas/pkg/service/model/wizard"
	"github.com/kpaas-io/kpaas/pkg/utils/h"
	"github.com/kpaas-io/kpaas/pkg/utils/log"
)

// @ID UploadBatchNodes
// @Summary Upload batch nodes configuration
// @Description Upload batch nodes configuration file to node list
// @Tags nodes
// @Accept text/plain
// @Produce application/json
// @Param nodes body string true "node list"
// @Success 201 {object} api.GetNodeListResponse
// @Router /api/v1/deploy/wizard/batchnodes [post]
func UploadBatchNodes(c *gin.Context) {

	requestData, err := getUploadBatchNodesRequestData(c)
	if err != nil {
		h.E(c, h.EParamsError.WithPayload(err.Error()))
		log.ReqEntry(c).Infof("parameter error: %v, %T", err, err)
		return
	}

	nodeList := make([]*wizard.Node, 0, len(requestData))
	for _, data := range requestData {

		node := wizard.NewNode()
		node.Name = data.Name
		node.DockerRootDirectory = data.DockerRootDirectory
		node.MachineRoles = data.MachineRoles

		node.IP = data.IP
		node.Port = data.Port
		node.Username = data.Username
		node.AuthenticationType = convertAPIAuthenticationTypeToModelAuthenticationType(data.AuthenticationType)
		switch data.AuthenticationType {
		case api.AuthenticationTypePassword:
			node.Password = data.Password
		case api.AuthenticationTypePrivateKey:
			node.PrivateKeyName = data.PrivateKeyName
		}

		nodeList = append(nodeList, node)
	}

	wizardData := wizard.GetCurrentWizard()
	err = wizardData.AddNodeList(nodeList)
	if err != nil {

		h.E(c, err)
		return
	}

	responseNodeList := getWizardNodes()
	h.R(c, api.GetNodeListResponse{
		Nodes: *responseNodeList,
	})
}

func getUploadBatchNodesRequestData(c *gin.Context) (nodeList []*api.NodeData, err error) {

	data, err := c.GetRawData()
	log.ReqEntry(c).Tracef("rawData: %v, err: %v", string(data), err)
	if err != nil {
		return
	}

	/**
	Excample Template
	#<hostname> <user>  <role,role,role>         <IP>             <ssh port>  <password>          <login key name>        <docker path>
	k8s-master1   root	    master,etcd          192.168.3.223    22          111111111111	      -		                  /var/lib/docker
	k8s-master2   root	    master,etcd          192.168.3.224    22          111111111111	      -   	                  /var/lib/docker
	k8s-master3   root	    master,etcd          192.168.3.227    22          111111111111	      -   	                  /var/lib/docker
	k8s-worker1   root	    worker,etcd          192.168.3.226    22          -	                  worker_key 	          /var/lib/docker
	k8s-worker2   root	    worker,etcd          192.168.3.229    22          -	                  worker_key 	          /var/lib/docker
	k8s-worker3   root	    worker               192.168.3.230    22          -	                  worker_key 	          /var/lib/docker
	*/
	matches, groupNames := tryToMatchBatchNodes(data)

	log.ReqEntry(c).Tracef("matches: %v", matches)

	if len(matches) <= 0 {
		err = fmt.Errorf("node list empty")
		return
	}

	nodeList = make([]*api.NodeData, 0, len(matches))

	log.ReqEntry(c).Tracef("match node count: %d", len(matches))

	ipList := make(map[string]bool)
	nameList := make(map[string]bool)

	for _, match := range matches {

		matchMap := make(map[string]string)
		for i, groupName := range groupNames {

			if i > 0 && i <= len(match) {
				matchMap[groupName] = match[i]
			}
		}

		log.ReqEntry(c).Tracef("matched roles: %s", matchMap["roles"])
		roles := splitInputRoles(matchMap)
		log.ReqEntry(c).Tracef("roles: %#v", roles)

		loginData := api.SSHLoginData{
			Username:           matchMap["username"],
			Password:           matchMap["password"],
			AuthenticationType: api.AuthenticationTypePassword,
		}

		if matchMap["password"] == "-" && matchMap["privateKeyName"] != "-" {
			loginData.AuthenticationType = api.AuthenticationTypePrivateKey
			loginData.PrivateKeyName = matchMap["privateKeyName"]
			loginData.Password = ""
		}

		var port int
		port, err = strconv.Atoi(matchMap["port"])
		if err != nil {
			return
		}

		node := &api.NodeData{
			NodeBaseData: api.NodeBaseData{
				Name:                matchMap["nodeName"],
				MachineRoles:        roles,
				DockerRootDirectory: matchMap["dockerPath"],
			},
			ConnectionData: api.ConnectionData{
				SSHLoginData: loginData,
				IP:           matchMap["ip"],
				Port:         uint16(port),
			},
		}
		log.ReqEntry(c).Tracef("node: %#v", node)

		err = node.Validate()
		if err != nil {
			return
		}

		if _, exist := ipList[matchMap["ip"]]; exist {

			err = fmt.Errorf("node ip %s was duplicated", matchMap["ip"])
			return
		}

		if _, exist := nameList[matchMap["nodeName"]]; exist {

			err = fmt.Errorf("node name %s was duplicated", matchMap["nodeName"])
			return
		}

		ipList[matchMap["ip"]] = true
		nameList[matchMap["nodeName"]] = true

		nodeList = append(nodeList, node)
	}

	log.ReqEntry(c).WithField("inputNodeList", nodeList).Trace("input node list")

	return
}

func splitInputRoles(matchMap map[string]string) []constant.MachineRole {
	roles := make([]constant.MachineRole, 0)
	if rolesString, exist := matchMap["roles"]; exist {
		logrus.WithField("func", "splitInputRoles").Tracef("rolesString: %s", rolesString)
		splitRoles := strings.Split(rolesString, ",")
		logrus.WithField("func", "splitInputRoles").Tracef("splitString: %s", splitRoles)
		for _, role := range splitRoles {
			trimmedRole := strings.TrimSpace(role)
			if role == "" {
				continue
			}
			roles = append(roles, constant.MachineRole(trimmedRole))
		}
	}
	return roles
}

func tryToMatchBatchNodes(data []byte) ([][]string, []string) {
	re := regexp.MustCompile(`(?m)^\s*` +
		`(?P<nodeName>[\w\-]+)\s+` +
		`(?P<username>[\w\-]+)\s+` +
		`(?P<roles>[\w,]+)\s+` +
		`(?P<ip>[\d.]+)\s+` +
		`(?P<port>[\d]+)\s+` +
		`(?P<password>[\w` + "`" + `~!@#$%^&*()\-+=\\|\[\]{};:'",./<>?]+)\s+` +
		`(?P<privateKeyName>[\-\w]+)\s+` +
		`(?P<dockerPath>[\w\-\/]+)`,
	)
	return re.FindAllStringSubmatch(string(data), -1), re.SubexpNames()
}
