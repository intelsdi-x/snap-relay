### Disclaimer: The location in docker-compose.yml for getting the task manifest publishInfluxdb.yml (line 76) is not going to work right now. 



# Running an example

## Requirements 
 * `docker` and `docker-compose` are **installed** and **configured** 
 * this plugin [downloaded and configured](../../README.md#installation) 
 * build snap-relay for linux by running `GOOS=linux go build -o snap-relay main.go` from the top level of the snap-relay repo

## Example
[This](/tasks/publishInfluxdb.yml) example task will publish metrics to **influxdb** from a relay collector plugin.

### Start your containers
In the docker-example directory of this plugin run,
```
$ docker-compose up -d
```

Check that the two plugins, and the task manifest were loaded correctly:
```
$ docker logs init
```

![docker-compose-up-d](https://user-images.githubusercontent.com/21182867/28698304-5c68b280-72f7-11e7-943c-87303b0945f0.gif)


### Explore the relay container
```
$ docker exec -it snap ash
```
The above command will open a bash terminal where you can perform normal snaptel commands such as the following,
` $ snaptel plugin list`, `$ snaptel metric list`, `$ snaptel task list`, etc. You can see the full list of snap commands [here](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPTEL.md). 

To exit the container, `$ exit`.

![relay-container-take2](https://user-images.githubusercontent.com/21182867/28698514-d7ba2d1e-72f8-11e7-921d-62e4d39010ff.gif)

### Explore the influxdb container and influx database
```
$ docker exec -it influxdb bash
```
The above command will drop you into a bash terminal in the influxdb container. 
To access the influx database first type, `$ influx`.  

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

To close the container, `$ exit` and `$ exit` again. 

![influxdb-container](https://user-images.githubusercontent.com/21182867/28698527-e22d0078-72f8-11e7-8c80-ca5f70c42900.gif)