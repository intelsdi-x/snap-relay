### Disclaimer: 
The location in docker-compose.yml for getting the task manifest publishInfluxdb.yml (line 76) will need to be updated after PR is merged. This is because the location of the docker-compose should no longer point to my branch (github.com/kjlyon/snap-relay/tree/update-readme/examples/docker-example/publishInfluxdb.yml) but should point to the permenant location of that file in intelsdi-x/snap-relay/examples/docker-example/publishInfluxdb.yml. 



# Running an example in docker-compose

## Requirements 
 * `docker` and `docker-compose` are **installed** and **configured** 
 * this plugin [downloaded and configured](../../README.md#installation) 
 * build snap-relay for Linux by running `GOOS=linux go build -o snap-relay main.go` from the top level of the snap-relay repo

## Example
This [docker-compose example](docker-compose.yml) will load two plugins: snap-relay and snap-plugin-publisher-influxdb, start a [task](publishInfluxdb.yml), and publish metrics to influxDB from the relay collector plugin.

### Start your containers
In a terminal window navigate to the docker-example directory of this plugin and run,
```
$ docker-compose up -d
```

Check that the two plugins and the task manifest were loaded correctly:
```
$ docker logs init
```

![docker-compose-new-take3](https://user-images.githubusercontent.com/21182867/28733581-b1e76b76-7391-11e7-810e-80bdcd219ec6.gif)


### Explore the relay container
```
$ docker exec -it snap ash
```
The above command will open a bash terminal where you can perform normal snaptel commands such as the following,
` $ snaptel plugin list`, `$ snaptel metric list`, `$ snaptel task list`, etc. You can see the full list of snap commands [here](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPTEL.md). 

To exit the container, `$ exit`.

![relay-container-take2](https://user-images.githubusercontent.com/21182867/28698514-d7ba2d1e-72f8-11e7-921d-62e4d39010ff.gif)

### Explore the influxDB container and influx database
```
$ docker exec -it influxdb bash
```
The above command will drop you into a bash terminal in the influxDB container. 
To access the influx database first type, 
```$ influx````  

Specify which database you want to use:
```
$ use snap
```
You can see the full list of series, 
```
$ show series
```
To see the metrics from a specific series,
```
$ select * from "<SOME_SERIES>"
```
Visit [docs.influxdata.com](https://docs.influxdata.com/influxdb/v1.3/tools/shell/) to see the full list of capabilities of the influxDB interactive shell. 

To close the container, `$ exit` and `$ exit` again. 

![influxdb-container](https://user-images.githubusercontent.com/21182867/28698527-e22d0078-72f8-11e7-8c80-ca5f70c42900.gif)
