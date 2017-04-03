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
		conn, err := net.ListenUDP("udp", udpAddr)
		So(err, ShouldBeNil)
		graphite := NewGraphite(UDPConnectionOption(conn))
		So(graphite, ShouldNotBeNil)
		err = graphite.Start()
		So(err, ShouldBeNil)
		Convey("starts already started server", func() {
			err = graphite.Start()
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, ErrAlreadyStaretd)
		})
		Convey("receives invalid data", func() {
			msgs := []string{"hello\n", "foo\n", "bar\n"}
			for _, msg := range msgs {
				_, err := conn.WriteTo([]byte(msg), conn.LocalAddr())
				So(err, ShouldBeNil)
			}
			time.Sleep(100 * time.Millisecond)
			So(len(graphite.metrics), ShouldEqual, 0)
		})
		Convey("send valid data", func() {
			msgs := []string{
				"myhost_example_com.cpu-2.cpu-idle 98.6103 1329168255\n",
				"myhost_example_com.cpu-2.cpu-nice 0 1329168255\n",
				"myhost_example_com.cpu-2.cpu-user 0.800076 1329168255\n",
			}
			for _, msg := range msgs {
				_, err := conn.WriteTo([]byte(msg), conn.LocalAddr())
				So(err, ShouldBeNil)
			}
			time.Sleep(100 * time.Millisecond)
			So(len(graphite.metrics), ShouldEqual, 3)
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

}
