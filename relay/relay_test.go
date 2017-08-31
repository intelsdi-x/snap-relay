// +build small

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
	"net"
	"testing"

	"github.com/intelsdi-x/snap-relay/graphite"
	"github.com/intelsdi-x/snap-relay/statsd"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRelay(t *testing.T) {
	r := relay{}

	Convey("Test GetMetricTypes", t, func() {
		Convey("Collect String", func() {
			mt, err := r.GetMetricTypes(nil)
			So(err, ShouldBeNil)
			So(len(mt), ShouldEqual, 2)
		})
	})

	Convey("Test GetConfigPolicy", t, func() {
		_, err := r.GetConfigPolicy()
		Convey("No error returned", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Test New", t, func() {
		udpAddr, err := net.ResolveUDPAddr("udp", "localhost:0")
		So(err, ShouldBeNil)
		So(udpAddr, ShouldNotBeNil)
		udpConn, err := net.ListenUDP("udp", udpAddr)
		So(err, ShouldBeNil)
		So(udpConn, ShouldNotBeNil)
		//create graphite option for test
		myGraphiteUDPOption := graphite.UDPConnectionOption(udpConn)
		//create statsd option for test
		myStatsdUDPOption := statsd.UDPConnectionOption(udpConn)

		newRelayTCP := New(myGraphiteUDPOption, myStatsdUDPOption)
		So(newRelayTCP, ShouldNotBeNil)
	})

	// Test for StreamMetrics can be found in client/client_test.go as a medium test
}
