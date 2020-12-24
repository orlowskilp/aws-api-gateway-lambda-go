#!/bin/bash

LIB_PATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
. $LIB_PATH/common.sh

parse_args $1

aws dynamodb create-table \
   $ENDPOINT \
   --table-name $TABLE_NAME \
   --attribute-definitions AttributeName=Key,AttributeType=S \
   --key-schema AttributeName=Key,KeyType=HASH \
   --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1