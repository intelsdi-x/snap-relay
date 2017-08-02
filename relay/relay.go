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
	"context"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/intelsdi-x/snap-relay/graphite"
	"github.com/intelsdi-x/snap-relay/statsd"
)

const (
	Name    = "relay"
	Version = 1
)

type relayMetrics interface {
	Metrics(context.Context) chan *plugin.Metric
	Start() error
}

type relay struct {
	graphiteServer relayMetrics
	statsdServer   relayMetrics
}

type relayOption interface {
	Type() string
}

func New(opts ...relayOption) plugin.StreamCollector {
	var gOpts []graphite.Option
	var sOpts []statsd.Option
	for _, x := range opts {

		switch t := x.(type) {
		case graphite.Option:
			gOpts = append(gOpts, t)
		case statsd.Option:
			sOpts = append(sOpts, t)
		}
	}

	return &relay{
		graphiteServer: graphite.NewGraphite(gOpts...),
		//statsdServer:   statsd.NewStatsd(sOpts...),
	}
}

func (r *relay) StreamMetrics(ctx context.Context, metrics_in chan []plugin.Metric, metrics_out chan []plugin.Metric, err chan string) error {
	log.SetLevel(log.Level(plugin.LogLevel))
	for metrics := range metrics_in {
		log.Debug("starting StreamMetrics")
		graphiteDispatchStarted := false
		statsdDispatchStarted := false
		r.graphiteServer.Start()
		r.statsdServer.Start()
		log.WithFields(
			log.Fields{
				"len(metrics)": len(metrics),
			},
		).Debug("received metrics")
		for idx, metric := range metrics {
			log.WithFields(
				log.Fields{
					"metric": metric.Namespace.String(),
				},
			).Debug("received metrics")

			//assign port values if any passed in
			if metric.Namespace[len(metric.Namespace)-1].Value == "collectd" {
				if val, err := metric.Config.GetString("collectdPort"); err == nil {
					metrics[idx].Data = val
				}
			} else if metric.Namespace[len(metric.Namespace)-1].Value == "statsd" {
				if val, err := metric.Config.GetString("statsdPort"); err == nil {
					metrics[idx].Data = val
				}
			}

			if !graphiteDispatchStarted && strings.Contains(metric.Namespace.String(), "collectd") {
				graphiteDispatchStarted = true
				go dispatchMetrics(ctx, r.graphiteServer.Metrics(ctx), metrics_out)
			}
			if !statsdDispatchStarted && strings.Contains(metric.Namespace.String(), "statsd") {
				statsdDispatchStarted = true
				go dispatchMetrics(ctx, r.statsdServer.Metrics(ctx), metrics_out)
			}
		}
	}
	return nil
}

func dispatchMetrics(ctx context.Context, in chan *plugin.Metric, out chan []plugin.Metric) {
	for {
		select {
		case metric := <-in:
			log.WithFields(
				log.Fields{
					"metric": metric.Namespace.String(),
					"data":   metric.Data,
				},
			).Debug("dispatching metrics")
			out <- []plugin.Metric{*metric}
		case <-ctx.Done():
			return
		}
	}
}

/*
	GetMetricTypes() returns the metric types for testing

	GetMetricTypes() will be called when your plugin is loaded in order to populate the metric catalog(where snaps stores all
	available metrics).

	Config info is passed in. This config information would come from global config snap settings.

	The metrics returned will be advertised to users who list all the metrics and will become targetable by tasks.
*/
func (r *relay) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	mts := []plugin.Metric{}
	vals := []string{"collectd", "statsd"}
	for _, val := range vals {
		metric := plugin.Metric{
			Namespace: plugin.NewNamespace(Name, val),
			Version:   Version,
		}
		mts = append(mts, metric)
	}

	return mts, nil
}

// GetConfigPolicy() returns the config policy for your plugin
func (r *relay) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()

	policy.AddNewStringRule([]string{"relay", "collectd"},
		"collectdPort",
		false)

	policy.AddNewStringRule([]string{"relay", "statsd"},
		"statsdPort",
		false)

	return *policy, nil
}
