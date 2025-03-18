package analyze

import (
	"context"
	"github.com/scylladb/scylla-operator/pkg/analyze/front"
	"github.com/scylladb/scylla-operator/pkg/analyze/snapshot"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms/rules"
	"k8s.io/klog/v2"
)

func Analyze(ctx context.Context, ds snapshot.Snapshot) error {

	for _, tree := range rules.Symptoms {
		diags, _, err := symptoms.MatchTree(tree, ds)
		if err != nil {
			klog.Warningf("Error when matching symptom %v", tree.Symptom().Name())
		}
		for _, issue := range diags {
			err := front.Print([]front.Diagnosis{front.NewDiagnosis(issue.Symptom, issue.Resources)})
			if err != nil {
				klog.Warningf("can't print diagnosis: %v", err)
			}
		}
	}

	klog.Infof("scanned the cluster for %d symptom trees", len(rules.Symptoms))
	return nil
}
