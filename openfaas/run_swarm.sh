cd faas
sh deploy_stack.sh

docker service create \
    -d \
    --mount type=bind,source="$(pwd)"/imgs,target=/data \
    -p "9000:9000" \
    --network=func_functions \
    --name minio \
    --env-file ".env" \
    minio/minio:RELEASE.2018-12-06T01-27-43Z \
    server /data

cd ../crop
faas-cli build -f crop.yml
faas-cli deploy -f crop.yml

cd ../most_significant_image
faas-cli build -f most-significant-image.yml
faas-cli deploy -f most-significant-image.yml

cd ../optimization
faas-cli build -f optimization.yml
faas-cli deploy -f optimization.yml

cd ../portrait
faas-cli build -f portrait.yml
faas-cli deploy -f portrait.yml

cd ../screenshot
faas-cli build -f screenshot.yml
faas-cli deploy -f screenshot.yml

cd ..
