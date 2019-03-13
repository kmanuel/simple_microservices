#!/bin/sh

sh scripts/swarm_build.sh
sh scripts/swarm_deploy.sh
python scripts/scale_service.py
