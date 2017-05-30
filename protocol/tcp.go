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
	"bufio"
	"net"
	"strings"

	"time"

	"io"

	log "github.com/Sirupsen/logrus"
)

type tcpListener struct {
	data     chan []byte
	listener *net.TCPListener
	done     chan struct{}
}

type tcpOption func(u *tcpListener) tcpOption

func NewTCPListener(opts ...tcpOption) *tcpListener {
	listener := &tcpListener{
		data: make(chan []byte, 100),
		done: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(listener)
	}
	return listener
}

func TCPListenerOption(listener *net.TCPListener) tcpOption {
	return func(t *tcpListener) tcpOption {
		prev := t.listener
		t.listener = listener
		return TCPListenerOption(prev)
	}
}

func (t *tcpListener) Data() chan []byte {
	return t.data
}

func (t *tcpListener) Stop() {
	close(t.done)
}

func (t *tcpListener) listen() error {
	if t.listener == nil {
		tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return err
		}
		t.listener, err = net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			return err
		}
		log.WithFields(
			log.Fields{
				"addr": t.listener.Addr().String(),
			},
		).Debug("tcp listening started")
	}
	return nil
}

func (t *tcpListener) handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-t.done:
			break
		default:
			log.WithFields(log.Fields{
				"ReadDeadLine": time.Now().Add(1 * time.Minute),
			}).Debug("reading line")
			conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					log.WithFields(log.Fields{
						"peer": conn.RemoteAddr().String(),
						"line": string(line),
						"msg":  "detected EOF before newline",
					}).Warn("invalid line")
				} else {
					log.WithFields(log.Fields{
						"peer": conn.RemoteAddr().String(),
					}).Error(err)
				}
				return
			}
			if len(line) > 0 {
				line = line[:len(line)-1] // removes trailing '/n'
				select {
				case t.data <- line:
					log.WithFields(log.Fields{
						"peer": conn.RemoteAddr().String(),
						"line": string(line),
					}).Debug("received line")
				default:
					log.WithFields(log.Fields{
						"peer":          conn.RemoteAddr().String(),
						"line":          line,
						"channel depth": len(t.data),
					}).Warn("channel full - discarding value")
				}
			}
		}
	}
}

func (t *tcpListener) Start() error {
	if err := t.listen(); err != nil {
		return err
	}
	log.WithField("addr", t.listener.Addr().String()).Debug("started TCP listener")

	go func() {
	L:
		for {
			select {
			case <-t.done:
				log.WithField("addr", t.listener.Addr().String()).Debug("stopped TCP listener")
				t.listener.Close()
				break L
			default:
				// listen for incoming requests
				conn, err := t.listener.AcceptTCP()
				if err != nil {
					if strings.Contains(err.Error(), "use of closed network connection") {
						break
					}
					log.WithFields(log.Fields{
						"addr": t.listener.Addr().String(),
					}).Error(err)
					break
				}
				// Handle connection
				go t.handleConn(conn)
			}
		}
	}()
	return nil
}
