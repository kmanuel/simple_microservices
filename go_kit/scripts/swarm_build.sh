#!/bin/sh

cd service/crop
docker build -t swarm_crop .
cd ../..

cd service/most_significant_image
docker build -t swarm_most_significant_image .
cd ../..

cd service/optimization
docker build -t swarm_optimization .
cd ../..

cd service/portrait
docker build -t swarm_portrait .
cd ../..

cd service/screenshot
docker build -t swarm_screenshot .
cd ../..

cd service/scale_service
docker build -t swarm_scale_service .
cd ../..


cd service/gateway
docker build -t swarm_gateway .
cd ..
