#!/bin/sh

zarf tools kubectl get secret zarf-state -n zarf -o jsonpath='{.data.state}' | base64 -d | jq -r .registryInfo.pushPassword
