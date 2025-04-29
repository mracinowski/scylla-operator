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

func Print(issue symptoms.Issue, compact bool) error {
	const (
		bold       = "\033[1m"
		red        = "\033[31m"
		yellow     = "\033[93m"
		blue       = "\033[34m"
		resetStyle = "\033[0m"
	)

	if issue.Symptom == nil {
		return errors.New("invalid argument: symptom cannot be nil")
	}
	symptom := *issue.Symptom

	fmt.Println(bold + red + symptom.Name() + resetStyle)

	if compact {
		return nil
	}

	diagnoses := symptom.Diagnoses()
	if len(diagnoses) > 0 {
		fmt.Println()
	}
	for _, diagnosis := range diagnoses {
		fmt.Println(red + diagnosis + resetStyle)
	}

	suggestions := symptom.Suggestions()
	if len(suggestions) > 0 {
		fmt.Println()
	}
	for _, suggestion := range suggestions {
		fmt.Println(blue + suggestion + resetStyle)
	}

	resources := issue.Resources
	if len(resources) > 0 {
		fmt.Println("\n" + yellow + "Related resources:" + resetStyle)
	}
	for key := range issue.Resources {
		fmt.Printf(yellow+"%16s -> %s\n", key, getResourceName(resources[key])+resetStyle)
	}

	fmt.Println()

	return nil
}
