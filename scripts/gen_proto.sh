#!/bin/bash
CURRENT_DIR=$1
for file in $(find ${CURRENT_DIR}/protos/* -type d)
do
  protoc -I=${file} -I=${CURRENT_DIR}/protos --go_out=${CURRENT_DIR} \
   --go-grpc_out=${CURRENT_DIR} ${file}/*.proto
done