# TODO: Add Travis build status
# Check metric types collected in plugin description
# Check description of the plugin
# Is go version section correct?
# If plugin name changes make sure all instances in this doc are fixed
# -- Q: Will snap-relay have a make file? 
# TODO: Add info on loading a stand-alone plugin in snap/README.md#examples
# Test all links
# TODO: Add METRICS.md file

# snap streaming collector plugin - relay

This plugin collects metrics from /relay/statsd and /relay/graphite which gather information about statsd and collectd relay protocols respectufully.  

It's used in the [Snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.7+](https://golang.org/dl/) - needed only for building

### Operating systems
All OSs currently supported by plugin:
* Linux/amd64
* Darwin/amd64

### Installation
You can get the pre-built binaries for your OS and architecture at Snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page. Download the plugins package from the latest release, unzip and store in a path you want `snapteld` to access.

### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-relay
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-relay.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/$GOOS/$GOARCH`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* If /proc resides in a different directory, say for example by mounting host /proc inside a container at /hostproc, a proc_path configuration item can be added to snapteld global config or as part of the task manifest for the metrics to be collected.

As part of snapteld global config

```yaml
---
control:
  plugins:
    collector:
      cpu:
        all:
          proc_path: /hostproc
```

Or as part of the task manifest

```json
{
...
    "workflow": {
        "collect": {
            "metrics": {
	      "/intel/relay/collectd" : {}
	    },
	    "config": {
	      "/intel/relay/collectd": {
                "collectdPort": "6126"
	      }
	    },
	    ...
       },
    },
...
```

* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-relay/blob/master/README.md#examples).

## Documentation
### Collected Metrics
Collected metrics have namespace in following format: `/intel/relay/graphite` and `/intel/relay/statsd`.
List of collected metrics in [METRICS.md](https://github.com/intelsdi-x/snap-relay/blob/master/METRICS.md)





//Editted up to this point 




### Examples
#### Run the example
```bash
./examples/run-cpu-file.sh
```

#### Run the plugin manually
Example running CPU collector plugin and writing data to a file using [file publisher plugin](https://github.com/intelsdi-x/snap-plugin-publisher-file).

Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

In one terminal window, open the Snap daemon (in this case with logging set to 1 and trust disabled):
```
$ snapteld -l 1 -t 0
```

In another terminal window:


Load snap-relay plugin:
```
$ snaptel plugin load snap-relay
```
See available metrics for your system:
```
$ snaptel metric list
```




Get influxdb plugin for publishing, appropriate for Linux or Darwin:
```
$ wget  http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-influxdb/latest/linux/x86_64/snap-plugin-publisher-influxdb
```
or
```
$ wget  http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-influxdb/latest/darwin/x86_64/snap-plugin-publisher-influxdb
```

Load influxdb plugin for publishing:
```
$ snaptel plugin load snap-plugin-publisher-influxdb
```

Create a task manifest file (see [exemplary files] (https://github.com/intelsdi-x/snap-plugin-collector-cpu/blob/master/examples/tasks/)):
    
```json
{
  "version": 1,
  "schedule": {
    "type": "simple",
    "interval": "5s"
  },
  "workflow": {
    "collect": {
      "metrics": {
        "/intel/procfs/cpu/*/active_jiffies": {},
        "/intel/procfs/cpu/*/active_percentage": {},
        "/intel/procfs/cpu/*/guest_jiffies": {},
        "/intel/procfs/cpu/*/guest_nice_jiffies": {},
        "/intel/procfs/cpu/*/guest_nice_percentage": {},
        "/intel/procfs/cpu/*/guest_percentage": {},
        "/intel/procfs/cpu/*/idle_jiffies": {},
        "/intel/procfs/cpu/*/idle_percentage": {},
        "/intel/procfs/cpu/*/iowait_jiffies": {},
        "/intel/procfs/cpu/*/iowait_percentage": {},
        "/intel/procfs/cpu/*/irq_jiffies": {},
        "/intel/procfs/cpu/*/irq_percentage": {},
        "/intel/procfs/cpu/*/nice_jiffies": {},
        "/intel/procfs/cpu/*/nice_percentage": {},
        "/intel/procfs/cpu/*/softirq_jiffies": {},
        "/intel/procfs/cpu/*/softirq_percentage": {},
        "/intel/procfs/cpu/*/steal_jiffies": {},
        "/intel/procfs/cpu/*/steal_percentage": {},
        "/intel/procfs/cpu/*/system_jiffies": {},
        "/intel/procfs/cpu/*/system_percentage": {},
        "/intel/procfs/cpu/*/user_jiffies": {},
        "/intel/procfs/cpu/*/user_percentage": {},
        "/intel/procfs/cpu/*/utilization_jiffies": {},
        "/intel/procfs/cpu/*/utilization_percentage": {}
      },
      "config": {
        "/intel/procfs/cpu": {
          "proc_path": "/proc"
        }
      },
      "publish": [
        {
          "plugin_name": "file",
          "config": {
            "file": "/tmp/published_cpu.log"
          }
        }
      ]
    }
  }
}
```


Create a task:
```
$ snaptel task create -t cpu-file.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

Stop previously created task:
```
$ snaptel task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-cpu/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-cpu/pulls).

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-cpu/issues).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. The full project is at http://github.com/intelsdi-x/snap.
To reach out on other use cases, visit [Slack](http://slack.snap-telemetry.io).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Katarzyna Zabrocka](https://github.com/katarzyna-z)