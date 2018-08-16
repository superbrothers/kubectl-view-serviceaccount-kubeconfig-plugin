# kubectl plugin view-serviceaccount-kubeconfig SERVICEACCOUNT

kubectl plugin that show a kubeconfig to access the apiserver with a specified serviceaccount.

```
$ kubectl plugin view-serviceaccount-kubeconfig -h
Show a kubeconfig setting to access the apiserver with a specified serviceaccount.

You can also use "view-sa-kubeconfig" as alias for this plugin.

Examples:
  # Show a kubeconfig setting of serviceaccount/default
  kubectl plugin view-serviceaccount-kubeconfig default

Usage:
  kubectl plugin view-serviceaccount-kubeconfig [flags] [options]

Use "kubectl options" for a list of global command-line options (applies to all commands).
```

## Install the plugin

You can install this plugin with [krew](https://github.com/GoogleContainerTools/krew) that is package manager for kubectl plugins.
```
$ kubectl plugin install view-serviceaccount-kubeconfig
```

If you are on macOS, you can install with homebrew:
```
$ brew tap superbrothers/kubectl-view-serviceaccount-kubeconfig-plugin
$ brew install kubectl-view-serviceaccount-kubeconfig-plugin
```

If you are on Linux, you can install with the following steps:
```
$ curl -sLO https://github.com/superbrothers/kubectl-view-serviceaccount-kubeconfig-plugin/releases/download/$(curl -sL https://raw.githubusercontent.com/superbrothers/kubectl-view-serviceaccount-kubeconfig-plugin/master/version.txt)/view-serviceaccount-kubeconfig-linux-amd64.zip
$ mkdir -p ~/.kube/plugins/view-serviceaccount-kubeconfig
$ unzip view-serviceaccount-kubeconfig-linux-amd64.zip -d ~/.kube/plugins/view-serviceaccount-kubeconfig
```

## License

This software is released under the MIT License.
