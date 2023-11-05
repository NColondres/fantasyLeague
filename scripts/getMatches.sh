#! /bin/bash

PUUIDs=$(curl -s -b cookiesCreator.txt -c cookiesCreator.txt -X GET localhost:8080/lobby | jq -r '.players[].puuid')

for PUUID in ${PUUIDs[@]}; do
    curl -s -b cookiesCreator.txt -c cookiesCreator.txt -X GET localhost:8080/player/matches/${PUUID} | jq
done