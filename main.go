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

package main

import (
	"context"
	"fmt"
	"github.com/namsral/flag"
	"os"

	"github.com/nlnwa/maalfrid-aggregator-client/pkg/aggregator"
	myFlag "github.com/nlnwa/maalfrid-aggregator-client/pkg/flag"
	"github.com/nlnwa/maalfrid-aggregator-client/pkg/version"
)

func main() {
	// global command parameters
	serviceHost := "localhost"
	servicePort := 3011

	// aggregate command parameters
	aggregateJobExecutionId := ""

	// filter command parameters
	filterJobExecutionId := ""
	filterSeedID := ""

	// sync command parameters
	var seedLabels myFlag.ArrayFlag

	// detect command parameters
	detectAll := false

	// global flags
	flag.StringVar(&serviceHost, "host", serviceHost, "maalfrid aggregator service host")
	flag.IntVar(&servicePort, "port", servicePort, "maalfrid aggregator service port")
	flag.Parse()

	address := fmt.Sprintf("%s:%d", serviceHost, servicePort)

	// aggregate command flags
	aggregateCommand := flag.NewFlagSet("aggregate", flag.ExitOnError)
	aggregateCommand.StringVar(&aggregateJobExecutionId, "job-execution-id", aggregateJobExecutionId, "id of job execution to aggregate")

	// sync command flags
	syncCommand := flag.NewFlagSet("sync", flag.ExitOnError)
	syncCommand.Var(&seedLabels, "label", "label selector on key:value format (can be specified multiple times)")

	// aggregate command flags
	filterCommand := flag.NewFlagSet("filter", flag.ExitOnError)
	filterCommand.StringVar(&filterSeedID, "seed-id", filterSeedID, "limit filtering to seed with this id")
	filterCommand.StringVar(&filterJobExecutionId, "job-execution-id", filterJobExecutionId, "id of job execution to aggregate")

	// detect command flags
	detectCommand := flag.NewFlagSet("detect", flag.ExitOnError)
	detectCommand.BoolVar(&detectAll, "all", detectAll, "if language detection should process already detected texts")

	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	var err error
	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Println(version.String())
		os.Exit(0)
	case "aggregate":
		err = aggregateCommand.Parse(os.Args[2:])
	case "sync":
		err = syncCommand.Parse(os.Args[2:])
	case "detect":
		err = detectCommand.Parse(os.Args[2:])
	case "filter":
		err = filterCommand.Parse(os.Args[2:])
	default:
		err = fmt.Errorf("unknown command \"%s\"", cmd)
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		usage()
		os.Exit(1)
	}

	if aggregateCommand.Parsed() {
		err = runAggregation(address, aggregateJobExecutionId)
	} else if syncCommand.Parsed() {
		err = syncEntities(address, seedLabels)
	} else if detectCommand.Parsed() {
		err = runLanguageDetection(address, detectAll)
	} else if filterCommand.Parsed() {
		err = filterAggregate(address, filterJobExecutionId, filterSeedID)
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func syncEntities(address string, labels []string) error {
	client := aggregator.NewClient(address)
	if err := client.Dial(); err != nil {
		return err
	}
	defer func() { _ = client.Hangup() }()
	return client.SyncEntities(context.Background(), labels)
}

func runAggregation(address string, jobExecutionId string) error {
	client := aggregator.NewClient(address)
	if err := client.Dial(); err != nil {
		return err
	}
	defer func() { _ = client.Hangup() }()
	return client.RunAggregation(context.Background(), jobExecutionId)
}

func filterAggregate(address string, jobExecutionId string, seedID string) error {
	client := aggregator.NewClient(address)
	if err := client.Dial(); err != nil {
		return err
	}
	defer func() { _ = client.Hangup() }()
	return client.FilterAggregate(context.Background(), jobExecutionId, seedID)
}

func runLanguageDetection(address string, detectAll bool) error {
	client := aggregator.NewClient(address)
	if err := client.Dial(); err != nil {
		return err
	}
	defer func() { _ = client.Hangup() }()
	return client.RunLanguageDetection(context.Background(), detectAll)
}

func usage() {
	fmt.Println(`Usage:
	maalfrid-aggregator-client <command>
	
Commands:
  detect	Initiate language detection of extracted texts in veidemann
  aggregate	Initiate aggregation of data from veidemann
  sync		Sync entities and seeds from veidemann
  filter	Filter aggregated data
  version	Print version information

Use "maalfrid-aggregator-client <command> -help" for more information about a given command.
Use "maalfrid-aggregator-client -help" for a list of global command-line options (applies to all commands).`)
}
