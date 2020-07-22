#!/bin/bash

docker rm rosetta
docker run -p 8080:8080 --name rosetta rosetta:1.0
