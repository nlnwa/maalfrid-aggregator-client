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
	"fmt"
	"os"
	"time"

	"github.com/namsral/flag"

	"github.com/nlnwa/maalfrid-aggregator-client/pkg/aggregator"
	myFlag "github.com/nlnwa/maalfrid-aggregator-client/pkg/flag"
	"github.com/nlnwa/maalfrid-aggregator-client/pkg/version"
)

func main() {
	// global command parameters
	serviceHost := "localhost"
	servicePort := 3011

	// aggregate command parameters
	aggregateStartTime := ""
	aggregateEndTime := ""

	// filter command parameters
	filterStartTime := ""
	filterEndTime := ""

	// sync command parameters
	var entityLabels myFlag.ArrayFlag
	entityName := ""

	// detect command parameters
	detectAll := false

	// global flags
	flag.StringVar(&serviceHost, "host", serviceHost, "maalfrid aggregator service host")
	flag.IntVar(&servicePort, "port", servicePort, "maalfrid aggregator service port")
	flag.Parse()

	address := fmt.Sprintf("%s:%d", serviceHost, servicePort)

	// aggregate command flags
	aggregateCommand := flag.NewFlagSet("aggregate", flag.ExitOnError)
	aggregateCommand.StringVar(&aggregateStartTime, "start-time", "", "lower bound of execution start time in RFC3339 format (inclusive)")
	aggregateCommand.StringVar(&aggregateEndTime, "end-time", "", "upper bound of execution start time in RFC3339 format (exclusive)")

	// sync command flags
	syncCommand := flag.NewFlagSet("sync", flag.ExitOnError)
	syncCommand.Var(&entityLabels, "label", "label selector on key:value format (can be specified multiple times)")
	syncCommand.StringVar(&entityName, "name", entityName, "only synchronize entity with matching name")

	// aggregate command flags
	filterCommand := flag.NewFlagSet("filter", flag.ExitOnError)
	filterCommand.StringVar(&aggregateStartTime, "start-time", "", "lower bound of execution start time in RFC3339 format (inclusive)")
	filterCommand.StringVar(&aggregateEndTime, "end-time", "", "upper bound of execution start time in RFC3339 format (exclusive)")

	// detect command flags
	detectCommand := flag.NewFlagSet("detect", flag.ExitOnError)
	detectCommand.BoolVar(&detectAll, "all", detectAll, "if language detection should process already detected texts")

	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Println(version.String())
		os.Exit(0)
	case "aggregate":
		aggregateCommand.Parse(os.Args[2:])
	case "sync":
		syncCommand.Parse(os.Args[2:])
	case "detect":
		detectCommand.Parse(os.Args[2:])
	case "filter":
		filterCommand.Parse(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Error: %s \"%s\"\n", "unknown command", cmd)
		fmt.Println("Run 'maalfrid-aggregator-client' (without subcommand) for usage")
		os.Exit(1)
	}

	var err error
	if aggregateCommand.Parsed() {
		err = runAggregation(address, aggregateStartTime, aggregateEndTime)
	}
	if syncCommand.Parsed() {
		err = syncEntities(address, entityLabels, entityName)
	}
	if detectCommand.Parsed() {
		err = runLanguageDetection(address, detectAll)
	}
	if filterCommand.Parsed() {
		err = filterAggregate(address, filterStartTime, filterEndTime)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func syncEntities(address string, labels []string, name string) error {
	client := aggregator.NewClient(address)
	if err := client.Dial(); err != nil {
		return err
	}
	defer client.Hangup()
	return client.SyncEntities(name, labels)
}

func runAggregation(address string, startTimeString string, endTimeString string) error {
	client := aggregator.NewClient(address)
	if err := client.Dial(); err != nil {
		return err
	}
	defer client.Hangup()
	var startTime time.Time
	var endTime time.Time
	var err error
	if len(startTimeString) > 0 {
		startTime, err = time.Parse(time.RFC3339, startTimeString)
	}
	if len(endTimeString) > 0 {
		endTime, err = time.Parse(time.RFC3339, endTimeString)
	}
	if err != nil {
		return err
	}
	return client.RunAggregation(startTime, endTime)
}

func filterAggregate(address string, startTimeString string, endTimeString string) error {
	client := aggregator.NewClient(address)
	if err := client.Dial(); err != nil {
		return err
	}
	defer client.Hangup()
	var startTime time.Time
	var endTime time.Time
	var err error
	if len(startTimeString) > 0 {
		startTime, err = time.Parse(time.RFC3339, startTimeString)
	}
	if len(endTimeString) > 0 {
		endTime, err = time.Parse(time.RFC3339, endTimeString)
	}
	if err != nil {
		return err
	}
	return client.FilterAggregate(startTime, endTime)
}

func runLanguageDetection(address string, detectAll bool) error {
	client := aggregator.NewClient(address)
	if err := client.Dial(); err != nil {
		return err
	}
	defer client.Hangup()
	return client.RunLanguageDetection(detectAll)
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

Use "maalfrid-aggregator-client <command> --help" for more information about a given command.
Use "maalfrid-aggregator-client --help" for a list of global command-line options (applies to all commands).`)
}
