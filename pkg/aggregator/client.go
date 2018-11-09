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

type Client struct {
	address string
	cc      *grpc.ClientConn
	client  api.AggregatorClient
	timeout time.Duration
}

func NewClient(address string) *Client {
	return &Client{address: address, timeout: 10}
}

func (ac *Client) Dial() (err error) {
	if ac.cc, err = grpc.Dial(ac.address, grpc.WithInsecure()); err != nil {
		return errors.Wrapf(err, "failed to dial: %s", ac.address)
	} else {
		ac.client = api.NewAggregatorClient(ac.cc)
		return
	}
}

func (ac *Client) Hangup() error {
	if ac.cc != nil {
		return ac.cc.Close()
	} else {
		return nil
	}
}

func (ac *Client) RunLanguageDetection(detectAll bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ac.timeout)
	defer cancel()

	req := &api.RunLanguageDetectionRequest{
		DetectAll: detectAll,
	}
	if _, err := ac.client.RunLanguageDetection(ctx, req); err != nil {
		return errors.Wrapf(err, "failed to run language detection")
	} else {
		return nil
	}
}

func (ac *Client) RunAggregation(startTime time.Time, endTime time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ac.timeout)
	defer cancel()

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
	} else {
		return nil
	}
}

func (ac *Client) FilterAggregate(startTime time.Time, endTime time.Time, seedId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ac.timeout)
	defer cancel()

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
		SeedId:    seedId,
	}
	if _, err := ac.client.FilterAggregate(ctx, req); err != nil {
		return errors.Wrapf(err, "failed to filter aggregate")
	} else {
		return nil
	}
}

func (ac *Client) SyncEntities(name string, labels []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ac.timeout)
	defer cancel()

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
	} else {
		return nil
	}
}
