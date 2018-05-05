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

```
$ curl -sLO https://github.com/superbrothers/kubectl-view-kubeconfig-plugin/releases/download/latest/view-kubeconfig.zip
$ mkdir -p ~/.kube/plugin/view-kubeconfig
$ unzip view-kubeconfig.zip -d ~/.kube/plugins/view-kubeconfig
Archive:  view-kubeconfig.zip
   creating: ~/.kube/plugins/view-kubeconfig/
  inflating: ~/.kube/plugins/view-kubeconfig/plugin.yaml
  inflating: ~/.kube/plugins/view-kubeconfig/kill
```

## License

This software is released under the MIT License.
