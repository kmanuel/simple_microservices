docker service create --name minio \
                        -p 9000:9000 \
                        --network=func_functions \
                        --env MINIO_ACCESS_KEY=123 \
                        --env MINIO_SECRET_KEY=12345678 \
                        minio/minio server /data
