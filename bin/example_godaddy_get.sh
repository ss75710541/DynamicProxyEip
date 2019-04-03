
#!/bin/bash


# Show all your records for 'domain'
# GET /v1/domains/{domain}

API_KEY=$1
API_SECRET=$2
DOMAIN=$3

curl -X GET -H "Authorization: sso-key ${API_KEY}:${API_SECRET}" "https://api.godaddy.com/v1/domains/${DOMAIN}/records/A/devopsproxy"
