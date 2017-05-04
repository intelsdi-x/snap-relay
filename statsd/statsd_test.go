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
	"net"
	"testing"

	"time"

	log "github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStatsd(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	Convey("Setup statsd server", t, func() {
		udpAddr, err := net.ResolveUDPAddr("udp", "localhost:0")
		So(err, ShouldBeNil)
		So(udpAddr, ShouldNotBeNil)
		udpConn, err := net.ListenUDP("udp", udpAddr)
		So(err, ShouldBeNil)
		So(udpConn, ShouldNotBeNil)
		tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		So(err, ShouldBeNil)
		So(tcpAddr, ShouldNotBeNil)
		listen, err := net.ListenTCP("tcp", tcpAddr)
		So(err, ShouldBeNil)
		statsd := NewStatsd(UDPConnectionOption(udpConn), TCPListenerOption(listen))
		So(statsd, ShouldNotBeNil)
		err = statsd.Start()
		So(err, ShouldBeNil)
		// create tcpClient
		tcpAddr, err = net.ResolveTCPAddr("tcp", listen.Addr().String())
		So(err, ShouldBeNil)
		So(tcpAddr, ShouldNotBeNil)
		tcpClient, err := net.DialTCP("tcp", nil, tcpAddr)
		Convey("starts already started server", func() {
			err = statsd.Start()
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, ErrAlreadyStarted)
		})
		Convey("receives invalid data over UDP", func() {
			msgs := []string{"hello\n", "foo\n", "bar\n"}
			for _, msg := range msgs {
				_, err := udpConn.WriteTo([]byte(msg), udpConn.LocalAddr())
				So(err, ShouldBeNil)
			}
			time.Sleep(100 * time.Millisecond)
			So(len(statsd.metrics), ShouldEqual, 0)
		})
		Convey("receives invalid data over TCP", func() {
			msgs := []string{"hello\n", "foo\n", "bar\n"}
			for _, msg := range msgs {
				_, err := tcpClient.Write([]byte(msg))
				So(err, ShouldBeNil)
			}
			time.Sleep(100 * time.Millisecond)
			So(len(statsd.metrics), ShouldEqual, 0)
		})
		Convey("sends valid UDP data", func() {
			msgs := []string{
				"foo.bar:7|c\n",
				"foo.bar:8|g\n",
				"foo.bar:97|s\n",
			}
			for _, msg := range msgs {
				_, err := udpConn.WriteTo([]byte(msg), udpConn.LocalAddr())
				So(err, ShouldBeNil)
			}
			time.Sleep(100 * time.Millisecond)
			So(len(statsd.metrics), ShouldEqual, 3)
			Convey("sends valid TCP data", func() {
				for _, msg := range msgs {
					_, err := tcpClient.Write([]byte(msg))
					So(err, ShouldBeNil)
				}
				time.Sleep(100 * time.Millisecond)
				So(len(statsd.metrics), ShouldEqual, 6)
			})
		})
	})

}
