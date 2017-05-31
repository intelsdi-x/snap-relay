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

package relay

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/intelsdi-x/snap-relay/graphite"
)

//TODO rename :)
type relayMetrics interface {
	Metrics() chan *plugin.Metric
	Start() error
}

//TODO consider not exporting
type Relay struct {
	graphiteServer relayMetrics
}

func New(opts ...graphite.Option) plugin.StreamCollector {
	return &Relay{
		graphiteServer: graphite.NewGraphite(opts...),
	}
}

func (r *Relay) StreamMetrics(metrics_in chan []plugin.Metric, metrics_out chan []plugin.Metric, err chan string) error {
	//TODO(JC) get log level from parent
	log.SetLevel(log.DebugLevel)
	// - listen on metrics_in
	//   - start server (collectd, statsd, etc) if requested
	//      - kick off go func to drain metrics from the server
	//		- emit messages received from server.Metrics() on metrics_out
	for metrics := range metrics_in {
		log.Debug("starting StreamMetrics")
		graphiteDispatchStarted := false
		r.graphiteServer.Start()
		log.WithFields(
			log.Fields{
				"len(metrics)": len(metrics),
			},
		).Debug("received metrics")
		for _, metric := range metrics {
			log.WithFields(
				log.Fields{
					"metric": metric.Namespace.String(),
				},
			).Debug("received metrics")
			if !graphiteDispatchStarted && strings.Contains(metric.Namespace.String(), "collectd") {
				graphiteDispatchStarted = true
				go dispatchMetrics(r.graphiteServer.Metrics(), metrics_out)
			}
		}
	}
	return nil
}

//TODO consider refactoring to make out chan (the lib) have a chan of plugin.Metric instead of the array
func dispatchMetrics(in chan *plugin.Metric, out chan []plugin.Metric) {
	for metric := range in {
		log.WithFields(
			log.Fields{
				"metric": metric.Namespace.String(),
				"data":   metric.Data,
			},
		).Debug("dispatching metrics")
		//TODO fix this weird derefence
		out <- []plugin.Metric{*metric}
	}
}

/*
	GetMetricTypes() returns the metric types for testing

	GetMetricTypes() will be called when your plugin is loaded in order to populate the metric catalog(where snaps stores all
	available metrics).

	Config info is passed in. This config information would come from global config snap settings.

	The metrics returned will be advertised to users who list all the metrics and will become targetable by tasks.
*/
func (r *Relay) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	mts := []plugin.Metric{}
	vals := []string{"collectd", "statsd"}
	for _, val := range vals {
		metric := plugin.Metric{
			Namespace: plugin.NewNamespace("relay", val),
			Version:   1,
		}
		mts = append(mts, metric)
	}

	return mts, nil
}

// GetConfigPolicy() returns the config policy for your plugin
func (r *Relay) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()

	policy.AddNewStringRule([]string{"relay", "collectd"},
		"graphite",
		false)

	policy.AddNewStringRule([]string{"relay", "statsd"},
		"statsd",
		false)

	return *policy, nil
}
