# kubectl plugin view-kubeconfig SERVICEACCOUNT

kubectl plugin that show a kubeconfig to access the apiserver with a specified serviceaccount.

```
$ kubectl plugin view-kubeconfig -h
Show a kubeconfig setting to access the apiserver with a specified serviceaccount.

Examples:
  # Show a kubeconfig setting of serviceaccount/default
  kubectl plugin view-kubeconfig default

Usage:
  kubectl plugin view-kubeconfig [options]

Use "kubectl options" for a list of global command-line options (applies to all commands).
```

## Install the plugin

If you are on macOS, you can install with homebrew:
```
$ brew tap superbrothers/kubectl-view-kubeconfig-plugin
$ brew install kubectl-view-kubeconfig-plugin
```

If you are on Linux, you can install with the following steps:
```
$ curl -sLO https://github.com/superbrothers/kubectl-view-kubeconfig-plugin/releases/download/$(curl -sL https://raw.githubusercontent.com/superbrothers/kubectl-view-kubeconfig-plugin/master/version.txt)/view-kubeconfig-linux-amd64.zip
$ mkdir -p ~/.kube/plugins/view-kubeconfig
$ unzip view-kubeconfig.zip -d ~/.kube/plugins/view-kubeconfig
Archive:  view-kubeconfig.zip
   creating: ~/.kube/plugins/view-kubeconfig/
  inflating: ~/.kube/plugins/view-kubeconfig/plugin.yaml
  inflating: ~/.kube/plugins/view-kubeconfig/view-kubeconfig
```

## License

This software is released under the MIT License.
