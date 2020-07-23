#!/bin/bash

docker stop rosetta
docker rm rosetta
docker run -d -p 8080:8080 --name rosetta rosetta:1.0
