#!/bin/bash



curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

current_dir="$(pwd)"
mkdir -p ${GOPATH}/src/github.com/operator-framework
cd ${GOPATH}/src/github.com/operator-framework
git clone https://github.com/operator-framework/operator-sdk
cd operator-sdk
git checkout master
make dep
make install
cd "${current_dir}" || exit 1
