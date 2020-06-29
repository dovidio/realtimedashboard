docker-compose up -d mongodb

echo "waiting for mongodb to be up and running"

is_healthy() {
    service="$1"
    container_id="$(docker-compose ps -q "$service")"
    health_status="$(docker inspect -f "{{.State.Status}}" "$container_id")"

    if [ "$health_status" = "running" ]; then
        return 0
    else
        return 1
    fi
}

while ! is_healthy  mongodb; do sleep 1; done

echo "setting up mongo replica set"
until docker exec mongodb mongo --eval "rs.initiate({_id : 'rs0', members: [{ _id : 0, host : \"mongodb:27017\" }]});rs.slaveOk(); db.getMongo().setReadPref('nearest');db.getMongo().setSlaveOk();"
do
    echo ...
    sleep 1
done