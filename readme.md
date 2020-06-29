# Real time app download visualization

The aim of this project is to show a world map dashboard containing app downloads information.
Every time a download is performed, the dashboard should update automatically.
Additionally the option to filter by timeframe is provided.

# Technologies
The technologies used are the following:
- MongoDB for the persistence, since it provides change streams
- Go backend server, without dependencies (except for the mongodb driver)
- An Angular Frontend Application
The backend server pushes changes to the frontend application using a web socket
![architecture overview](architecture-overview.png)

# Prerequisites
bash
docker-compose 

# Setup
Run docker.sh