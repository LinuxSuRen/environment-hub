#!/bin/bash
set -e

mkdir -p /var/data

atest run -p test-suite.yaml --report md --level debug

if [[ "$K3D" == "true" ]]
then
    atest run -p test-suite-k3d.yaml --report md --level debug
fi
