// Copyright 2018 National Library of Norway
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

// Package aggregator contains an aggregator service client
package aggregator

import (
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/timestamp"
	api "github.com/nlnwa/maalfrid-api/gen/go/maalfrid/service/aggregator"
)

// Client represents the client to the aggregator service.
type Client struct {
	address string               // address in the form "host:port"
	cc      *grpc.ClientConn     // gRPC connection
	client  api.AggregatorClient // gRPC client
	timeout time.Duration        // timeout duration
}

// NewClient creates a new client with the specified address and timeout.
func NewClient(address string) *Client {
	return &Client{address: address, timeout: 10}
}

// Dial makes a connection to the gRPC service.
func (ac *Client) Dial() (err error) {
	if ac.cc, err = grpc.Dial(ac.address, grpc.WithInsecure()); err != nil {
		return errors.Wrapf(err, "failed to dial: %s", ac.address)
	}
	ac.client = api.NewAggregatorClient(ac.cc)
	return
}

// Hangup closes the connection to the gRPC service.
func (ac *Client) Hangup() error {
	if ac.cc != nil {
		return ac.cc.Close()
	}
	return nil
}

// RunLanguageDetection calls the gRPC method with the same name.
func (ac *Client) RunLanguageDetection(ctx context.Context, detectAll bool) error {
	req := &api.RunLanguageDetectionRequest{
		DetectAll: detectAll,
	}
	if _, err := ac.client.RunLanguageDetection(ctx, req); err != nil {
		return errors.Wrapf(err, "failed to run language detection")
	}
	return nil
}

// RunAggregation calls the gRPC method with the same name.
func (ac *Client) RunAggregation(ctx context.Context, startTime time.Time, endTime time.Time) error {
	var err error
	var defaultTime time.Time
	var startTimeProto *timestamp.Timestamp
	var endTimeProto *timestamp.Timestamp

	if defaultTime.Equal(startTime) {
		startTimeProto = nil
	} else {
		startTimeProto, err = ptypes.TimestampProto(startTime)
	}
	if defaultTime.Equal(endTime) {
		endTimeProto = nil
	} else {
		endTimeProto, err = ptypes.TimestampProto(endTime)
	}
	if err != nil {
		return err
	}
	req := &api.RunAggregationRequest{
		StartTime: startTimeProto,
		EndTime:   endTimeProto,
	}
	if _, err := ac.client.RunAggregation(ctx, req); err != nil {
		return errors.Wrapf(err, "failed to run aggregation")
	}
	return nil
}

// FilterAggregate calls the gRPC method with the same name.
func (ac *Client) FilterAggregate(ctx context.Context, startTime time.Time, endTime time.Time, seedID string) error {
	var err error
	var defaultTime time.Time
	var startTimeProto *timestamp.Timestamp
	var endTimeProto *timestamp.Timestamp

	if defaultTime.Equal(startTime) {
		startTimeProto = nil
	} else {
		startTimeProto, err = ptypes.TimestampProto(startTime)
	}
	if defaultTime.Equal(endTime) {
		endTimeProto = nil
	} else {
		endTimeProto, err = ptypes.TimestampProto(endTime)
	}
	if err != nil {
		return err
	}
	req := &api.FilterAggregateRequest{
		StartTime: startTimeProto,
		EndTime:   endTimeProto,
		SeedId:    seedID,
	}
	if _, err := ac.client.FilterAggregate(ctx, req); err != nil {
		return errors.Wrapf(err, "failed to filter aggregate")
	}
	return nil
}

// SyncEntities calls the gRPC method with the same name.
func (ac *Client) SyncEntities(ctx context.Context, name string, labels []string) error {
	var labelsProto []*api.Label
	for _, label := range labels {
		parts := strings.Split(label, ":")
		if len(parts) == 1 {
			labelsProto = append(labelsProto, &api.Label{Value: parts[0]})
		} else {
			labelsProto = append(labelsProto, &api.Label{Key: parts[0], Value: parts[1]})
		}
	}
	req := &api.SyncEntitiesRequest{
		Name:   name,
		Labels: labelsProto,
	}
	if _, err := ac.client.SyncEntities(ctx, req); err != nil {
		return errors.Wrapf(err, "failed to sync entities")
	}
	return nil
}
