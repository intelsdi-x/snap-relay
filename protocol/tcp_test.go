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
		select {
		case <-server.done:
			time.Sleep(100 * time.Millisecond)
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out waiting for TCP server to stop")
		}
	})
}
