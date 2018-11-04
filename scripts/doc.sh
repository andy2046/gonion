#!/usr/bin/env bash

set -euo pipefail

godoc2md github.com/andy2046/gonion \
    > $GOPATH/src/github.com/andy2046/gonion/docs.md
