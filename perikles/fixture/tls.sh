#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

kubectx do-ams3-odysseia-do-prod
kubens odysseia

cat perikles.crt | base64 | tr -d '\n' > test.txt
kubectl create secret tls -n odysseia perikles-certs \
  --cert=./perikles.crt \
  --key=./perikles.key
