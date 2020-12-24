#!/bin/bash

LIB_PATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
. $LIB_PATH/common.sh

parse_args $1

aws dynamodb delete-table \
   $ENDPOINT \
   --table-name $TABLE_NAME