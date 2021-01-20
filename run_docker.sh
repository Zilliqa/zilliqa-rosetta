#!/bin/bash

docker stop rosetta
docker rm rosetta
docker run -d -p 4201:4201 -p 4301:4301 -p 4501:4501 -p 33133:33133 -p 8080:8080 --name rosetta rosetta:1.0
