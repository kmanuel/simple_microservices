docker service create --name faktory \
                -p 7419:7419 \
                -p 7420:7420 \
                --network=func_functions \
                contribsys/faktory
