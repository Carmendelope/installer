#!/bin/bash
BASEPATH=$HOME/go/src/github.com/nalej/
PROJECTS=("authx" "application-manager" "login-api" "system-model" "user-manager" "web" "infrastructure-manager" "public-api" "conductor" "installer" "signup")

cd "${BASEPATH}"

eval $(minikube docker-env)
for p in "${PROJECTS[@]}"; do
    echo "Loading $p"
    #docker import $p/bin/images/$p/image.tar
    docker load < $p/bin/images/$p/image.tar
done

