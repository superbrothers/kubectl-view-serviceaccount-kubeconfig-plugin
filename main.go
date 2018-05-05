package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubectl/pkg/pluginutils"
)

func init() {
	// Initialize glog flags
	flag.CommandLine.Set("logtostderr", "true")
	flag.CommandLine.Set("v", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_V"))
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: kubectl plugin view-kubeconfig SERVICEACCOUNT")
	}

	serviceaccountName := os.Args[1]

	restConfig, kubeConfig, err := pluginutils.InitClientAndConfig()
	if err != nil {
		log.Fatalf("Failed to init client and config: %v", err)
	}

	client := kubernetes.NewForConfigOrDie(restConfig)
	namespace, _, _ := kubeConfig.Namespace()

	serviceaccount, err := client.CoreV1().ServiceAccounts(namespace).Get(serviceaccountName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get a serviceaccount: %v", err)
	}

	if len(serviceaccount.Secrets) < 1 {
		log.Fatalf("serviceaccount %s has no secrets", serviceaccount.GetName())
	}

	secret, err := client.CoreV1().Secrets(namespace).Get(serviceaccount.Secrets[0].Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get a secret: %v", err)
	}

	caCrt, ok := secret.Data["ca.crt"]
	if !ok {
		log.Fatalf("key 'ca.crt' not found in %s", secret.GetName())
	}

	token, ok := secret.Data["token"]
	if !ok {
		log.Fatalf("key 'token' not found in %s", secret.GetName())
	}

	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		log.Fatal("Failed to get current kubeconfig data")
	}

	context := rawConfig.CurrentContext
	cluster := rawConfig.Contexts[context].Cluster
	server := rawConfig.Clusters[cluster].Server

	config := clientcmdapi.Config{
		CurrentContext: context,
		Clusters: map[string]*clientcmdapi.Cluster{
			cluster: &clientcmdapi.Cluster{
				Server: server,
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
		log.Fatalf("Failed to serialize a kubeconfig to yaml: %v", err)
	}

	fmt.Print(string(data[:]))
}
