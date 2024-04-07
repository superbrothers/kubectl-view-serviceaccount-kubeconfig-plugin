package cmd

import (
	"context"
	"fmt"

	"github.com/go-jose/go-jose/v3/jwt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

func getTokenForServiceAccount(ctx context.Context, client *kubernetes.Clientset, namespace, serviceaccountName string) (string, []byte, error) {
	serviceaccount, err := client.CoreV1().ServiceAccounts(namespace).Get(ctx, serviceaccountName, metav1.GetOptions{})
	if err != nil {
		return "", nil, fmt.Errorf("Failed to get serviceaccount %s/%s: %v", namespace, serviceaccountName, err)
	}

	var serviceAccountSecrets []v1.ObjectReference

	if len(serviceaccount.Secrets) > 0 {
		serviceAccountSecrets = serviceaccount.Secrets
	} else {
		secrets, err := client.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return "", nil, fmt.Errorf("failed to list secrets: %v", err)
		}

		for _, secret := range secrets.Items {
			if secret.Annotations["kubernetes.io/service-account.uid"] == string(serviceaccount.UID) {
				serviceAccountSecrets = append(serviceAccountSecrets, v1.ObjectReference{Name: secret.Name})
				break
			}
		}
	}

	if len(serviceAccountSecrets) < 1 {
		return "", nil, fmt.Errorf(`"serviceaccount %s/%s has no secrets.

In Kubernetes 1.24+, secret-based tokens are no longer auto-created
by default for new service accounts. Using bound tokens created by "kubectl
create token" command to access the Kubernetes API is recommended instead.

Alternatively, you can attach a long-lived token to the service account;
see https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/#create-token
for more information.

Check the help message of this command to see how to show the kubeconfig
setting with a bound token.`, namespace, serviceaccountName)
	}

	var secret *v1.Secret
	for _, secretRef := range serviceAccountSecrets {
		secret, err = client.CoreV1().Secrets(namespace).Get(ctx, secretRef.Name, metav1.GetOptions{})
		if err != nil {
			return "", nil, fmt.Errorf("Failed to get a secret: %v", err)
		}

		if secret.Type == v1.SecretTypeServiceAccountToken {
			break
		}
	}

	if secret == nil {
		return "", nil, fmt.Errorf("serviceAccount %s/%s has no secret type %q", namespace, serviceaccountName, v1.SecretTypeServiceAccountToken)
	}

	token, ok := secret.Data["token"]
	if !ok {
		return "", nil, fmt.Errorf("key 'token' not found in %s", secret.GetName())
	}

	caCrt, ok := secret.Data["ca.crt"]
	if !ok {
		return "", nil, fmt.Errorf("key 'ca.crt' not found in %s", secret.GetName())
	}

	return string(token), caCrt, nil
}

type kubeName struct {
	Name string `json:"name"`
}

type kubeClaims struct {
	Namespace      string   `json:"namespace"`
	ServiceAccount kubeName `json:"serviceaccount"`
}

type claims struct {
	Kubernetes kubeClaims `json:"kubernetes.io"`
}

func getServiceAccountNamespacedNameFromBoundToken(tokenData string) (*types.NamespacedName, error) {
	tok, err := jwt.ParseSigned(tokenData)
	if err != nil {
		return nil, err
	}

	var unsafeClaims claims
	if err := tok.UnsafeClaimsWithoutVerification(&unsafeClaims); err != nil {
		return nil, err
	}

	namespacedName := &types.NamespacedName{
		Namespace: unsafeClaims.Kubernetes.Namespace,
		Name:      unsafeClaims.Kubernetes.ServiceAccount.Name,
	}

	return namespacedName, nil
}
