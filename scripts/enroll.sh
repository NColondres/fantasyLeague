#!/bin/bash

# Purpose of this script is to enroll two players (provided through arguments).
# Then start the lobby all while using cookies to represent each user.
# Only the creator of the lobby is able to start the lobby.

if [ $# -eq 0 ]
    then
        echo -e "Not enough arguments supplied\nNeed 2 sets of summoner name and region in that order"
        exit 1
fi

echo -e "Enrolling $1 as a creator\n"
curl -s -b cookiesCreator.txt -c cookiesCreator.txt -X POST localhost:8080/enroll \
    -d "{\"summoner\":\"$1\", \"region\":\"$2\"}" | jq

LOBBY_ID=$(awk '/lobby_id/{ print $7 }' cookiesCreator.txt)

# Copy cookiesCreator.txt and remove the puuid
sed '/puuid/d' cookiesCreator.txt > cookiesJoiner.txt
sed '/puuid/d' cookiesCreator.txt > cookiesJoiner2.txt
sed '/puuid/d' cookiesCreator.txt > cookiesJoiner3.txt
echo
echo -e "Enrolling $3 as a joiner to the lobby\n"
curl -s -b cookiesJoiner.txt -c cookiesJoiner.txt -X POST localhost:8080/enroll \
    -d "{\"summoner\":\"$3\", \"region\":\"$4\"}" | jq

echo -e "Enrolling $5 as a joiner to the lobby\n"
curl -s -b cookiesJoiner2.txt -c cookiesJoiner2.txt -X POST localhost:8080/enroll \
    -d "{\"summoner\":\"$5\", \"region\":\"$6\"}" | jq

echo -e "Enrolling $7 as a joiner to the lobby\n"
curl -s -b cookiesJoiner3.txt -c cookiesJoiner3.txt -X POST localhost:8080/enroll \
    -d "{\"summoner\":\"$7\", \"region\":\"$8\"}" | jq

echo -e "Starting lobby: ${LOBBY_ID}\n"
curl -s -b cookiesCreator.txt -c cookiesCreator.txt -X POST localhost:8080/start \
    -d @testLobby.json  | jq

echo
echo -e "Lobby ${LOBBY_ID} Info\n"
curl -s -b cookiesCreator.txt -c cookiesCreator.txt -X GET localhost:8080/lobby | jq

# PUUID_TO_DELETE=$(grep puuid cookiesJoiner3.txt  | awk '{print $7}')
# echo
# echo -e "Deleting player with PUUID: $PUUID_TO_DELETE"
# curl -s -b cookiesCreator.txt -c cookiesCreator.txt -X DELETE localhost:8080/player/$PUUID_TO_DELETE | jq

 rm cookies*



