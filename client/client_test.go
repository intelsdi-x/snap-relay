// +build medium

/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2017 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/client"
	"github.com/intelsdi-x/snap/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	// Start the plugin in stand-alone mode
	cmd := exec.Command("../build/darwin/x86_64/snap-relay", "--stand-alone")
	err := cmd.Start()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Fatal(err, "\nHINT: Binary not found at specified location. Try running 'make' in the root directory then retry this test.\n")
		}
		log.Fatal(err)
	}

	defer cmd.Process.Kill()
	time.Sleep(2 * time.Second)

	resp, err := http.Get("http://localhost:8182")
	if err != nil {
		log.Fatal(err)
	}
	if resp == nil {
		log.Fatal("response from http.Get should not equal nil.")
	}

	// Unmarshal preamble to get ListenAddress
	preamble := make(map[string]string)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &preamble)

	// Create client
	c, err := client.NewStreamCollectorGrpcClient(
		preamble["ListenAddress"],
		5*time.Second,
		client.SecurityTLSOff(),
	)
	if err != nil {
		panic(err)
	}

	Convey("Test StreamMetrics", t, func() {
		// Start streaming the requested metrics
		requestedMetrics := []core.Metric{
			plugin.MetricType{
				Namespace_: core.NewNamespace("relay", "collectd"),
			},
		}
		metricsChan, errChan, err := c.StreamMetrics("test-taskID", requestedMetrics)
		So(err, ShouldBeNil)
		So(metricsChan, ShouldNotBeNil)
		So(errChan, ShouldNotBeNil)

		// Open default port and send data
		go func() {
			// 6123 is default port for graphite tcp
			conn, err := net.Dial("tcp", "localhost:6123")
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			time.Sleep(2 * time.Second)
			conn.Write([]byte("test.first 13 1502916834\n"))
		}()

		// Listen for sent metric. Exit if metric not received in 10s
		t := time.NewTimer(time.Second * 10).C
		select {
		case <-t:
			fmt.Printf("Metric not received in 10 seconds, exiting. \n")
			So(true, ShouldBeFalse)
			break
		case receivedFromMetricsChan := <-metricsChan:
			So(len(receivedFromMetricsChan), ShouldEqual, 1)
			break
		}

	})
}
