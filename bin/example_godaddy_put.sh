#!/bin/bash

API_KEY=$1
API_SECRET=$2
DOMAIN=$3
NAME=$4
IP=$5

curl -X PUT "https://api.godaddy.com/v1/domains/${DOMAIN}/records/A/${NAME}" -H "accept: application/json" -H "Content-Type: application/json" -H "Authorization: sso-key ${API_KEY}:${API_SECRET}" -d "[ { \"data\": \"${IP}\", \"ttl\": 600 }]"
