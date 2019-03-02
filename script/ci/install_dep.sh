#!/bin/bash

curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/darwin/amd64/kubectl
chmod +x ./kubectl
sudo mv -v ./kubectl /usr/bin/

curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

mkdir -p ${GOPATH}/src/github.com/operator-framework
cd ${GOPATH}/src/github.com/operator-framework
git clone https://github.com/operator-framework/operator-sdk
cd operator-sdk
git checkout master
make dep
make install