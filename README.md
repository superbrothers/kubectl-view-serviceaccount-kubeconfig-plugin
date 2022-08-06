# kubectl plugin view-serviceaccount-kubeconfig SERVICEACCOUNT

This is a kubectl plugin that show a kubeconfig to access the kube-apiserver with a specified serviceaccount.

:warning: In Kubernetes 1.24+, secret-based tokens are no longer auto-created by default for new service accounts. You can use this plugin with `--type request-token` flag to create a new token instead of using the existing secret-based token. See https://kubernetes.io/docs/concepts/configuration/secret/#service-account-token-secrets for more details about the service account token.

```
Usage:
  kubectl view-serviceaccount-kubeconfig SERVICEACCOUNT [options] [flags]

Examples:
  # Show a kubeconfig setting of the service account "myapp" in the current namespace
  kubectl view-serviceaccount-kubeconfig myapp

  # Show a kubeconfig setting of a service account in a custom namespace
  kubectl view-serviceaccount-kubeconfig myapp --namespace myns

  # Request a bound token to authenticate to the kube-apiserver as the service account
  # "myapp" in the current namespace and show it as a kubeconfig setting
  kubectl view-serviceaccount-kubeconfig myapp --type request-token

  # Request a bound token with a custom expiration and show it as a kubeconfig setting
  kubectl view-serviceaccount-kubeconfig myapp --type request-token --duration 10m
```

## Try the plugin

```
# Create a service account "myapp"  in the current namespace
$ kubectl create serviceaccount myapp

# Create a bount token of the service account "myapp" with a custom expiration
# and save it as a kubeconfig file
$ kubectl view-serviceaccount-kubeconfig myapp --request-token --duration 10d > ./kubeconfig

# List pods by the operation of the service account "myapp"
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
