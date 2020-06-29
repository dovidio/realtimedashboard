# Real time app download visualization

The aim of this project is to show a world map dashboard containing app downloads information.
Every time a download is performed, the dashboard should update automatically.
![real time dashboard](real-time-dashboard.png)


## Technologies
The technologies used are the following:
- MongoDB for the persistence using the [change streams ](https://www.mongodb.com/blog/post/an-introduction-to-change-streams)
 feature for real time updates
- Go backend server, with pretty much no dependencies beside the mongodb driver (server folder)
- An Angular Frontend Application (frontend folder)
- Mapbox GL JS for the map visualization

The backend server pushes changes to the frontend application using a web socket
![architecture overview](architecture-overview.png)

## Prerequisites
bash
docker-compose 

## Run the app

### With Docker
The easiest way to run the whole app is to run the run script
```bash
./run.sh
```

This will spin up all the necessary containers (mongodb, the server app, the angular app) by using docker compose.
Additionally it will take care of setting up mongodb as a replicaset, which is needed in order to have [change streams ](https://www.mongodb.com/blog/post/an-introduction-to-change-streams)
