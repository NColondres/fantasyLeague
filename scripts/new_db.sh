#!/bin/bash


cd $(dirname $0) && cd ..
CURRENT_DIR=$(pwd)

# Stop and cleanup Old Container; Create a new one
docker container stop my_db
docker container rm my_db
docker build -t sql_db_test ${CURRENT_DIR}/db/
docker run --name my_db -d -p 32769:3306 sql_db_test:latest
docker system prune -af --volumes