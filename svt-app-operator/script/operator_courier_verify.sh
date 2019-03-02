#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

readonly BUNDLE_DIR="${GOPATH}/src/github.com/hongkailiu/operators/svt-app-operator/deploy/BUNDLE_DIR"
echo "ls -al ${BUNDLE_DIR}"
ls -al "${BUNDLE_DIR}"
operator-courier verify "${BUNDLE_DIR}"