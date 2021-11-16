#!/bin/bash
#
# lego-letsencrypt.sh
#

cd $(dirname $0)

source .env

export DOMAIN1="*.tangx.in"

lego  --email="${EMAIL}" \
      --key-type rsa2048 \
      --domains="${DOMAIN1}" \
      --path=$(pwd) --dns $DNS_PROVIDER --accept-tos run
