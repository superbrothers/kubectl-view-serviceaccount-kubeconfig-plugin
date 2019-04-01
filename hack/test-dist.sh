#!/usr/bin/env bash

set -e -x -o pipefail

export KREW_ROOT="$(mktemp -d)"
trap "rm -rf $KREW_ROOT" EXIT

"$HOME/.krew/bin/kubectl-krew" install \
    --manifest _dist/view-serviceaccount-kubeconfig.yaml \
    --archive _dist/kubectl-view_serviceaccount_kubeconfig-linux-amd64.zip
