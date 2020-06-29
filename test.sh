./setup-mongo-rs.sh > /dev/null

docker exec mongodb mongo --eval "db.getSiblingDB(\"appdownloads\").dropDatabase();" > /dev/null

# By default the server doesn't generate data
docker-compose up -d server > /dev/null

sleep 3

echo "Testing that server pushes data on new record in DB"
appdownloads=$(curl -s http://localhost:8080/appdownloads | jq length | cat)
code=0
if [[ $appdownlads -eq 0 ]]; then
    # Connect to websocket
    rm -f test
    ls > test 2>&1 &
    websocat --origin http://someorigin ws://localhost:8080/appdownloadssocket > test 2>&1 &

    # when inserting a new record
    docker exec mongodb mongo --eval "db.getSiblingDB(\"appdownloads\").appdownloads.insert({latitude: 0,longitude: 0,country: \"Austria\",app_id: \"IOS_ALERT\",downloaded_at: 1593368544952});" > /dev/null
    docker exec mongodb mongo --eval "db.getSiblingDB(\"appdownloads\").appdownloads.insert({latitude: 0,longitude: 0,country: \"Austria\",app_id: \"IOS_ALERT\",downloaded_at: 1593368544952});" > /dev/null
    docker exec mongodb mongo --eval "db.getSiblingDB(\"appdownloads\").appdownloads.insert({latitude: 0,longitude: 0,country: \"Austria\",app_id: \"IOS_ALERT\",downloaded_at: 1593368544952});" > /dev/null

    sleep 3
    # then websocket pushes new values
    valuesPushed=$(wc -l < test)
    if [[ $valuesPushed -ne 3 ]]; then
        cat test
        echo "\e[91mLightASSERTION ERROR: expected 3 values to be pushed to web socket. found" $valuesPushed
    else 
        echo -e "\e[92mTEST PASSED"
    fi

    # kill the websocket
else
    echo "Test setup failure: previous data was not deleted from db. Number of app downloads: " $appdownloads
    code=-1
fi

exit $code