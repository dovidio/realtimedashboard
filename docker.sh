docker-compose up -d mongodb

sleep 5

# setup single node replica set
docker exec -it mongodb mongo --eval "rs.initiate({_id : 'rs0', members: [{ _id : 0, host : \"mongodb:27017\" }]});rs.slaveOk(); db.getMongo().setReadPref('nearest');db.getMongo().setSlaveOk();"

exit