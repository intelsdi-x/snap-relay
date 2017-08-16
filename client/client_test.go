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
	"testing"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/client"
	"github.com/intelsdi-x/snap/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {

	cmd := exec.Command("../build/darwin/x86_64/snap-relay", "--stand-alone")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer cmd.Process.Kill()
	time.Sleep(2 * time.Second)

	// fmt.Printf("****pid: %v\n", cmd.Process.Pid)
	// fmt.Printf("****status: %v\n", cmd.ProcessState)
	// fmt.Printf("****stderr: %v\n", cmd.Stderr)
	// fmt.Printf("****stdOut: %v\n", cmd.Stdout)

	resp, err := http.Get("http://localhost:8182")
	//	So(err, ShouldBeNil)
	//	So(resp, ShouldNotBeNil)
	if err != nil {
		log.Fatal(err)
	}
	if resp == nil {
		log.Fatal("resp should equal nil, actually equals ", resp)
	}
	preamble := make(map[string]string)
	body, err := ioutil.ReadAll(resp.Body)
	//So(err, ShouldBeNil)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &preamble)

	c, err := client.NewStreamCollectorGrpcClient(
		preamble["ListenAddress"],
		5*time.Second,
		client.SecurityTLSOff(),
	)
	if err != nil {
		panic(err)
	}

	Convey("Test StreamMetrics", t, func() {
		requested_metrics := []core.Metric{
			plugin.MetricType{
				Namespace_: core.NewNamespace("relay", "collectd"),
			},
		}
		metricsChan, errChan, err := c.StreamMetrics("test-taskID", requested_metrics)
		So(err, ShouldBeNil)
		So(metricsChan, ShouldNotBeNil)
		So(errChan, ShouldNotBeNil)

		go func() {
			// 6123 is default port for graphite tcp
			conn, err := net.Dial("tcp", "localhost:6123")

			if err != nil {
				panic(err)
			}
			defer conn.Close()
			time.Sleep(2 * time.Second)

			//var timer int64 = 10
			//t := time.Duration(timer)
			//for now := time.Now(); time.Since(now) < t; {
			conn.Write([]byte("THIS IS A MESSAGE"))
			fmt.Println("Sent message: THIS IS A MESSAGE")
			//}
		}()

		t := time.NewTimer(time.Second * 10).C
		//for {
		fmt.Println("in select statement... ")
		select {
		case <-t:
			fmt.Printf("In timer case.... \n")
			break
		case receivedFromMetricsChan := <-metricsChan:
			fmt.Printf("In metricsChan case.... received: %v\n", receivedFromMetricsChan)
			So(receivedFromMetricsChan, ShouldContain, "THIS IS A MESSAGE")
			break
		}
		fmt.Println("after select statement... ")

		//}

		time.Sleep(2 * time.Second)

	})

	Convey("Test GetMetricTypes", t, func() {
		Convey("Collect String", func() {
			mt, err := c.GetMetricTypes(plugin.NewPluginConfigType())
			So(err, ShouldBeNil)
			So(len(mt), ShouldEqual, 2)
		})

	})

	Convey("Test GetConfigPolicy", t, func() {
		_, err := c.GetConfigPolicy()
		Convey("No error returned", func() {
			So(err, ShouldBeNil)
		})

	})

}
