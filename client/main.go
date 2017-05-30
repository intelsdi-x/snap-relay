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
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/client"
	"github.com/intelsdi-x/snap/core"

	log "github.com/Sirupsen/logrus"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	log.SetLevel(log.DebugLevel)
	if len(os.Args) != 2 {
		log.Fatal("this program expects a single arg for the host and port example: localhost:1234")
	}
	c, err := client.NewStreamCollectorGrpcClient(
		os.Args[1],
		5*time.Second,
		nil,
		false,
	)
	if err != nil {
		panic(err)
	}
	metricsOut, errOut, err := c.StreamMetrics([]core.Metric{plugin.MetricType{Namespace_: core.NewNamespace("relay", "collectd")}})
	if err != nil {
		panic(err)
	}
	go func() {
		for metrics := range metricsOut {
			for _, metric := range metrics {
				log.WithFields(
					log.Fields{
						"metric": metric,
					},
				).Debug("received metric")
			}
		}
	}()

	go func() {
		for err := range errOut {
			log.WithFields(
				log.Fields{
					"error": err,
				},
			).Error("error received")
		}
	}()
	<-done
}