package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd/api/latest"
	"k8s.io/kubectl/pkg/util/templates"

	// Initialize all known client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	viewServiceaccountKubeconfigLong = templates.LongDesc(`
Show a kubeconfig setting for serviceaccount from bound token or secret-based token.

Note that in Kubernetes 1.24+, secret-based tokens are no longer auto-created
by default for new service accounts. Using bound tokens created by "kubectl
create token" command to access the Kubernetes API is recommended instead.
`)

	viewServiceaccountKubeconfigExample = templates.Examples(`
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
`)
)

// ViewServiceaccountKubeconfig provides information required to show the
// KUBECONFIG setting of serviceaccount
type ViewServiceaccountKubeconfigOptions struct {
	configFlags *genericclioptions.ConfigFlags
	printFlags  *genericclioptions.PrintFlags
	printObj    printers.ResourcePrinterFunc

	serviceaccountName string

	genericclioptions.IOStreams
}

// NewViewServiceaccountKubeconfig provides an instance of
// ViewServiceaccountKubeconfig with default value
func NewViewServiceaccountKubeconfigOptions(streams genericclioptions.IOStreams) *ViewServiceaccountKubeconfigOptions {
	return &ViewServiceaccountKubeconfigOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		// disabled all output flags except yaml/json.
		printFlags: (&genericclioptions.PrintFlags{
			JSONYamlPrintFlags: genericclioptions.NewJSONYamlPrintFlags(),
		}).WithDefaultOutput("yaml"),

		IOStreams: streams,
	}
}

// NewCmdViewServiceaccountKubeconfig provides a cobra command wrapping
// ViewServiceaccountKubeconfigOptions
func NewCmdViewServiceaccountKubeconfig(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewViewServiceaccountKubeconfigOptions(streams)

	cmd := &cobra.Command{
		Use:     "kubectl view-serviceaccount-kubeconfig [SERVICEACCOUNT] [options]",
		Short:   "Show a kubeconfig setting for serviceaccount from bound token or secret-based token.",
		Long:    viewServiceaccountKubeconfigLong,
		Example: viewServiceaccountKubeconfigExample,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(args); err != nil {
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

	cmd.Flags().StringVarP(o.printFlags.OutputFormat, "output", "o", *o.printFlags.OutputFormat, fmt.Sprintf("Output format. One of : %s.", strings.Join(o.printFlags.AllowedFormats(), "|")))

	o.configFlags.AddFlags(cmd.Flags())

	// add the klog flags
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// Workaround for this issue:
	// https://github.com/kubernetes/kubernetes/issues/17162
	_ = flag.CommandLine.Parse([]string{})

	return cmd
}

// Complete sets all information required for showing the KUBECONFIG setting
// of serviceaccount
func (o *ViewServiceaccountKubeconfigOptions) Complete(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("exactly one SERVICEACCOUT is required, got %d", len(args))
	}

	if len(args) > 0 {
		o.serviceaccountName = args[0]
	}

	printer, err := o.printFlags.ToPrinter()
	if err != nil {
		return err
	}
	o.printObj = printer.PrintObj

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *ViewServiceaccountKubeconfigOptions) Validate() error {
	return nil
}

// Run shows a kubeconfig to access the apiserver with a specified
// serviceaccount
func (o *ViewServiceaccountKubeconfigOptions) Run() error {
	kubeConfig := o.configFlags.ToRawKubeConfigLoader()

	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return fmt.Errorf("Failed to get current kubeconfig data")
	}

	var currentContext string
	if *o.configFlags.Context != "" {
		currentContext = *o.configFlags.Context
	} else {
		currentContext = rawConfig.CurrentContext
	}

	cluster := rawConfig.Contexts[currentContext].Cluster
	server := rawConfig.Clusters[cluster].Server

	var (
		serviceaccountName string
		namespace          string
		token              string
		caCrt              []byte
	)

	// We expect the serviceaccount bound token can be read from the stdin
	// if no arguments is specified.
	if o.serviceaccountName == "" {
		tokenData, err := io.ReadAll(o.IOStreams.In)
		if err != nil {
			return err
		}
		token = string(tokenData)

		namespacedName, err := getServiceAccountNamespacedNameFromBoundToken(token)
		if err != nil {
			return err
		}

		namespace = namespacedName.Namespace
		serviceaccountName = namespacedName.Name

		// We get CA certificate data from the kubeconfig file
		caCrt = rawConfig.Clusters[cluster].CertificateAuthorityData
	} else {
		restConfig, err := o.configFlags.ToRESTConfig()
		if err != nil {
			return err
		}

		client := kubernetes.NewForConfigOrDie(restConfig)

		serviceaccountName = o.serviceaccountName
		namespace, _, err = kubeConfig.Namespace()
		if err != nil {
			return err
		}

		token, caCrt, err = getTokenForServiceAccount(context.Background(), client, namespace, serviceaccountName)
		if err != nil {
			return err
		}
	}

	config := &clientcmdapi.Config{
		CurrentContext: currentContext,
		Clusters: map[string]*clientcmdapi.Cluster{
			cluster: {
				Server:                   server,
				CertificateAuthorityData: caCrt,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			serviceaccountName: {
				Token: token,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			currentContext: {
				Cluster:   cluster,
				AuthInfo:  serviceaccountName,
				Namespace: namespace,
			},
		},
	}

	convertedObj, err := latest.Scheme.ConvertToVersion(config, latest.ExternalVersion)
	if err != nil {
		return err
	}

	return o.printObj(convertedObj, o.Out)
}
