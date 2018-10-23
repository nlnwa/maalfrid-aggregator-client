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

	"github.com/namsral/flag"

	"github.com/nlnwa/maalfrid-aggregator-client/pkg/aggregator"
	"github.com/nlnwa/maalfrid-aggregator-client/pkg/version"
)

func main() {
	serviceHost := "localhost"
	servicePort := 3011

	flag.StringVar(&serviceHost, "host", serviceHost, "maalfrid aggregator service host")
	flag.IntVar(&servicePort, "port", servicePort, "maalfrid aggregator service port")
	flag.Parse()

	address := fmt.Sprintf("%s:%d", serviceHost, servicePort)

	if len(os.Args) < 2 {
		usage()
	}
	cmd := os.Args[1]
	switch cmd {
	case
		"detect",
		"aggregate",
		"sync-seeds",
		"sync-entities",
		"version":
		if err := runCommand(cmd, address); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	default:
		usage()
	}
}

func runCommand(cmd string, address string) error {
	client := aggregator.NewClient(address)

	if err := client.Dial(); err != nil {
		return err
	}
	defer client.Hangup()

	switch cmd {
	case "detect":
		return client.RunLanguageDetection()
	case "aggregate":
		return client.RunAggregation()
	case "sync-seeds":
		return client.SyncSeeds()
	case "sync-entities":
		return client.SyncEntities()
	case "version":
		_, err := fmt.Println(version.String())
		return err
	default:
		return fmt.Errorf("%s: %s", "Unknown subcommand", cmd)
	}
}

func usage() {
	fmt.Println(`Usage:
	maalfrid-aggregator-client <command>
	
Commands:
	detect		Initiate language detection of extracted texts (missing language code)
	aggregate	Initiate aggregation of data from veidemann to maalfrid
	sync-entities	Sync entities from veidemann to maalfrid
	sync-seeds	Sync seeds from veidemann to maalfrid
	version		Print version information`)
	os.Exit(1)
}
