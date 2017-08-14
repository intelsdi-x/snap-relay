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
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/intelsdi-x/snap-relay/graphite"
	"github.com/intelsdi-x/snap-relay/relay"
	"github.com/intelsdi-x/snap-relay/statsd"
)

const (
	pluginName    = "plugin-relay"
	pluginVersion = 1
)

func main() {
	plugin.Flags = append(plugin.Flags, graphite.GraphiteTCPListenPortFlag)
	plugin.Flags = append(plugin.Flags, graphite.GraphiteUDPListenPortFlag)
	plugin.Flags = append(plugin.Flags, statsd.StatsdTCPListenPortFlag)
	plugin.Flags = append(plugin.Flags, statsd.StatsdUDPListenPortFlag)
	plugin.StartStreamCollector(
		relay.New(
			graphite.TCPListenPortOption(&graphite.GraphiteTCPPort),
			graphite.UDPListenPortOption(&graphite.GraphiteUDPPort),
			statsd.TCPListenPortOption(&statsd.StatsdTCPPort),
			statsd.UDPListenPortOption(&statsd.StatsdUDPPort),
		),
		pluginName,
		pluginVersion,
	)
}
