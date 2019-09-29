#!/bin/sh

cd service/crop
docker build -t kit_crop .
cd ../..

cd service/most_significant_image
docker build -t kit_most_significant_image .
cd ../..

cd service/optimization
docker build -t kit_optimization .
cd ../..

cd service/portrait
docker build -t kit_portrait .
cd ../..

cd service/screenshot
docker build -t kit_screenshot .
cd ../..

cd service/scale_service
docker build -t kit_scale_service .
cd ../..


cd service/gateway
docker build -t kit_gateway .
cd ..
