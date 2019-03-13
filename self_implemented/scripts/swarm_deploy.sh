#!/bin/sh

docker stack deploy --compose-file scripts/swarm-compose.yml self_impl_swarm
