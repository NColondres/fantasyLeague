#!/bin/bash

for ((i=1; i<20; i=i+2))
do
    thing1=thing$i
    thing2=thing$(($i + 1))

    echo -e "Enrolling $thing1 as a creator\n"
    curl -s -b cookiesCreator.txt -c cookiesCreator.txt -X POST localhost:8080/enroll \
        -d "{\"summoner\":\"$thing1\", \"region\":\"na1\"}" | jq

    LOBBY_ID=$(awk '/lobby_id/{ print $7 }' cookiesCreator.txt)

    # Copy cookiesCreator.txt and remove the puuid
    sed '/puuid/d' cookiesCreator.txt > cookiesJoiner.txt
    echo
    echo -e "Enrolling $thing2 as a joiner to the lobby\n"
    curl -s -b cookiesJoiner.txt -c cookiesJoiner.txt -X POST localhost:8080/enroll \
        -d "{\"summoner\":\"$thing2\", \"region\":\"na1\"}" | jq

    echo -e "Starting lobby: ${LOBBY_ID}\n"
    curl -s -b cookiesCreator.txt -c cookiesCreator.txt -X POST localhost:8080/start \
        -d "{\"k_d_a\": $i, \"win\": 2, \"creep\": 0.8, \"double\": 1500}" | jq

    echo
    echo -e "Lobby ${LOBBY_ID} Info\n"
    curl -s -b cookiesCreator.txt -c cookiesCreator.txt -X GET localhost:8080/lobby | jq

    rm cookiesCreator.txt cookiesJoiner.txt
    echo 
    echo 
done