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
	client := aggregator.NewClient(address)

	cmd := getSubcommand()

	err := runCommand(cmd, client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func runCommand(cmd string, client *aggregator.Client) error {

	if err := client.Dial(); err != nil {
		return err
	}
	defer client.Hangup()

	switch cmd {
	case "detect":
		return client.RunLanguageDetection()
	case "aggregate":
		return client.RunLanguageDetection()
	default:
		return fmt.Errorf("%s", "unknown command")
	}
}

func getSubcommand() string {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "detect":
		break
	case "aggregate":
		break
	case "version":
		fmt.Printf("%s, %s\n", "Aggregator client", version.String())
		os.Exit(0)
	case "sync-entities":
		notImplemented()
	case "sync-seeds":
		notImplemented()
	default:
		usage()
	}

	return os.Args[1]
}

func notImplemented() {
	fmt.Println("Method not implemented")
	os.Exit(1)
}

func usage() {
	fmt.Printf(`usage: %s <command>
	
Commands:
	* detect
	* aggregate
`, os.Args[0])
	os.Exit(1)
}
