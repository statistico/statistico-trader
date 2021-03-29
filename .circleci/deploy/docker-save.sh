#!/bin/bash

set -e

mkdir -p /tmp/workspace/docker-cache

docker save -o /tmp/workspace/docker-cache/statisticotrader_console.tar statisticotrader_console:latest
