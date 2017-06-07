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

package graphite

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/intelsdi-x/snap-relay/protocol"
	"github.com/urfave/cli"
)

var (
	// ErrAlreadyStarted error
	ErrAlreadyStarted = errors.New("server already started")
	// TCPAddr
	TCPAddr = "127.0.0.1:6123"
	// TCPListAddrFlag for overriding the listen
	TCPListAddrFlag cli.StringFlag = cli.StringFlag{
		Name:        "graphite-tcp",
		Usage:       "graphite TCP listen address and port",
		Value:       TCPAddr,
		Destination: &TCPAddr,
	}
)

type graphite struct {
	udp       protocol.Receiver
	tcp       protocol.Receiver
	metrics   chan *plugin.Metric
	done      chan struct{}
	isStarted bool
}

func NewGraphite(opts ...Option) *graphite {
	graphite := &graphite{
		udp:       protocol.NewUDPListener(),
		tcp:       protocol.NewTCPListener(),
		metrics:   make(chan *plugin.Metric, 1000),
		done:      make(chan struct{}),
		isStarted: false,
	}

	for _, opt := range opts {
		opt(graphite)
	}
	return graphite
}

type Option func(g *graphite) Option

func (g *graphite) Metrics() chan *plugin.Metric {
	return g.metrics
}

func UDPConnectionOption(conn *net.UDPConn) Option {
	return func(g *graphite) Option {
		if g.isStarted {
			log.WithFields(log.Fields{
				"_block": "UDPConnectionOption",
			}).Warn("option cannot be set.  service already started")
			return UDPConnectionOption(nil)
		}
		g.udp = protocol.NewUDPListener(protocol.UDPConnectionOption(conn))
		return UDPConnectionOption(conn)
	}
}

func TCPListenAddrOption(addr *string) Option {
	return func(g *graphite) Option {
		if g.isStarted {
			log.WithFields(log.Fields{
				"_block": "TCPListenAddrOption",
			}).Warn("option cannot be set.  service already started")
			return TCPListenAddrOption(addr)
		}
		g.tcp = protocol.NewTCPListener(protocol.TCPListenAddrOption(addr))
		return TCPListenAddrOption(addr)
	}
}

func TCPListenerOption(conn *net.TCPListener) Option {
	return func(g *graphite) Option {
		if g.isStarted {
			log.WithFields(log.Fields{
				"_block": "TCPConnectionOption",
			}).Warn("option cannot be set.  service already started")
			return TCPListenerOption(nil)
		}
		g.tcp = protocol.NewTCPListener(protocol.TCPListenerOption(conn))
		return TCPListenerOption(conn)
	}
}

func (g *graphite) Start() error {
	log.Info("Starting graphite relay")
	if g.isStarted {
		return ErrAlreadyStarted
	}
	if err := g.udp.Start(); err != nil {
		return err
	}
	if err := g.tcp.Start(); err != nil {
		return err
	}
	g.isStarted = true
	go func() {
		for {
			select {
			case data := <-g.udp.Data():
				if metric := parse(string(data)); metric != nil {
					select {
					case g.metrics <- metric:
					default:
						log.WithFields(log.Fields{
							"transport":        "udp",
							"_block":           "graphite",
							"metric_namespace": strings.Join(metric.Namespace.Strings(), "/"),
						}).Warn("Dropping metric.  Channel is full")
					}
				}
			case data := <-g.tcp.Data():
				if metric := parse(string(data)); metric != nil {
					select {
					case g.metrics <- metric:
					default:
						log.WithFields(log.Fields{
							"transport":        "tcp",
							"_block":           "graphite",
							"metric_namespace": strings.Join(metric.Namespace.Strings(), "/"),
						}).Warn("Dropping metric.  Channel is full")
					}
				}
			case <-g.done:
				break
			}
		}
	}()
	return nil
}

func (g *graphite) stop() {
	g.udp.Stop()
	g.tcp.Stop()
	close(g.done)
}

func parse(data string) *plugin.Metric {
	data = strings.Trim(data, "\r")
	line := strings.Split(data, " ")
	if len(line) != 3 {
		log.WithFields(log.Fields{
			"data": data,
		}).Warnln("unable to parse graphite data")
		return nil
	}
	ns := plugin.NewNamespace("collectd")
	ns = ns.AddStaticElements(strings.Split(line[0], ".")...)
	epoch, err := strconv.ParseInt(line[2], 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"_block": "toMetric",
			"data":   epoch,
			"error":  err.Error(),
		}).Error("failed to parse timestamp")
		return nil
	}
	timestamp := time.Unix(epoch, 0)
	return &plugin.Metric{
		Namespace: ns,
		Timestamp: timestamp,
		Data:      line[1],
	}
}
