package rules

import "github.com/scylladb/scylla-operator/pkg/analyze/symptoms"

var Symptoms = symptoms.NewSymptomSet("All", []*symptoms.SymptomSet{
	&StorageSymptoms,
})
