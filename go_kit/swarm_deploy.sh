#!/bin/sh

docker stack deploy --compose-file swarm-compose.yml self_impl_swarm
