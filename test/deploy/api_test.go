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
	"context"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/kpaas-io/kpaas/pkg/constant"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
	"github.com/kpaas-io/kpaas/pkg/deploy/server"
	"github.com/sirupsen/logrus"
)

var (
	client pb.DeployContollerClient
	conn   *grpc.ClientConn
	stopCh chan struct{}
)

func setup() {
	if _testConfig.Skip {
		return
	}

	logrus.SetLevel(logrus.DebugLevel)

	serverAddress := _testConfig.RemoteServerAddress

	if _testConfig.LaunchLocalServer {
		var port uint16 = 9999
		// Setup and start gRpc server
		options := server.ServerOptions{
			Port:       port,
			LogFileLoc: "./tmp/logs",
		}
		stopCh = make(chan struct{})
		go server.New(options).Run(stopCh)
		serverAddress = fmt.Sprintf("localhost:%d", port)
	}

	// Create a gRpc client
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var err error
	conn, err = grpc.DialContext(ctx, serverAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println("did not connect:", err)
		os.Exit(1)
	}
	client = pb.NewDeployContollerClient(conn)
}

func tearDown() {
	if _testConfig.Skip {
		return
	}

	conn.Close()

	if _testConfig.LaunchLocalServer {
		stopCh <- struct{}{}
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func TestTestConnection(t *testing.T) {
	if _testConfig.Skip {
		t.SkipNow()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	request, expetecdReply := getTestConnectionData()
	actualReply, err := client.TestConnection(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, actualReply)
	assert.Equal(t, expetecdReply, actualReply)
}

func TestCheckNodes(t *testing.T) {
	if _testConfig.Skip {
		t.SkipNow()
	}

	// CheckNodes request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	checkRequest, expetecdCheckReply := getCheckNodesData()
	actualCheckReply, err := client.CheckNodes(ctx, checkRequest)
	assert.NoError(t, err)
	assert.NotNil(t, actualCheckReply)
	assert.Equal(t, expetecdCheckReply, actualCheckReply)

	// GetCheckNodesResult request
	var actualResultReply *pb.GetCheckNodesResultReply
	resultRequest, expetecdResultReply := getGetCheckNodesResultData()
	// Call GetCheckNodesResult repeatly until the related task is done or failed.
	err = wait.Poll(3*time.Second, 1*time.Minute, func() (done bool, err error) {
		ctxPoll, cancelPoll := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelPoll()
		actualResultReply, err = client.GetCheckNodesResult(ctxPoll, resultRequest)
		if err != nil {
			return false, err
		}
		if actualResultReply.Status == string(constant.OperationStatusFailed) ||
			actualResultReply.Status == string(constant.OperationStatusSuccessful) {
			return true, nil
		}
		return false, nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, actualResultReply)
	sortCheckNodesResult(actualResultReply)
	sortCheckNodesResult(expetecdResultReply)
	assert.Equal(t, expetecdResultReply, actualResultReply)

	// Test GetCheckNodesLog
	ctxGetLog, cancelGetLog := context.WithTimeout(context.Background(), 10000*time.Second)
	defer cancelGetLog()
	requestGetLog, _ := getGetCheckNodesLogData()
	actualGetLogReply, errGetLog := client.GetCheckNodesLog(ctxGetLog, requestGetLog)
	assert.NoError(t, errGetLog)
	assert.NotNil(t, actualGetLogReply)
	logStr := string(actualGetLogReply.Log)
	t.Log(logStr)
	// Just a simple check on the content of the log
	assert.Equal(t, true, len(actualGetLogReply.Log) > 100)
}

func TestDeployMultipleNodes(t *testing.T) {
	if _testConfig.Skip {
		t.SkipNow()
	}

	// Deploy request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	deployRequest, expetecdDeployReply := getDeployMultipleNodesData()
	actualDeployReply, err := client.Deploy(ctx, deployRequest)
	assert.NoError(t, err)
	assert.NotNil(t, actualDeployReply)
	assert.Equal(t, expetecdDeployReply, actualDeployReply)

	// GetDeployResult request
	var actualResultReply *pb.GetDeployResultReply
	resultRequest, expetecdResultReply := getDeployResultData()
	// Call GetDeployResult repeatly until the related task is done or failed.
	err = wait.Poll(5*time.Second, 10*time.Minute, func() (done bool, err error) {
		ctxPoll, cancelPoll := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelPoll()
		actualResultReply, err = client.GetDeployResult(ctxPoll, resultRequest)
		if err != nil {
			return false, err
		}
		if actualResultReply.Status == string(constant.OperationStatusFailed) ||
			actualResultReply.Status == string(constant.OperationStatusSuccessful) {
			return true, nil
		}
		return false, nil
	})
	assert.NoError(t, err)
	sortDeployResults(expetecdResultReply)
	sortDeployResults(actualResultReply)
	assert.NotNil(t, actualResultReply)
	assert.Equal(t, expetecdResultReply, actualResultReply)

	// Test GetDeployLog
	ctxGetLog, cancelGetLog := context.WithTimeout(context.Background(), 10000*time.Second)
	defer cancelGetLog()
	requestGetLog, _ := getGetDeployLogData()
	actualGetLogReply, errGetLog := client.GetDeployLog(ctxGetLog, requestGetLog)
	assert.NoError(t, errGetLog)
	assert.NotNil(t, actualGetLogReply)
	logStr := string(actualGetLogReply.Log)
	t.Log(logStr)
	// Just a simple check on the content of the log
	assert.Equal(t, true, len(actualGetLogReply.Log) > 100)
}

func TestDeployAllInOne(t *testing.T) {
	if _testConfig.Skip {
		t.SkipNow()
	}

	// Deploy request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	deployRequest, expetecdDeployReply := getDeployAllInOneData()
	actualDeployReply, err := client.Deploy(ctx, deployRequest)
	assert.NoError(t, err)
	assert.NotNil(t, actualDeployReply)
	assert.Equal(t, expetecdDeployReply, actualDeployReply)

	// GetDeployResult request
	var actualResultReply *pb.GetDeployResultReply
	resultRequest, expetecdResultReply := getDeployAllInOneResultData()
	// Call GetDeployResult repeatly until the related task is done or failed.
	err = wait.Poll(5*time.Second, 10*time.Minute, func() (done bool, err error) {
		ctxPoll, cancelPoll := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelPoll()
		actualResultReply, err = client.GetDeployResult(ctxPoll, resultRequest)
		if err != nil {
			return false, err
		}
		if actualResultReply.Status == string(constant.OperationStatusFailed) ||
			actualResultReply.Status == string(constant.OperationStatusSuccessful) {
			return true, nil
		}
		return false, nil
	})
	assert.NoError(t, err)
	sortDeployResults(expetecdResultReply)
	sortDeployResults(actualResultReply)
	assert.NotNil(t, actualResultReply)
	assert.Equal(t, expetecdResultReply, actualResultReply)

	// Test GetDeployLog
	ctxGetLog, cancelGetLog := context.WithTimeout(context.Background(), 10000*time.Second)
	defer cancelGetLog()
	requestGetLog, _ := getGetDeployLogData()
	actualGetLogReply, errGetLog := client.GetDeployLog(ctxGetLog, requestGetLog)
	assert.NoError(t, errGetLog)
	assert.NotNil(t, actualGetLogReply)
	logStr := string(actualGetLogReply.Log)
	t.Log(logStr)
	// Just a simple check on the content of the log
	assert.Equal(t, true, len(actualGetLogReply.Log) > 100)
}

func TestFetchKubeConfig(t *testing.T) {
	if _testConfig.Skip {
		t.SkipNow()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	request, _ := getFetchKubeConfigData()
	actualReply, err := client.FetchKubeConfig(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, actualReply)
	// Just a simple check on the content of kube config
	assert.Equal(t, true, len(actualReply.KubeConfig) > 1000)
}

func sortItemCheckResults(results []*pb.ItemCheckResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Item.Name <= results[j].Item.Name
	})
}

func sortCheckNodesResult(r *pb.GetCheckNodesResultReply) {
	for _, nodeCheckResult := range r.Nodes {
		sortItemCheckResults(nodeCheckResult.Items)
	}
}

func sortDeployItemResults(items []*pb.DeployItemResult) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].DeployItem.Role == items[j].DeployItem.Role {
			return items[i].DeployItem.NodeName <= items[j].DeployItem.NodeName
		}
		return items[i].DeployItem.Role < items[j].DeployItem.Role
	})
}

func sortDeployResults(result *pb.GetDeployResultReply) {
	if result == nil {
		return
	}

	sortDeployItemResults(result.Items)
}
