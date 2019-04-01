package cmd

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubernetes/pkg/kubectl/util/templates"
	// Initialize all known client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	viewServiceaccountKubeconfigLong = templates.LongDesc(`
Show a kubeconfig setting to access the apiserver with a specified serviceaccount.

The kubeconfig setting will reference the serviceaccount token and use the
current server and namespace
`)

	viewServiceaccountKubeconfigExample = templates.Examples(`
# Show a kubeconfig setting of serviceaccount/default
kubectl view-serviceaccount-kubeconfig default

# Show a kubeconfig setting of serviceaccount/bot in namespace/kube-system
kubectl view-serviceaccount-kubeconfig bot -n kube-system
`)
)

// ViewServiceaccountKubeconfig provides information required to show the
// KUBECONFIG setting of serviceaccount
type ViewServiceaccountKubeconfigOptions struct {
	configFlags *genericclioptions.ConfigFlags

	args []string

	genericclioptions.IOStreams
}

// NewViewServiceaccountKubeconfig provides an instance of
// ViewServiceaccountKubeconfig with default value
func NewViewServiceaccountKubeconfigOptions(streams genericclioptions.IOStreams) *ViewServiceaccountKubeconfigOptions {
	return &ViewServiceaccountKubeconfigOptions{
		configFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}

// NewCmdViewServiceaccountKubeconfig provides a cobra command wrapping
// ViewServiceaccountKubeconfigOptions
func NewCmdViewServiceaccountKubeconfig(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewViewServiceaccountKubeconfigOptions(streams)

	cmd := &cobra.Command{
		Use:     "kubectl view-serviceaccount-kubeconfig SERVICEACCOUNT [options]",
		Short:   "Show a kubeconfig setting to access the apiserver with a specified serviceaccount.",
		Long:    viewServiceaccountKubeconfigLong,
		Example: viewServiceaccountKubeconfigExample,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			c.SilenceUsage = true
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	o.configFlags.AddFlags(cmd.Flags())

	// add the klog flags
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// Workaround for this issue:
	// https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})

	return cmd
}

// Complete sets all information required for showing the KUBECONFIG setting
// of serviceaccount
func (o *ViewServiceaccountKubeconfigOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *ViewServiceaccountKubeconfigOptions) Validate() error {
	if len(o.args) != 1 {
		return fmt.Errorf("exactly one SERVICEACCOUT is required, got %d", len(o.args))
	}

	return nil
}

// Run shows a kubeconfig to access the apiserver with a specified
// serviceaccount
func (o *ViewServiceaccountKubeconfigOptions) Run() error {
	serviceaccountName := o.args[0]

	restConfig, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	kubeConfig := o.configFlags.ToRawKubeConfigLoader()

	client := kubernetes.NewForConfigOrDie(restConfig)
	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		return err
	}

	serviceaccount, err := client.CoreV1().ServiceAccounts(namespace).Get(serviceaccountName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Failed to get a serviceaccount/%s in namespace/%s: %v", serviceaccountName, namespace, err)
	}

	if len(serviceaccount.Secrets) < 1 {
		return fmt.Errorf("serviceaccount %s has no secrets", serviceaccount.GetName())
	}

	secret, err := client.CoreV1().Secrets(namespace).Get(serviceaccount.Secrets[0].Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Failed to get a secret: %v", err)
	}

	caCrt, ok := secret.Data["ca.crt"]
	if !ok {
		return fmt.Errorf("key 'ca.crt' not found in %s", secret.GetName())
	}

	token, ok := secret.Data["token"]
	if !ok {
		return fmt.Errorf("key 'token' not found in %s", secret.GetName())
	}

	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return fmt.Errorf("Failed to get current kubeconfig data")
	}

	context := rawConfig.CurrentContext
	cluster := rawConfig.Contexts[context].Cluster
	server := rawConfig.Clusters[cluster].Server

	config := clientcmdapi.Config{
		CurrentContext: context,
		Clusters: map[string]*clientcmdapi.Cluster{
			cluster: &clientcmdapi.Cluster{
				Server:                   server,
				CertificateAuthorityData: caCrt,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			serviceaccount.GetName(): &clientcmdapi.AuthInfo{
				Token: string(token[:]),
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			context: &clientcmdapi.Context{
				Cluster:   cluster,
				AuthInfo:  serviceaccount.GetName(),
				Namespace: namespace,
			},
		},
	}

	data, err := clientcmd.Write(config)
	if err != nil {
		return fmt.Errorf("Failed to serialize a kubeconfig to yaml: %v", err)
	}

	fmt.Fprint(o.Out, string(data[:]))

	return nil
}
