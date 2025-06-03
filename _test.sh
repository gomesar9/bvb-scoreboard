#!/bin/bash

jsongz_file='test.json.gz'
gzip -c test.json >"$jsongz_file"

url='http://127.0.0.1:8080/team'

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Content-Encoding: gzip" \
    --data-binary "@$jsongz_file" \
    "$url"
