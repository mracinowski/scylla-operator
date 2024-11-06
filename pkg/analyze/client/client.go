package client

import (
	"github.com/scylladb/scylla-operator/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

func Start(x genericclioptions.ClientConfig) error {
	_, err := kubernetes.NewForConfig(x.ProtoConfig)
	if err != nil {
		return err
	}

	klog.Infof("Kubeconfig")

	return nil
}
