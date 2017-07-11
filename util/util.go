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

package util

import (
	"sync"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

type ChannelManager interface {
	Add(chan *plugin.Metric)
	Remove(chan *plugin.Metric)
	DispatchMetric(*plugin.Metric)
	Count() int
}

type channelMgr struct {
	*sync.Mutex
	channels map[chan *plugin.Metric]interface{}
}

func NewChannelMgr() *channelMgr {
	return &channelMgr{
		&sync.Mutex{},
		make(map[chan *plugin.Metric]interface{}),
	}
}

func (c *channelMgr) Add(ch chan *plugin.Metric) {
	c.Lock()
	defer c.Unlock()
	c.channels[ch] = nil
}

func (c *channelMgr) Remove(ch chan *plugin.Metric) {
	c.Lock()
	defer c.Unlock()
	delete(c.channels, ch)
	close(ch)
}

func (c *channelMgr) Count() int {
	c.Lock()
	defer c.Unlock()
	return len(c.channels)
}

func (c *channelMgr) DispatchMetric(m *plugin.Metric) {
	c.Lock()
	defer c.Unlock()
	for ch := range c.channels {
		ch <- m
	}
}
