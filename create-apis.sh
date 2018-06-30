#!/bin/bash

for i in {0..100}
do
  wso2am-cli api create \
    --name myapi_$i \
    --context myapi_$i \
    --version 1.0 \
    --definition ./swagger.json \
    --production-url http://localhost/ \
    --sandbox-url http://localhost/ \
    --gateway-env "Production and Sandbox" \
    --publish \
    --update
done
