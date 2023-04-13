package cmd

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/types"
)

func TestGetServiceAccountNamespacedNameFromBoundToken(t *testing.T) {
	tests := []struct {
		token     string
		wantError bool
		expected  *types.NamespacedName
	}{
		{
			token:     "this is invalid token",
			wantError: true,
			expected:  nil,
		},
		{
			token:     "eyJhbGciOiJSUzI1NiIsImtpZCI6IjRMX0dHcnlxX3BxdzdIVmhabjBBWlJnMm1qZ1VNbUJadmM4aTF4SEZTVlUifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNjc4NjM0MTczLCJpYXQiOjE2Nzg2MzA1NzMsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJkZWZhdWx0IiwidWlkIjoiNWY0NDE0MjMtZmRlMS00N2JiLTlkNmQtMDM3NDNjZTU1MWZmIn19LCJuYmYiOjE2Nzg2MzA1NzMsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTpkZWZhdWx0In0.mXBIw5Z_QXRfNR2Ic1-Uwq8tSMaUM7I1xZ4W1LBB8chFTyHNCP0o_G-7JtQH1AOOUPRWUw_sg_sypHtij05OsCx7C61NoBC8NoL4jq6V28duvqj2AhOXTW08Je_wXLe_FD4E1i5n1TU6Wg3tXFzROjB3pt8koM-MBxGo8LiTsJrC41ipZEMZZBF44pmOXG1LKxks3q1JzLcthdVt4yFyy9tW33MPg_i50If1JeD0Cqbb-UQ-PIPNDsGDFaNd0R-KBQy3tnMxP86Hoj1oRymCgQIz5HRkgBsjoYkDul7Cc6uomy4hcM5amV8aonuxLy0mO0xhcMSkEGgEREd3Xv_bFw",
			wantError: false,
			expected: &types.NamespacedName{
				Namespace: "kube-system",
				Name:      "default",
			},
		},
	}

	for i, tt := range tests {
		actual, err := getServiceAccountNamespacedNameFromBoundToken(tt.token)
		gotError := err != nil

		if tt.wantError != gotError {
			t.Errorf("%d: expected %t, got %t, err: %q", i, tt.wantError, gotError, err)
		}

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("%d: expected %#v, got %#v", i, tt.expected, actual)
		}
	}
}
