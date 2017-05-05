# snap-relay

## How to start streaming plugin with Snap:
* start Snap `snapteld -l 1 -t 0`
* start relay plugin `go run main.go -stand-alone -log-level 5 -stand-alone-port 8182`
* load plugin in Snap `snaptel plugin load http://localhost:8182 `

From here you can unload the plugin, see metric list, and use Snap as normal.



## How to test streaming plugin without Snap:
Terminal 1: 
1. `cd relay-plugin`
2. `go run main.go -stand-alone --log-level 5`

Terminal 2:
1. `curl localhost:8181`

Terminal 3:
1. `cd relay-plugin`
2. `go run client/main.go <number from preamble print out>`

Terminal 2: 
1. `echo "test.first 10 `date +%s`"|nc -c localhost <number from TCP listener>`

Terminal 1:
output debug messages!  <- it worked! 


![snap-relay-take1](https://cloud.githubusercontent.com/assets/21182867/25767820/9bf9e176-31b1-11e7-82a4-88b0fd5368f2.gif)