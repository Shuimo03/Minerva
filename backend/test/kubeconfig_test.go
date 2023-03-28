package test

import (
	"backend/config/kubeconfig"
	"net/url"
	"testing"
)

func TestKubeClient(t *testing.T) {
	uri, uriError := url.Parse("")
	if uriError != nil {
		t.Fatalf("TODO", uriError)
	}
	_, err := kubeconfig.KubeClient(uri)
	if err != nil {
		t.Fatalf("TODO", err)
	}
}
