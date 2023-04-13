# kubectl plugin view-serviceaccount-kubeconfig SERVICEACCOUNT

Show a kubeconfig setting for serviceaccount from bound token or secret-based token.

Note that in Kubernetes 1.24+, secret-based tokens are no longer auto-created
by default for new service accounts. Using bound tokens created by "kubectl
create token" command to access the Kubernetes API is recommended instead.

See [Configure Service Accounts for Pods \| Kubernetes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/) for more information.

```
Usage:
  kubectl view-serviceaccount-kubeconfig [SERVICEACCOUNT] [options] [flags]

Examples:
  # Show a kubeconfig setting using bound token for serviceaccount/myapp in the current namespace
  kubectl create token myapp | kubectl view-serviceaccount-kubeconfig

  # Show a kubeconfig setting using bound token for a serviceaccount in a custom namspcae
  kubectl create token myapp --namespace myns | kubectl view-serviceaccount-kubeconfig

  # Show a kubeconfig setting using bound token with a custom expiration
  kubectl create token myapp --duration 10m | kubectl view-serviceaccount-kubeconfig

  # Show a kubeconfig setting using bound token in JSON format
  kubectl create token myapp | kubectl view-serviceaccount-kubeconfig --output json

  # Show a kubeconfig setting using secret-based token for serviceaccount/myapp in the current namespace
  kubectl view-serviceaccount-kubeconfig default

  # Show a kubeconfig setting using secret-based token for serviceaccount/bot in namespace/kube-system
  kubectl view-serviceaccount-kubeconfig bot -n kube-system
```

## Try the plugin

```
# create a serviceaccount "myapp" in the current namespace
$ kubectl create serviceaccount myapp

# save a kubeconfig setting for serviceaccount "myapp"
$ kubectl create token myapp | kubectl view-serviceaccount-kubeconfig >./kubeconfig

# list pods as serviceaccount "myapp" from outside of kubernetes cluster
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
