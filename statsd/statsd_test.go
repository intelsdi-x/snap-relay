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
	"net"
	"testing"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	. "github.com/smartystreets/goconvey/convey"
)

type testContext struct{ context.Context }

func (tc testContext) Deadline() (deadline time.Time, ok bool) {
	return time.Now(), true
}
func (tc testContext) Done() <-chan struct{} {
	return nil
}
func (tc testContext) Err() error {
	return nil
}
func (tc testContext) Value(key interface{}) interface{} {
	return nil
}

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
		udpOption := UDPConnectionOption(udpConn)
		So(udpOption.Type(), ShouldEqual, "statsd")
		port := 1234
		statsd := NewStatsd(UDPListenPortOption(&port), udpOption, TCPListenPortOption(&port), TCPListenerOption(listen))
		So(statsd, ShouldNotBeNil)
		err = statsd.Start()
		So(err, ShouldBeNil)
		So(statsd.isStarted, ShouldBeTrue)

		MyOtherUDPOption := UDPConnectionOption(udpConn)(statsd)
		So(MyOtherUDPOption, ShouldEqual, UDPConnectionOption(nil))
		myOtherTCPOption := TCPListenerOption(listen)(statsd)
		So(myOtherTCPOption, ShouldEqual, TCPListenerOption(nil))
		myOtherTCPPortOption := TCPListenPortOption(&port)(statsd)
		So(myOtherTCPPortOption, ShouldEqual, TCPListenPortOption(nil))
		MyOtherUDPPortOption := UDPListenPortOption(&port)(statsd)
		So(MyOtherUDPPortOption, ShouldEqual, UDPListenPortOption(nil))

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
			So(len(statsd.Metrics(testContext{})), ShouldEqual, 0)
		})
		Convey("receives invalid data over TCP", func() {
			msgs := []string{"hello\n", "foo\n", "bar\n"}
			for _, msg := range msgs {
				_, err := tcpClient.Write([]byte(msg))
				So(err, ShouldBeNil)
			}
			time.Sleep(100 * time.Millisecond)
			So(len(statsd.metrics), ShouldEqual, 0)
			So(len(statsd.Metrics(testContext{})), ShouldEqual, 0)
		})
		Convey("sends valid UDP data", func() {
			myCh := make(chan *plugin.Metric, 1000)
			statsd.channelMgr.Add(myCh)

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
			So(len(myCh), ShouldEqual, 3)

			Convey("sends valid TCP data", func() {
				for _, msg := range msgs {
					_, err := tcpClient.Write([]byte(msg))
					So(err, ShouldBeNil)
				}
				time.Sleep(100 * time.Millisecond)
				So(len(myCh), ShouldEqual, 6)
			})
		})
		statsd.stop()
		time.Sleep(100 * time.Millisecond)
		select {
		case done := <-statsd.done:
			So(done, ShouldNotBeNil)
		case <-time.After(100 * time.Millisecond):
			t.Fail()
		}
	})

}
