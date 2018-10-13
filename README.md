# kubectl plugin view-serviceaccount-kubeconfig SERVICEACCOUNT

This is a kubectl plugin that show a kubeconfig to access the apiserver with a specified serviceaccount.

```
Usage:
  kubectl view-serviceaccount-kubeconfig SERVICEACCOUNT [options] [flags]

Examples:
  # Show a kubeconfig setting of serviceaccount/default
  kubectl view-serviceaccount-kubeconfig default
  
  # Show a kubeconfig setting of serviceaccount/bot in namespace/kube-system
  kubectl view-serviceaccount-kubeconfig bot -n kube-system
```

## Try the plugin

```
# create a serviceaccount/bot in namepace/default
$ kubectl create serviceaccount bot

# save a kubeconfig setting of serviceaccount/bot
$ kubectl view-serviceaccount-kubeconfig bot > ./kubeconfig

# list pods as serviceaccount/bot from outside of kubernetes cluster
$ kubectl get pods --kubeconfig=./kubeconfig
```

## Install the plugin

<!-- Krew does not support kubectl 1.12 plugin model yet. See https://github.com/GoogleContainerTools/krew/issues/33
You can install this plugin with [krew](https://github.com/GoogleContainerTools/krew) that is package manager for kubectl plugins.
```
$ kubectl plugin install view-serviceaccount-kubeconfig
```
-->

If you are on macOS, you can install with homebrew:
```
$ brew tap superbrothers/kubectl-view-serviceaccount-kubeconfig-plugin
$ brew install kubectl-view-serviceaccount-kubeconfig-plugin
```

If you are on Linux, you can install with the following steps:
```
$ curl -sLO https://github.com/superbrothers/kubectl-view-serviceaccount-kubeconfig-plugin/releases/download/$(curl -sL https://raw.githubusercontent.com/superbrothers/kubectl-view-serviceaccount-kubeconfig-plugin/master/version.txt)/kubectl-view_serviceaccount_kubeconfig-linux-amd64.zip
$ unzip kubectl-view_serviceaccount_kubeconfig-linux-amd64.zip -d /tmp
$ cp /tmp/kubectl-view_serviceaccount_kubeconfig /tmp/kubectl-view_sa_kubeconfig /usr/local/bin/
```

## License

This software is released under the MIT License.
