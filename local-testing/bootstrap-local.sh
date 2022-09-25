#!/bin/bash
echo "stopping miriconf containers..."
docker-compose -f ./local-testing/docker-compose.yaml down 

echo "forcing docker build..."
docker-compose -f ./local-testing/docker-compose.yaml build

echo "bringing up docker-compose..."
docker-compose -f ./local-testing/docker-compose.yaml up -d

echo "waiting 15 seconds for mongodb to come up"
sleep 15

echo "loading test data..."
mongorestore --host=localhost --port=27017 --username=root --password=testing ./local-testing/test-mongo-data

echo "complete... access http://localhost:8080/swagger/index.html for API docs..."