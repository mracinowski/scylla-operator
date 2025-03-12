package rules

import "github.com/scylladb/scylla-operator/pkg/analyze/symptoms"

var symptomSlices = [][]symptoms.SymptomTreeNode{
	StorageSymptoms,
	DummySymptoms,
}

var Symptoms []symptoms.SymptomTreeNode

var SymptomTests = []symptoms.SymptomTreeNode{
	OrTestTree(),
	AndTestTree(),
}

func init() {
	for _, s := range symptomSlices {
		Symptoms = append(Symptoms, s...)
	}
}
