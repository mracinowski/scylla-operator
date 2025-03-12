package symptoms

import (
	"context"
	"github.com/scylladb/scylla-operator/pkg/analyze/snapshot"
)

type Job struct {
	Symptom     Symptom
	SubIssues   []Issue
	ResultsChan chan JobStatus
}

type JobStatus struct {
	Job       *Job
	Error     error
	Issues    []Issue
	SubIssues []Issue
}

func (j JobStatus) matched() bool {
	return len(j.Issues) > 0
}

type MatchWorkerPool struct {
	ds            snapshot.Snapshot
	jobs          chan *Job
	statusChan    chan JobStatus
	numWorkers    int
	started       bool
	worker        func(ctx context.Context, pool *MatchWorkerPool)
	workerContext context.Context
	workerCancel  context.CancelFunc
}

func NewDefaultMatchWorkerPool(
	ctx context.Context,
	ds snapshot.Snapshot,
	statusChan chan JobStatus,
	numWorkers int,
	worker func(ctx context.Context, pool *MatchWorkerPool),
) *MatchWorkerPool {
	workerContext, workerCancel := context.WithCancel(ctx)
	return &MatchWorkerPool{
		ds:            ds,
		jobs:          make(chan *Job, numWorkers),
		statusChan:    statusChan,
		numWorkers:    numWorkers,
		started:       false,
		worker:        worker,
		workerContext: workerContext,
		workerCancel:  workerCancel,
	}
}

func (w *MatchWorkerPool) EnqueueTree(root SymptomTreeNode, results chan JobStatus) {
	if root.IsLeaf() {
		w.EnqueueNode(root.Symptom(), results, nil)
	} else {
		c := make(chan JobStatus)
		go root.Handler()(w, root.Symptom(), len(root.Children()), c, results)
		for _, child := range root.Children() {
			w.EnqueueTree(child, c)
		}
	}
}

func (w *MatchWorkerPool) EnqueueNode(symptom Symptom, results chan JobStatus, subIssues []Issue) {
	w.Enqueue(Job{
		Symptom:     symptom,
		ResultsChan: results,
		SubIssues:   subIssues,
	})
}

func (w *MatchWorkerPool) Enqueue(job Job) {
	w.jobs <- &job
}

// Start initializes and starts the worker pool if it has not been started yet; panics if already started.
// This method is not thread safe.
func (w *MatchWorkerPool) Start() {
	if w.started {
		panic("MatchWorkerPool already started")
	}
	w.started = true
	for i := 0; i < w.numWorkers; i++ {
		go w.worker(w.workerContext, w)
	}
}

func (w *MatchWorkerPool) Finish() {
	w.workerCancel()
	close(w.jobs)
}

func Worker(ctx context.Context, pool *MatchWorkerPool) {
	for {
		select {
		case <-ctx.Done():
			break
		case job := <-pool.jobs:
			diag, err := job.Symptom.Match(pool.ds)
			job.ResultsChan <- JobStatus{
				Job:       job,
				Error:     err,
				Issues:    diag,
				SubIssues: job.SubIssues,
			}
		}
	}
}
