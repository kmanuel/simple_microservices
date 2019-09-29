cd faas
sh deploy_stack.sh
cd ..

docker service create \
    -d \
    --mount type=bind,source="$(pwd)"/imgs,target=/data \
    -p "9000:9000" \
    --network=func_functions \
    --name minio \
    --env-file ".env" \
    minio/minio:RELEASE.2018-12-06T01-27-43Z \
    server /data

