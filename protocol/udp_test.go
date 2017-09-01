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

func TestUDPListen(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	Convey("Setup UDP server and client", t, func(c C) {
		udpAddr, err := net.ResolveUDPAddr("udp", "localhost:0")
		So(err, ShouldBeNil)
		So(udpAddr, ShouldNotBeNil)
		conn, err := net.ListenUDP("udp", udpAddr)
		So(err, ShouldBeNil)
		server := NewUDPListener(UDPConnectionOption(conn))
		err = server.Start()
		c.So(err, ShouldBeNil)
		Convey("Send/receive udp messages (with newline)", func() {
			msgs := []string{"hello\n", "foo\n", "bar\n"}
			clientConn := server.conn
			for _, msg := range msgs {
				_, err := clientConn.WriteTo([]byte(msg), clientConn.LocalAddr())
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
					_, err := clientConn.WriteTo([]byte(msg), clientConn.LocalAddr())
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
		select {
		case <-server.done:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out waiting for UDP server to stop")
		}
	})

	Convey("Setup failing UDP server", t, func(c C) {
		//good ResolveUDPAddr
		udpAddr, err := net.ResolveUDPAddr("udp", "localhost:0")
		So(err, ShouldBeNil)
		So(udpAddr, ShouldNotBeNil)

		//bad server.Start
		BadConn, err := net.ListenUDP("udppdu", udpAddr)
		So(err, ShouldNotBeNil)

		//start server with badConn
		listenPort := 5
		server := NewUDPListener(UDPConnectionOption(BadConn), UDPListenPortOption(&listenPort))
		err = server.Start()
		c.So(err, ShouldNotBeNil)

		//start server with no conn
		server = NewUDPListener()
		err = server.Start()
		c.So(err, ShouldBeNil)
	})

}
