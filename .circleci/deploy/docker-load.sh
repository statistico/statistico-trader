#!/bin/bash

set -e

docker load -i /tmp/workspace/docker-cache/statisticotrader_console.tar
docker load -i /tmp/workspace/docker-cache/statisticotrader_envoyproxy.tar
