#!/bin/bash

LIB_PATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
. $LIB_PATH/common.sh

parse_args $1

aws dynamodb put-item \
    $ENDPOINT \
    --table-name $TABLE_NAME \
    --item '{
        "Key": {"S": "a"},
        "Value": {"S": "ahsfkahkfahfsla"}
      }'
aws dynamodb put-item \
    $ENDPOINT \
    --table-name $TABLE_NAME \
    --item '{
        "Key": {"S": "b"},
        "Value": {"S": "klsudhfgakjhaasf"}
      }'
aws dynamodb put-item \
    $ENDPOINT \
    --table-name $TABLE_NAME \
    --item '{
        "Key": {"S": "c"},
        "Value": {"S": "akljsndakjsndkja"}
      }'