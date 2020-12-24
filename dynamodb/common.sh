#!/bin/bash

LOCAL_DDB="http://dynamodb:8000"
ENDPOINT="--endpoint-url $LOCAL_DDB"

# This value is in accordance with the `defaultTableName` constant
# in pkg/dynamodb/dynamodb.go. If these two values are different
# integration tests may break.
TABLE_NAME="sample-table"

function parse_args() {
    while test $# -gt 0; do
    case "$1" in
        -h|--help)
        echo "$0 [-r|--remote]"
        exit 0
        ;;
        -r|--remote)
        shift
        ENDPOINT=""
        ;;
        *)
        break
        ;;
    esac
    done
}