package analyze

import (
	"context"
	"runtime"

	"github.com/scylladb/scylla-operator/pkg/analyze/front"
	"github.com/scylladb/scylla-operator/pkg/analyze/sources"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms/rules"
	"k8s.io/klog/v2"
)

func Analyze(ctx context.Context, ds *sources.DataSource2) error {
	statusChan := make(chan symptoms.JobStatus)
	matchWorkerPool := symptoms.NewMatchWorkerPool(ctx, ds, statusChan, runtime.NumCPU())
	matchWorkerPool.Start()
	defer close(statusChan)
	defer matchWorkerPool.Finish()

	enqueued := matchWorkerPool.EnqueueAll(&rules.Symptoms)
	klog.Infof("enqueued %d symptoms", enqueued)

	finished := 0
	for {
		done := false

		select {
		case <-ctx.Done():
			done = true
		case status := <-statusChan:
			finished++

			if status.Error != nil {
				klog.Warningf("symptom %s error: %v", (*status.Job.Symptom).Name(), status.Error)
			}
			if status.Issues != nil {
				for _, issue := range status.Issues {
					err := front.Print(issue, false)
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

	klog.Infof("scanned the cluster for %d symptoms", enqueued)
	return nil
}
