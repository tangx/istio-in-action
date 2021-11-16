#!/bin/bash

cd $(dirname $0)

kubectl delete secret wild-tangx-in -n istio-system

kubectl create secret generic wild-tangx-in \
    --from-file=key=./certificates/_.tangx.in.key  \
    --from-file=cert=./certificates/_.tangx.in.crt  \
    -n istio-system

