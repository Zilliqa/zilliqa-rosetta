#!/bin/bash

docker stop rosetta_standalone
docker rm rosetta_standalone
docker run -d -p 8080:8080 --name rosetta_standalone rosetta_standalone:1.0
