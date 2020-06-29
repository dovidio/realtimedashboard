# Server Architecture

The server is mainly responsible to subscribe to the database changes and serve them to the web application.
Some important files are the following:

- **/appdonwnload/repository.go**: provides the possibility to list all downloads and insert a single download (used by the simulator)
- **/appdownload/simulator.go**: simulates random download data. It can be enabled by setting the env variable GENERATE_DATA to true.
GENERATE_DATA_PERIOD specifies the period in milliseconds of data point creation.
- **/appdownload/watcher**: watches for database changes and notifies observers. In this case the observer holds a websocket connection
and forward the data to the client, however the watcher does not know about websocket so if in the future a new protocol will be introduced
(e.g. MQTT) one can implement the Observer interface pretty easily.
- **/appdownload/websocket**: setup a websocket connection endpoint and handle client requests while sending them updates

Special mention:

**/db/db-abstaction**: this is an abstraction layer that I introduced to make easier unit testing classes that use the mongo-db client.
This is of course an overkill for this app but if the app but for bigger app it's a good practice.

# Running the server locally
You need golang 1:14 for this.
Go to server, and execute the following commands
```bash
go mod download
go run main.go
```
The server will stop if it cannot connect with mongodb after 10 seconds, so make sure that mongo is up and running.

# Run with docker
The easiest way to run the server would be to use docker or docker-compose in the root directory. This will also run mongodb and the frontend
Notice that mongo needs to be setup to run in replica set mode. For this reason I provided a bash script that spin up all the containers and takes
care of this as well. You only need to run
```bash
./run.sh
```
