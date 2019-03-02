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

echo "installing dep"
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

echo "installing operator-framework"
current_dir="$(pwd)"
mkdir -p ${GOPATH}/src/github.com/operator-framework
cd ${GOPATH}/src/github.com/operator-framework
git clone https://github.com/operator-framework/operator-sdk
cd operator-sdk
git checkout master
make dep
make install
cd "${current_dir}" || exit 1

### https://blog.travis-ci.com/2017-10-26-running-kubernetes-on-travis-ci-with-minikube
echo "installing minikube"
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && \
  chmod +x minikube && sudo mv minikube /usr/local/bin/
###https://github.com/kubernetes/minikube/issues/2176
#export MINIKUBE_HOME=$HOME
export CHANGE_MINIKUBE_NONE_USER=true
#export KUBECONFIG=$HOME/.kube/config
echo "starting minikube"
sudo minikube start --vm-driver=none --kubernetes-version=v1.13.3
echo "minikube update-context ..."
ls -al "$HOME/.kube/config"
cat "$HOME/.kube/config"
sudo minikube update-context
ls -al "$HOME/.minikube/"
readonly UUU="$(id -u -n)"
readonly GGG="$(id -u -n)"
echo "UUU=${UUU}; GGG=${GGG}"
sudo chown -R "${UUU}:${GGG}" "$HOME/.minikube/"
ls -al "$HOME/.minikube/"
echo "waiting node to be ready ..."
JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; \
  until kubectl get nodes -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do echo "sleeping 1 sec"; kubectl get node; sleep 1; done

