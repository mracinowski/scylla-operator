package symptoms

import (
	"fmt"
	"github.com/scylladb/scylla-operator/pkg/analyze/snapshot"
)

func MatchTree(node SymptomTreeNode, ds snapshot.Snapshot) ([]Issue, bool, error) {
	if node.IsLeaf() {
		diag, err := node.Symptom().Match(ds)
		if err != nil {
			return nil, false, err
		}
		return diag, len(diag) > 0, nil
	} else {
		if node.getCallback() == nil {
			return nil, false, fmt.Errorf("Symptom %v: Can't match non-leaf symptom with no callback", node.Name())
		}
		if len(node.Children()) == 0 {
			return nil, false, fmt.Errorf("Symptom %v: Can't match non-leaf symptom with no children", node.Name())
		}
		if node.Symptom() == nil {
			return nil, false, fmt.Errorf("Symptom %v: Can't match symptom node with nil symptom", node.Name())
		}
		diags := make([]Issue, 0)
		matchedCount := 0
		for _, child := range node.Children() {
			diag, matched, err := MatchTree(child, ds)
			if err != nil {
				return nil, false, err
			}
			if matched {
				matchedCount++
			}
			diags = append(diags, diag...)
		}
		if node.ConditionMet(matchedCount) {
			diag, err := node.Symptom().Match(ds)
			if err != nil {
				return nil, false, err
			}
			return append(diag, diags...), len(diag) > 0, nil
		} else {
			return diags, false, nil
		}
	}
}
