#!/bin/sh

cat /tmp/kubeconfig.yaml | sed 's/0.0.0.0/127.0.0.1/g' > /tmp/kubeconfig.tmp.yaml
cat /tmp/kubeconfig.tmp.yaml > /tmp/kubeconfig.yaml
