#!/bin/bash

url='http://127.0.0.1:8080/team'

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Content-Encoding: gzip" \
    --data-binary @test.json.gz \
    "$url"
