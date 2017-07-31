<!-- 
# Check metric types collected in plugin description
# Check description of the plugin
# Check description of metrics collected
# TODO: Add info on loading a stand-alone plugin in snap/README.md#examples
# Run pluginsync -> Will add CONTRIBUTING.md and makefile and travis build status (?)-->

# snap streaming collector plugin - relay

This plugin collects metrics from /relay/statsd and /relay/graphite which gather information about statsd and collectd relay protocols respectufully.  

It's used in the [Snap framework](https://github.com/intelsdi-x/snap).

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
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-relay/blob/master/README.md#examples).

## Documentation
### Collected Metrics
The snap-relay plugin allows access to any metric that is exposed by [Collectd](https://collectd.org/) or [Statsd](https://github.com/etsy/statsd) and is publishable to [Graphite](https://graphiteapp.org/). The collected metrics have namespace in following format: `/intel/relay/graphite` and `/intel/relay/statsd`.


## Examples
### Download and run the docker-compose example

Details can be found in [docker-compose example](/examples/docker-example/) folder.


### Run the plugin manually
This example demonstrates running the relay collector plugin and writing to a file using [file publisher plugin](https://github.com/intelsdi-x/snap-plugin-publisher-file).

#### Start Snap and load plugins
In one terminal window, start the Snap daemon (in this case with logging set to 1 and trust disabled):
```
$ snapteld -l 1 -t 0
```

There are two ways of loading plugins: normally which uses the plugin's binary, and remotely which is available when you run the plugin in stand-alone mode. Below we will demonstrate both ways. 

To load snap-relay plugin in stand-alone mode you must first start the plugin. In another terminal window navigate to your local copy of the snap-relay repository and start the plugin:

```
$ go run main.go stand-alone
```

The plugin will list its stand-alone-port value (default is 8182). Open another terminal window and load the plugin remotely by using the stand-alone port value as shown below:
```
$ snaptel plugin load http://localhost:8182
```

Next, we will load the file plugin by using the binary. We must first get the appropriate version for Linux or Darwin:
```
$ wget  http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
```
or
```
$ wget  http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/darwin/x86_64/snap-plugin-publisher-file
```
Load the file plugin for publishing:
```
$ snaptel plugin load snap-plugin-publisher-file
```

Create a task manifest (see [exemplary files](/examples/tasks/))
```
---
  version: 1
  schedule:
    type: "streaming"
  workflow:
    collect:
      metrics:
       /relay/collectd: {}
      publish:
        -
            plugin_name: "file"
            config:
                file: "/tmp/published_relay.log"
```

Create a task:
```
$ snaptel task create -t /examples/tasks/collectd.yml
```

Send data (do this a few times):
```
$ echo "test.first 13 `date +%s`"|nc -u -c localhost 6124
```

See the results:
```
$ cat /tmp/published_relay.log
```
![screen shot 2017-07-27 at 4 07 03 pm](https://user-images.githubusercontent.com/21182867/28695723-d4b6cc66-72e5-11e7-9057-0c8a2690df80.png)

### Running the Built-In Client

Same as the example above, **start snap-relay** by running the following command in the root of your snap-relay repo:
```
go run main.go --stand-alone --log-level 5
```

Now, open a new terminal and type,
```
curl localhost:8182
```
This will print out the **preamble** for the snap-relay plugin. From this, look for where it says `"ListenAddress"`. Copy the address that is printed there, it will look something like this: `"127.0.0.1:62283"`.

In a third terminal, navigate to your snap-relay repo again and **start the built-in client**,
```
go run client/main.go "<number_from_preamble>"
```

Now we will **send data** and watch it be sent by snap-relay and received in the client. Back in your second terminal type the following command. The default TCP_listen_port is `6124`. Unless you manually set it, that is what it will be,  
```
echo "test.first 10 `date +%s`"|nc -c localhost 6124
```

Repeat that above command a couple times. Each time, you should see a `dispatching metrics` log message in snap-relay and a new metric appear in the client. 

![run-builtin-client-take2](https://user-images.githubusercontent.com/21182867/28794816-86d6a692-75ec-11e7-8cb0-0b5f44c29e62.gif)

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-relay/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-relay/pulls).

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-relay/issues).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. The full project is at http://github.com/intelsdi-x/snap.
To reach out on other use cases, visit [Slack](http://slack.snap-telemetry.io).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](https://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Kelly Lyon](https://github.com/kjlyon)