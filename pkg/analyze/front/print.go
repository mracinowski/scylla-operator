package front

import (
	"errors"
	"fmt"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
)

func getResourceName(resource any) string {
	val := reflect.ValueOf(resource)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	meta, ok := val.FieldByName("ObjectMeta").Interface().(metav1.ObjectMeta)
	if !ok {
		return "<unnamed resource>"
	}
	if len(meta.Namespace) == 0 {
		return meta.Name
	}
	return meta.Namespace + "." + meta.Name
}

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
		fmt.Println("Resources:")
		for key := range issue.Resources {
			fmt.Printf("\t%16s -> %s\n", key, getResourceName(resources[key]))
		}
	}

	fmt.Println()

	return nil
}
