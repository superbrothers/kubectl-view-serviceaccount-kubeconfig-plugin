package cmd

import (
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestViewServiceaccountKubeconfigOptionComplete(t *testing.T) {
	tests := []struct {
		args      []string
		wantError bool
	}{
		{
			args:      []string{},
			wantError: false,
		},
		{
			args:      []string{"default"},
			wantError: false,
		},
		{
			args:      []string{"default", "mysa"},
			wantError: true,
		},
	}

	for i, tt := range tests {
		o := NewViewServiceaccountKubeconfigOptions(genericclioptions.NewTestIOStreamsDiscard())
		err := o.Complete(tt.args)
		if tt.wantError {
			if err == nil {
				t.Errorf("%d: expected err, but got nil", i)
			}
			continue
		}

		if err != nil {
			t.Errorf("%d: unexpected err: %v", i, err)
		}
	}
}
