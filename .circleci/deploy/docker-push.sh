#!/bin/bash

set -e

aws ecr get-login --no-include-email --region $AWS_DEFAULT_REGION | bash

docker tag "statisticotrader_console" "$AWS_ECR_ACCOUNT_URL/statistico-trader:$CIRCLE_SHA1"
docker push "$AWS_ECR_ACCOUNT_URL/statistico-trader:$CIRCLE_SHA1"
