#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

echo "installing kubectl"
curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
chmod +x ./kubectl
echo "./kubectl version --short --client=true"
./kubectl version --short --client=true
sudo mv -v ./kubectl /usr/local/bin/

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
