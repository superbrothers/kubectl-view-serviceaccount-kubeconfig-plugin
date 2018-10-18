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

1. Install [krew](https://github.com/GoogleContainerTools/krew) that is a plugin manager for kubectl.
2. Run:

        kubectl krew install view-serviceaccount-kubeconfig

3. Try it out:

        kubectl view-serviceaccount-kubeconfig default

## License

This software is released under the MIT License.
