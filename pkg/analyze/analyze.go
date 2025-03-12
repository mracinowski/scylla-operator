package analyze

import (
	"context"
	"fmt"
	"github.com/scylladb/scylla-operator/pkg/analyze/front"
	"github.com/scylladb/scylla-operator/pkg/analyze/snapshot"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms/rules"
	"k8s.io/klog/v2"
	"runtime"
)

func Analyze(ctx context.Context, ds snapshot.Snapshot) error {
	statusChan := make(chan symptoms.JobStatus)
	matchWorkerPool := symptoms.NewDefaultMatchWorkerPool(ctx, ds, statusChan, runtime.NumCPU(), symptoms.Worker)
	matchWorkerPool.Start()
	defer close(statusChan)
	defer matchWorkerPool.Finish()

	/*
		for _, tree := range rules.SymptomTests {
			matchWorkerPool.EnqueueTree(tree, statusChan)
		}
		enqueued := len(rules.SymptomTests)
	*/

	for _, tree := range rules.Symptoms {
		matchWorkerPool.EnqueueTree(tree, statusChan)
	}
	enqueued := len(rules.Symptoms)

	klog.Infof("enqueued %d symptom trees", enqueued)

	finished := 0
	for {
		done := false

		select {
		case <-ctx.Done():
			done = true
		case status := <-statusChan:
			finished++

			if status.Error != nil {
				klog.Warningf("symptom %s error: %v", status.Job.Symptom.Name(), status.Error)
			}
			if status.Issues != nil {
				fmt.Println("Main issue:")
				for _, issue := range status.Issues {
					err := front.Print([]front.Diagnosis{front.NewDiagnosis(issue.Symptom, issue.Resources)})
					if err != nil {
						klog.Warningf("can't print diagnosis: %v", err)
					}
				}
			}
			if status.SubIssues != nil {
				fmt.Println("Sub Issues:")
				for _, issue := range status.SubIssues {
					err := front.Print([]front.Diagnosis{front.NewDiagnosis(issue.Symptom, issue.Resources)})
					if err != nil {
						klog.Warningf("can't print diagnosis: %v", err)
					}
				}
			}

			if finished == enqueued {
				done = true
			}
		}

		if done {
			break
		}
	}

	klog.Infof("scanned the cluster for %d symptom trees", enqueued)
	return nil
}
