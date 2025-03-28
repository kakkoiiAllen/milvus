// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpcdatacoord

import (
	"context"
	"errors"
	"testing"

	"github.com/milvus-io/milvus/internal/proto/commonpb"
	"github.com/milvus-io/milvus/internal/proto/datapb"
	"github.com/milvus-io/milvus/internal/proto/internalpb"
	"github.com/milvus-io/milvus/internal/proto/milvuspb"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type MockDataCoord struct {
	states               *internalpb.ComponentStates
	status               *commonpb.Status
	err                  error
	initErr              error
	startErr             error
	stopErr              error
	regErr               error
	strResp              *milvuspb.StringResponse
	infoResp             *datapb.GetSegmentInfoResponse
	flushResp            *datapb.FlushResponse
	assignResp           *datapb.AssignSegmentIDResponse
	segStateResp         *datapb.GetSegmentStatesResponse
	binResp              *datapb.GetInsertBinlogPathsResponse
	colStatResp          *datapb.GetCollectionStatisticsResponse
	partStatResp         *datapb.GetPartitionStatisticsResponse
	recoverResp          *datapb.GetRecoveryInfoResponse
	flushSegResp         *datapb.GetFlushedSegmentsResponse
	metricResp           *milvuspb.GetMetricsResponse
	compactionStateResp  *milvuspb.GetCompactionStateResponse
	manualCompactionResp *milvuspb.ManualCompactionResponse
	compactionPlansResp  *milvuspb.GetCompactionPlansResponse
	watchChannelsResp    *datapb.WatchChannelsResponse
	getFlushStateResp    *milvuspb.GetFlushStateResponse
	dropVChanResp        *datapb.DropVirtualChannelResponse
	setSegmentStateResp  *datapb.SetSegmentStateResponse
	importResp           *datapb.ImportTaskResponse
	updateSegStatResp    *commonpb.Status
	acquireSegLockResp   *commonpb.Status
	releaseSegLockResp   *commonpb.Status
	addSegmentResp       *commonpb.Status
}

func (m *MockDataCoord) Init() error {
	return m.initErr
}

func (m *MockDataCoord) Start() error {
	return m.startErr
}

func (m *MockDataCoord) Stop() error {
	return m.stopErr
}

func (m *MockDataCoord) Register() error {
	return m.regErr
}

func (m *MockDataCoord) SetEtcdClient(etcdClient *clientv3.Client) {
}

func (m *MockDataCoord) GetComponentStates(ctx context.Context) (*internalpb.ComponentStates, error) {
	return m.states, m.err
}

func (m *MockDataCoord) GetTimeTickChannel(ctx context.Context) (*milvuspb.StringResponse, error) {
	return m.strResp, m.err
}

func (m *MockDataCoord) GetStatisticsChannel(ctx context.Context) (*milvuspb.StringResponse, error) {
	return m.strResp, m.err
}

func (m *MockDataCoord) GetSegmentInfo(ctx context.Context, req *datapb.GetSegmentInfoRequest) (*datapb.GetSegmentInfoResponse, error) {
	return m.infoResp, m.err
}

func (m *MockDataCoord) Flush(ctx context.Context, req *datapb.FlushRequest) (*datapb.FlushResponse, error) {
	return m.flushResp, m.err
}

func (m *MockDataCoord) AssignSegmentID(ctx context.Context, req *datapb.AssignSegmentIDRequest) (*datapb.AssignSegmentIDResponse, error) {
	return m.assignResp, m.err
}

func (m *MockDataCoord) GetSegmentStates(ctx context.Context, req *datapb.GetSegmentStatesRequest) (*datapb.GetSegmentStatesResponse, error) {
	return m.segStateResp, m.err
}

func (m *MockDataCoord) GetInsertBinlogPaths(ctx context.Context, req *datapb.GetInsertBinlogPathsRequest) (*datapb.GetInsertBinlogPathsResponse, error) {
	return m.binResp, m.err
}

func (m *MockDataCoord) GetCollectionStatistics(ctx context.Context, req *datapb.GetCollectionStatisticsRequest) (*datapb.GetCollectionStatisticsResponse, error) {
	return m.colStatResp, m.err
}

func (m *MockDataCoord) GetPartitionStatistics(ctx context.Context, req *datapb.GetPartitionStatisticsRequest) (*datapb.GetPartitionStatisticsResponse, error) {
	return m.partStatResp, m.err
}

func (m *MockDataCoord) GetSegmentInfoChannel(ctx context.Context) (*milvuspb.StringResponse, error) {
	return m.strResp, m.err
}

func (m *MockDataCoord) SaveBinlogPaths(ctx context.Context, req *datapb.SaveBinlogPathsRequest) (*commonpb.Status, error) {
	return m.status, m.err
}

func (m *MockDataCoord) GetRecoveryInfo(ctx context.Context, req *datapb.GetRecoveryInfoRequest) (*datapb.GetRecoveryInfoResponse, error) {
	return m.recoverResp, m.err
}

func (m *MockDataCoord) GetFlushedSegments(ctx context.Context, req *datapb.GetFlushedSegmentsRequest) (*datapb.GetFlushedSegmentsResponse, error) {
	return m.flushSegResp, m.err
}

func (m *MockDataCoord) GetMetrics(ctx context.Context, req *milvuspb.GetMetricsRequest) (*milvuspb.GetMetricsResponse, error) {
	return m.metricResp, m.err
}

func (m *MockDataCoord) CompleteCompaction(ctx context.Context, req *datapb.CompactionResult) (*commonpb.Status, error) {
	return m.status, m.err
}

func (m *MockDataCoord) ManualCompaction(ctx context.Context, req *milvuspb.ManualCompactionRequest) (*milvuspb.ManualCompactionResponse, error) {
	return m.manualCompactionResp, m.err
}

func (m *MockDataCoord) GetCompactionState(ctx context.Context, req *milvuspb.GetCompactionStateRequest) (*milvuspb.GetCompactionStateResponse, error) {
	return m.compactionStateResp, m.err
}

func (m *MockDataCoord) GetCompactionStateWithPlans(ctx context.Context, req *milvuspb.GetCompactionPlansRequest) (*milvuspb.GetCompactionPlansResponse, error) {
	return m.compactionPlansResp, m.err
}

func (m *MockDataCoord) WatchChannels(ctx context.Context, req *datapb.WatchChannelsRequest) (*datapb.WatchChannelsResponse, error) {
	return m.watchChannelsResp, m.err
}

func (m *MockDataCoord) GetFlushState(ctx context.Context, req *milvuspb.GetFlushStateRequest) (*milvuspb.GetFlushStateResponse, error) {
	return m.getFlushStateResp, m.err
}

func (m *MockDataCoord) DropVirtualChannel(ctx context.Context, req *datapb.DropVirtualChannelRequest) (*datapb.DropVirtualChannelResponse, error) {
	return m.dropVChanResp, m.err
}

func (m *MockDataCoord) SetSegmentState(ctx context.Context, req *datapb.SetSegmentStateRequest) (*datapb.SetSegmentStateResponse, error) {
	return m.setSegmentStateResp, m.err
}

func (m *MockDataCoord) Import(ctx context.Context, req *datapb.ImportTaskRequest) (*datapb.ImportTaskResponse, error) {
	return m.importResp, m.err
}

func (m *MockDataCoord) UpdateSegmentStatistics(ctx context.Context, req *datapb.UpdateSegmentStatisticsRequest) (*commonpb.Status, error) {
	return m.updateSegStatResp, m.err
}

func (m *MockDataCoord) AcquireSegmentLock(ctx context.Context, req *datapb.AcquireSegmentLockRequest) (*commonpb.Status, error) {
	return m.acquireSegLockResp, m.err
}

func (m *MockDataCoord) ReleaseSegmentLock(ctx context.Context, req *datapb.ReleaseSegmentLockRequest) (*commonpb.Status, error) {
	return m.releaseSegLockResp, m.err
}

func (m *MockDataCoord) AddSegment(ctx context.Context, req *datapb.AddSegmentRequest) (*commonpb.Status, error) {
	return m.addSegmentResp, m.err
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Test_NewServer(t *testing.T) {
	ctx := context.Background()
	server := NewServer(ctx, nil)
	assert.NotNil(t, server)

	t.Run("Run", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{}
		err := server.Run()
		assert.Nil(t, err)
	})

	t.Run("GetComponentStates", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			states: &internalpb.ComponentStates{},
		}
		states, err := server.GetComponentStates(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, states)
	})

	t.Run("GetTimeTickChannel", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			strResp: &milvuspb.StringResponse{},
		}
		resp, err := server.GetTimeTickChannel(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetStatisticsChannel", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			strResp: &milvuspb.StringResponse{},
		}
		resp, err := server.GetStatisticsChannel(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetSegmentInfo", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			infoResp: &datapb.GetSegmentInfoResponse{},
		}
		resp, err := server.GetSegmentInfo(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("Flush", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			flushResp: &datapb.FlushResponse{},
		}
		resp, err := server.Flush(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("AssignSegmentID", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			assignResp: &datapb.AssignSegmentIDResponse{},
		}
		resp, err := server.AssignSegmentID(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetSegmentStates", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			segStateResp: &datapb.GetSegmentStatesResponse{},
		}
		resp, err := server.GetSegmentStates(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetInsertBinlogPaths", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			binResp: &datapb.GetInsertBinlogPathsResponse{},
		}
		resp, err := server.GetInsertBinlogPaths(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetCollectionStatistics", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			colStatResp: &datapb.GetCollectionStatisticsResponse{},
		}
		resp, err := server.GetCollectionStatistics(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetPartitionStatistics", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			partStatResp: &datapb.GetPartitionStatisticsResponse{},
		}
		resp, err := server.GetPartitionStatistics(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetSegmentInfoChannel", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			strResp: &milvuspb.StringResponse{},
		}
		resp, err := server.GetSegmentInfoChannel(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("SaveBinlogPaths", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			status: &commonpb.Status{},
		}
		resp, err := server.SaveBinlogPaths(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetRecoveryInfo", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			recoverResp: &datapb.GetRecoveryInfoResponse{},
		}
		resp, err := server.GetRecoveryInfo(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetFlushedSegments", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			flushSegResp: &datapb.GetFlushedSegmentsResponse{},
		}
		resp, err := server.GetFlushedSegments(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetMetrics", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			metricResp: &milvuspb.GetMetricsResponse{},
		}
		resp, err := server.GetMetrics(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("WatchChannels", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			watchChannelsResp: &datapb.WatchChannelsResponse{},
		}
		resp, err := server.WatchChannels(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetFlushState", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			getFlushStateResp: &milvuspb.GetFlushStateResponse{},
		}
		resp, err := server.GetFlushState(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("DropVirtualChannel", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			dropVChanResp: &datapb.DropVirtualChannelResponse{},
		}
		resp, err := server.DropVirtualChannel(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("CompleteCompaction", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			status: &commonpb.Status{},
		}
		resp, err := server.CompleteCompaction(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("ManualCompaction", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			manualCompactionResp: &milvuspb.ManualCompactionResponse{},
		}
		resp, err := server.ManualCompaction(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetCompactionState", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			compactionStateResp: &milvuspb.GetCompactionStateResponse{},
		}
		resp, err := server.GetCompactionState(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetCompactionStateWithPlans", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			compactionPlansResp: &milvuspb.GetCompactionPlansResponse{},
		}
		resp, err := server.GetCompactionStateWithPlans(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("set segment state", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			setSegmentStateResp: &datapb.SetSegmentStateResponse{},
		}
		resp, err := server.SetSegmentState(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("import", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			importResp: &datapb.ImportTaskResponse{
				Status: &commonpb.Status{},
			},
		}
		resp, err := server.Import(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update seg stat", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			updateSegStatResp: &commonpb.Status{
				ErrorCode: commonpb.ErrorCode_Success,
			},
		}
		resp, err := server.UpdateSegmentStatistics(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("acquire segment reference lock", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			acquireSegLockResp: &commonpb.Status{
				ErrorCode: commonpb.ErrorCode_Success,
			},
		}
		resp, err := server.AcquireSegmentLock(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("release segment reference lock", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			releaseSegLockResp: &commonpb.Status{
				ErrorCode: commonpb.ErrorCode_Success,
			},
		}
		resp, err := server.ReleaseSegmentLock(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("add segment", func(t *testing.T) {
		server.dataCoord = &MockDataCoord{
			addSegmentResp: &commonpb.Status{
				ErrorCode: commonpb.ErrorCode_Success,
			},
		}
		resp, err := server.AddSegment(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	err := server.Stop()
	assert.Nil(t, err)
}

func Test_Run(t *testing.T) {
	ctx := context.Background()
	server := NewServer(ctx, nil)
	assert.NotNil(t, server)

	server.dataCoord = &MockDataCoord{
		regErr: errors.New("error"),
	}

	err := server.Run()
	assert.Error(t, err)

	server.dataCoord = &MockDataCoord{
		startErr: errors.New("error"),
	}

	err = server.Run()
	assert.Error(t, err)

	server.dataCoord = &MockDataCoord{
		initErr: errors.New("error"),
	}

	err = server.Run()
	assert.Error(t, err)

	server.dataCoord = &MockDataCoord{
		stopErr: errors.New("error"),
	}

	err = server.Stop()
	assert.Error(t, err)
}
