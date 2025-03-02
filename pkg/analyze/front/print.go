package front

import (
	"errors"
	"fmt"

	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
)

func PrintSymptom(symptom symptoms.Symptom, compact bool) {
	fmt.Println(symptom.Name())

	if compact {
		return
	}

	diagnoses := symptom.Diagnoses()
	if len(diagnoses) == 0 {
		fmt.Println("No diagnoses found for this issue.")
	} else {
		fmt.Println("Diagnoses:")
		for _, diagnosis := range diagnoses {
			fmt.Println("\t", diagnosis)
		}
	}

	suggestions := symptom.Suggestions()
	if len(suggestions) == 0 {
		fmt.Println("No suggestions found for this issue.")
	} else {
		fmt.Println("Suggestions:")
		for _, suggestion := range suggestions {
			fmt.Println("\t", suggestion)
		}
	}
}

func Print(issue symptoms.Issue, compact bool) error {
	if issue.Symptom == nil {
		return errors.New("invalid argument: symptom cannot be nil")
	}
	PrintSymptom(*issue.Symptom, compact)

	if compact {
		return nil
	}

	resources := issue.Resources
	if len(resources) == 0 {
		fmt.Println("No resources related to this issue.")
	} else {
		for key := range issue.Resources {
			fmt.Println(key)
		}
	}

	return nil
}

func FindSymptom(root *symptoms.SymptomSet, name string) *symptoms.Symptom {
	if root == nil {
		return nil
	}

	var result *symptoms.Symptom

	for _, symptom := range (*root).Symptoms() {
		if (*symptom).Name() == name {
			return symptom
		}
	}

	for _, child := range (*root).DerivedSets() {
		child_result := FindSymptom(child, name)
		if child_result != nil {
			return child_result
		}
	}

	return result
}
