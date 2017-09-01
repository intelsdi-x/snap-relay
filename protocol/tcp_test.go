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

package protocol

import (
	"net"
	"testing"

	"time"

	log "github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTCPListen(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	Convey("Setup TCP server and client", t, func() {
		tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		So(err, ShouldBeNil)
		So(tcpAddr, ShouldNotBeNil)
		conn, err := net.ListenTCP("tcp", tcpAddr)
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)
		server := NewTCPListener(TCPListenerOption(conn))
		So(server, ShouldNotBeNil)
		err = server.Start()
		time.Sleep(100 * time.Millisecond)
		So(err, ShouldBeNil)
		Convey("Send/receive messages (with newline)", func() {
			msgs := []string{"hello\n", "foo\n", "bar\n"}
			tcpAddr, err := net.ResolveTCPAddr("tcp", conn.Addr().String())
			So(err, ShouldBeNil)
			So(tcpAddr, ShouldNotBeNil)
			clientConn, err := net.DialTCP("tcp", nil, tcpAddr)
			So(err, ShouldBeNil)
			for _, msg := range msgs {
				_, err := clientConn.Write([]byte(msg))
				So(err, ShouldBeNil)
			}
			for _, msg := range msgs {
				select {
				case data := <-server.data:
					So(data, ShouldResemble, []byte(msg[:len(msg)-1]))
				case <-time.After(time.Millisecond * 200):
					t.Fatalf("timed out while reading sent data")
				}
			}
			So(len(server.data), ShouldEqual, 0)
			So(len(server.Data()), ShouldEqual, 0)

			Convey("without newline", func() {
				msgs := []string{"hello", "foo", "bar"}
				for _, msg := range msgs {
					_, err := clientConn.Write([]byte(msg))
					So(err, ShouldBeNil)
				}
				select {
				case <-server.data:
					t.Fatalf("messages without a newline should be ignored")
				case <-time.After(time.Millisecond * 100):
					break
				}
			})
		})
		server.Stop()
		reachedDone := false
		select {
		case <-server.done:
			time.Sleep(100 * time.Millisecond)
			reachedDone = true
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out waiting for TCP server to stop")
		}
		So(reachedDone, ShouldBeTrue)
	})

	Convey("Setup failing TCP server", t, func(c C) {
		//good ResolveTCPAddr
		tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		So(err, ShouldBeNil)
		So(tcpAddr, ShouldNotBeNil)

		//bad server.Start
		BadConn, err := net.ListenTCP("tcppct", tcpAddr)
		So(err, ShouldNotBeNil)

		//start server with badConn
		listenPort := 5
		server := NewTCPListener(TCPListenerOption(BadConn), TCPListenPortOption(&listenPort))
		err = server.Start()
		c.So(err, ShouldNotBeNil)

		//start server with no conn
		server = NewTCPListener()
		err = server.Start()
		c.So(err, ShouldBeNil)
	})

}
