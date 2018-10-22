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
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	api "github.com/nlnwa/maalfrid-api/gen/go/maalfrid/service/aggregator"
)

type Client struct {
	address string
	cc      *grpc.ClientConn
	client  api.AggregatorClient
}

func NewClient(address string) *Client {
	return &Client{address: address}
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

func (ac *Client) RunLanguageDetection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := ac.client.RunLanguageDetection(ctx, &empty.Empty{})
	return err
}

func (ac *Client) RunAggregation() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := ac.client.RunAggregation(ctx, &empty.Empty{})
	return err
}

func (ac *Client) SyncEntities() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := ac.client.SyncEntities(ctx, &empty.Empty{})
	return err
}

func (ac *Client) SyncSeeds() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := ac.client.SyncSeeds(ctx, &empty.Empty{})
	return err
}
