package analyze

import (
	"context"
	"github.com/scylladb/scylla-operator/pkg/analyze/front"
	"github.com/scylladb/scylla-operator/pkg/analyze/snapshot"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms/rules"
	"k8s.io/klog/v2"
	"os"
)

func Analyze(ctx context.Context, ds snapshot.Snapshot) error {
	klog.Infof("Analyzing the cluster for %d available symptom trees...", len(rules.Symptoms))
	diags := make([]symptoms.Issue, 0)
	for _, tree := range rules.Symptoms {
		diag, _, err := symptoms.MatchTree(tree, ds)
		if err != nil {
			klog.Warningf("Error when matching symptom %s: %v", tree.Symptom().Name(), err)
			return err
		}
		if diag != nil {
			diags = append(diags, diag...)
		}
	}
	front.Print(os.Stdout, diags)

	klog.Infof("Scanned the cluster for %d symptom trees", len(rules.Symptoms))
	return nil
}
