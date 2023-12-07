#!/bin/bash
set -e

mkdir -p /var/data

atest run -p test-suite.yaml --report md --level debug
