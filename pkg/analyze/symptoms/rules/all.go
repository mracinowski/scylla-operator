package rules

import "github.com/scylladb/scylla-operator/pkg/analyze/symptoms"

var symptomSlices = [][]symptoms.SymptomTreeNode{
	StorageSymptoms,
}

var Symptoms []symptoms.SymptomTreeNode

func init() {
	for _, s := range symptomSlices {
		Symptoms = append(Symptoms, s...)
	}
}
