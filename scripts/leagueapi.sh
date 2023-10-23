#!/bin/bash

cd $(dirname $0) && cd ..
CURRENT_DIR=$(pwd)

docker build ${CURRENT_DIR}/api/ -t leagueofgoonsapi:latest --no-cache
docker run --name leagueofgoonsapi -dp 8080:8080 leagueofgoonsapi:latest
docker logs leagueofgoonsapi -f