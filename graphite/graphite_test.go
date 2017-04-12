package graphite

import (
	"net"
	"testing"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGraphite(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	Convey("Setup graphite server", t, func() {
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
		graphite := NewGraphite(UDPConnectionOption(udpConn), TCPListenerOption(listen))
		So(graphite, ShouldNotBeNil)
		err = graphite.Start()
		So(err, ShouldBeNil)
		// create tcpClient
		tcpAddr, err = net.ResolveTCPAddr("tcp", listen.Addr().String())
		So(err, ShouldBeNil)
		So(tcpAddr, ShouldNotBeNil)
		tcpClient, err := net.DialTCP("tcp", nil, tcpAddr)
		Convey("starts already started server", func() {
			err = graphite.Start()
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
			So(len(graphite.metrics), ShouldEqual, 0)
		})
		Convey("receives invalid data over TCP", func() {
			msgs := []string{"hello\n", "foo\n", "bar\n"}
			for _, msg := range msgs {
				_, err := tcpClient.Write([]byte(msg))
				So(err, ShouldBeNil)
			}
			time.Sleep(100 * time.Millisecond)
			So(len(graphite.metrics), ShouldEqual, 0)
		})
		Convey("sends valid UDP data", func() {
			msgs := []string{
				"myhost_example_com.cpu-2.cpu-idle 98.6103 1329168255\n",
				"myhost_example_com.cpu-2.cpu-nice 0 1329168255\n",
				"myhost_example_com.cpu-2.cpu-user 0.800076 1329168255\n",
			}
			for _, msg := range msgs {
				_, err := udpConn.WriteTo([]byte(msg), udpConn.LocalAddr())
				So(err, ShouldBeNil)
			}
			time.Sleep(100 * time.Millisecond)
			So(len(graphite.metrics), ShouldEqual, 3)
			Convey("sends valid TCP data", func() {
				for _, msg := range msgs {
					_, err := tcpClient.Write([]byte(msg))
					So(err, ShouldBeNil)
				}
				time.Sleep(100 * time.Millisecond)
				So(len(graphite.metrics), ShouldEqual, 6)
				Convey("Reads metrics from the buffer", func() {
					So(
						<-graphite.metrics,
						ShouldResemble,
						&plugin.Metric{
							Namespace: plugin.NewNamespace(
								"collectd", "myhost_example_com", "cpu-2", "cpu-idle"),
							Timestamp: time.Unix(1329168255, 0),
						},
					)
				})
			})
		})
	})

}
