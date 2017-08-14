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
package statsd

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/intelsdi-x/snap-relay/protocol"
	"github.com/intelsdi-x/snap-relay/util"
	"github.com/urfave/cli"
)

var (
	ErrAlreadyStarted = errors.New("server already started")
	// StatsdTCPPort
	StatsdTCPPort = 6125
	// StatsdTCPListenPortFlag for overriding the listen address
	StatsdTCPListenPortFlag cli.IntFlag = cli.IntFlag{
		Name:        "statsd-tcp-port",
		Usage:       "statsd TCP listen port",
		Value:       StatsdTCPPort,
		Destination: &StatsdTCPPort,
	}
	// StatsdUDPAddr
	StatsdUDPPort = 6126
	// StatsdUDPListenAddrFlag for overriding the listen address
	StatsdUDPListenPortFlag cli.IntFlag = cli.IntFlag{
		Name:        "Statsd-udp-port",
		Usage:       "statsd UDP listen port",
		Value:       StatsdUDPPort,
		Destination: &StatsdUDPPort,
	}
)

type statsd struct {
	udp        protocol.Receiver
	tcp        protocol.Receiver
	metrics    chan *plugin.Metric
	channelMgr util.ChannelManager
	done       chan struct{}
	isStarted  bool
}

func NewStatsd(opts ...Option) *statsd {
	statsd := &statsd{
		udp:        protocol.NewUDPListener(),
		tcp:        protocol.NewTCPListener(),
		metrics:    make(chan *plugin.Metric, 1000),
		done:       make(chan struct{}),
		isStarted:  false,
		channelMgr: util.NewChannelMgr(),
	}

	for _, opt := range opts {
		opt(statsd)
	}
	return statsd
}

type Option func(sd *statsd) Option

func (o Option) Type() string {
	return "statsd"
}

func UDPConnectionOption(conn *net.UDPConn) Option {
	return func(sd *statsd) Option {
		if sd.isStarted {
			log.WithFields(log.Fields{
				"_block": "UDPConnectionOption",
			}).Warn("option cannot be set.  service already started")
			return UDPConnectionOption(nil)
		}
		sd.udp = protocol.NewUDPListener(protocol.UDPConnectionOption(conn))
		return UDPConnectionOption(conn)
	}
}

func UDPListenPortOption(port *int) Option {
	return func(sd *statsd) Option {
		if sd.isStarted {
			log.WithFields(log.Fields{
				"_block": "UDPListenPortOption",
				"detail": "service already started",
			}).Warn("option cannot be set")
			return UDPListenPortOption(port)
		}
		sd.udp = protocol.NewUDPListener(protocol.UDPListenPortOption(port))
		return UDPListenPortOption(port)
	}
}

func TCPListenPortOption(port *int) Option {
	return func(sd *statsd) Option {
		if sd.isStarted {
			log.WithFields(log.Fields{
				"_block": "TCPListenPortOption",
				"detail": "service already started",
			}).Warn("option cannot be set")
			return TCPListenPortOption(port)
		}
		sd.tcp = protocol.NewTCPListener(protocol.TCPListenPortOption(port))
		return TCPListenPortOption(port)
	}
}

func TCPListenerOption(conn *net.TCPListener) Option {
	return func(sd *statsd) Option {
		if sd.isStarted {
			log.WithFields(log.Fields{
				"_block": "TCPConnectionOption",
			}).Warn("option cannot be set.  service already started")
			return TCPListenerOption(nil)
		}
		sd.tcp = protocol.NewTCPListener(protocol.TCPListenerOption(conn))
		return TCPListenerOption(conn)
	}
}

func (sd *statsd) Start() error {
	if sd.isStarted {
		return ErrAlreadyStarted
	}
	log.Info("Starting statsd relay")
	if err := sd.udp.Start(); err != nil {
		return err
	}
	if err := sd.tcp.Start(); err != nil {
		return err
	}
	sd.isStarted = true
	go func() {
		for {
			select {
			case data := <-sd.udp.Data():
				lines := strings.Split(string(data), "\n")
				for _, line := range lines {
					if metric := parseData(string(line)); metric != nil {
						select {
						case sd.metrics <- metric:
						default:
							log.WithFields(log.Fields{
								"transport":        "udp",
								"_block":           "statsd",
								"metric_namespace": strings.Join(metric.Namespace.Strings(), "/"),
							}).Warn("dropping metric.  Channel is full")
						}
					}
				}
			case data := <-sd.tcp.Data():
				lines := strings.Split(string(data), "\n")
				for _, line := range lines {
					if metric := parseData(string(line)); metric != nil {
						select {
						case sd.metrics <- metric:
						default:
							log.WithFields(log.Fields{
								"transport":        "tcp",
								"_block":           "statsd",
								"metric_namespace": strings.Join(metric.Namespace.Strings(), "/"),
							}).Warn("dropping metric.  Channel is full")
						}
					}
				}
			case <-sd.done:
				break
			}
		}
	}()

	// routine that dispatches statsd metrics to all available streams
	go func() {
		for {
			select {
			case m := <-sd.metrics:
				log.Debugf("dispatching metrics to %v streams", sd.channelMgr.Count())
				sd.channelMgr.DispatchMetric(m)
			case <-sd.done:
				return
			}
		}
	}()
	return nil
}

func (sd *statsd) Metrics(ctx context.Context) chan *plugin.Metric {
	mchan := make(chan *plugin.Metric, 1000)
	sd.channelMgr.Add(mchan)
	go func() {
		select {
		case <-ctx.Done():
			sd.channelMgr.Remove(mchan)
		}
	}()
	return mchan
}

func (sd *statsd) stop() {
	sd.udp.Stop()
	sd.tcp.Stop()
	close(sd.done)
}

func parseMetricType(t string) string {
	switch t {
	case "c":
		return "counter"
	case "g":
		return "gauge"
	case "s":
		return "set"
	case "ms":
		return "timer"
	default:
		return t
	}

}

func parseData(data string) *plugin.Metric {
	tags := map[string]string{}
	lineElems := strings.Split(data, "|")
	if len(lineElems) > 2 {
		log.WithFields(log.Fields{
			"_block":        "parseData",
			"received_data": data,
			"expected_data": "<metric>:<value>|<type>",
		}).Error("invalid metric line")
		return nil
	}
	tags["data_type"] = parseMetricType(lineElems[1])
	metricElements := strings.Split(lineElems[0], ":")
	if len(metricElements) != 2 {
		log.WithFields(log.Fields{
			"_block":        "parseData",
			"received_data": lineElems[0],
			"expected_data": "<metric>:<value>",
		}).Error("invalid data in metric line")
		return nil
	}
	ns := plugin.NewNamespace("statsd")
	ns = ns.AddStaticElements(strings.Split(metricElements[0], ".")...)
	value, err := strconv.ParseInt(metricElements[1], 10, 64)

	if err != nil {
		log.WithFields(log.Fields{
			"_block":    "parseData",
			"namespace": ns,
			"data":      value,
			"error":     err.Error(),
		}).Error("failed to parse data")
		return nil
	}
	return &plugin.Metric{
		Namespace: ns,
		Data:      value,
		Tags:      tags,
		Timestamp: time.Now(),
	}
}
