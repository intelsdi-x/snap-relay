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
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	// string "test -run"
	fmt.Printf("%s\n", os.Args)
	// using os.Exec (or something similiar) to start the server with stand-alone flag
	// as an alt to the go func
	// The rest of the test should work as intended
	//os.Args = []string{os.Args[0], "--stand-alone"}

	cmd := exec.Command(os.Args[0], "--stand-alone")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	fmt.Printf("****stderr: %v\n", cmd.Stderr)
	fmt.Printf("****stdOut: %v\n", cmd.Stdout)

	// go func() {
	// 	retcode := plugin.StartStreamCollector(New(), "snap-relay", 1) //starts relay
	// 	fmt.Printf("retcode=%v\n", retcode)
	// }()

	Convey("Test StreamMetrics", t, func() {
		So(nil, ShouldBeNil)

		_, err := http.Get("http://localhost:8182")
		So(err, ShouldBeNil)

		// c, err := client.NewStreamCollectorGrpcClient(
		// 	"localhost:8182",
		// 	5*time.Second,
		// 	client.SecurityTLSOff(),
		// )
		// if err != nil {
		// 	panic(err)
		// }

		// requested_metrics := []core.Metric{
		// 	ctlplugin.MetricType{
		// 		Namespace_: core.NewNamespace("relay", "collectd"),
		// 	},
		// }
		// metricsChan, errChan, err := c.StreamMetrics(requested_metrics)
		// So(err, ShouldBeNil)
		// So(metricsChan, ShouldNotBeNil)
		// So(errChan, ShouldNotBeNil)
		// cfg := cdata.NewNode()
		// cfg.AddItem("MaxCollectDuration", ctypes.ConfigValueInt{Value: 5000000000})
		// cfg.AddItem("MaxMetricsBuffer", ctypes.ConfigValueInt{Value: 2})

		// requested_metrics := []core.Metric{
		// 	ctlplugin.MetricType{
		// 		Namespace_: core.NewNamespace("animal", "cat"),
		// 		Config_:    cfg,
		// 	},
		// }

		//  rq <- requested_metrics

		//Need to use this (the one actually defined in relay.go)
		//func (r *relay) StreamMetrics(ctx context.Context, metrics_in chan []plugin.Metric, metrics_out chan []plugin.Metric, err chan string) error {
		//instead of one below:

		// var d time.Duration = 5

		// ctx, _ := context.WithTimeout(context.Background(), d)

		// chanIn := make(chan []plugin.Metric)
		// chanOut := make(chan []plugin.Metric)
		// chanErr := make(chan string)

		// r := relay{}
		// err = r.StreamMetrics(ctx, chanIn, chanOut, chanErr)
		// if err != nil {
		//  panic(err)
		// }

		// x, y, z := <-chanIn, <-chanOut, <-chanErr
		// fmt.Printf("x: %v \ny: %v \nz: %v", x, y, z)

		// metricsOut, errOut, err := c.StreamMetrics(requested_metrics)
		// So(metricsOut, ShouldNotBeEmpty)
		// So(errOut, ShouldBeNil)
		// So(err, ShouldBeNil)

	})

	// Convey("Test GetMetricTypes", t, func() {
	// 	r := relay.New()

	// 	Convey("Collect String", func() {
	// 		mt, err := r.GetMetricTypes(nil)
	// 		So(err, ShouldBeNil)
	// 		So(len(mt), ShouldEqual, 2)
	// 	})

	// })

	// Convey("Test GetConfigPolicy", t, func() {
	// 	r := relay.New()
	// 	_, err := r.GetConfigPolicy()

	// 	Convey("No error returned", func() {
	// 		So(err, ShouldBeNil)
	// 	})

	// })

}
