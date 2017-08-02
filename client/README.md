The intention of this client is to aid in the development of snap-relay by mocking out the framework. By using this client, we can ensure the  plugin is sending the framework the metrics we were expecting. The client is also used in tests for snap-relay. 

### Running the Built-In Client

**start snap-relay** by running the following command in the root of your snap-relay repo:
```
snap-relay --stand-alone --log-level 5
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