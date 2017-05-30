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
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type udpListener struct {
	data chan []byte
	conn *net.UDPConn
	done chan struct{}
}

func NewUDPListener(opts ...option) *udpListener {
	listener := &udpListener{
		data: make(chan []byte, 100),
		done: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(listener)
	}
	return listener
}

type option func(u *udpListener) option

func UDPConnectionOption(conn *net.UDPConn) option {
	return func(u *udpListener) option {
		prev := u.conn
		u.conn = conn
		return UDPConnectionOption(prev)
	}
}

func (u *udpListener) listen() error {
	if u.conn == nil {
		udpAddr, err := net.ResolveUDPAddr("udp", "localhost:0")
		if err != nil {
			return err
		}
		u.conn, err = net.ListenUDP("udp", udpAddr)
		if err != nil {
			return err
		}
		log.WithFields(
			log.Fields{
				"addr": u.conn.LocalAddr().String(),
			},
		).Debug("udp listening started")
	}

	return nil
}

func (u *udpListener) Data() chan []byte {
	return u.data
}

func (u *udpListener) Stop() {
	close(u.done)
}

func (u *udpListener) Start() error {
	var buf [65535]byte
	var data *bytes.Buffer
	if err := u.listen(); err != nil {
		return err
	}
	log.WithField("addr", u.conn.LocalAddr().String()).Debug("started UDP listener")

	go func() {
		for {
			select {
			case <-u.done:
				u.conn.Close()
				log.WithField("addr", u.conn.LocalAddr().String()).Debug("stopped UDP listener")
				break
			default:
				rlen, peer, err := u.conn.ReadFromUDP(buf[:])
				if err != nil {
					if !strings.Contains(err.Error(), "use of closed network connection") {
						log.WithFields(log.Fields{
							"peer": peer.String(),
							"err":  err.Error(),
						}).Error("error reading from UDP")
					}
					close(u.done)
				}

				data = bytes.NewBuffer(buf[:rlen])

				line, err := data.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						log.WithFields(log.Fields{
							"peer": peer.String(),
							"line": string(line),
							"msg":  "detected EOF before newline",
						}).Warn("invalid line")
					} else {
						log.WithFields(log.Fields{
							"peer": fmt.Sprintf("%v:%v", peer.IP.String(), peer.Port),
						}).Error(err)
					}
					return
				}
				if len(line) > 0 {
					line = line[:len(line)-1] // removes trailing '/n'
					select {
					case u.data <- line:
						log.WithFields(log.Fields{
							"peer": peer.String(),
							"line": string(line),
						}).Debug("received line")
					default:
						log.WithFields(log.Fields{
							"peer":          peer.String(),
							"line":          line,
							"channel depth": len(u.data),
						}).Warn("channel full - discarding value")
					}
				}

			}
		}
	}()
	return nil
}
