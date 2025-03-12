package rules

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/snapshot"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
)

func trueSelector(ss snapshot.Snapshot) []map[string]any {
	s := make([]map[string]any, 0)
	m := make(map[string]any)
	m["true"] = "true"
	s = append(s, m)
	return s
}

func falseSelector(ss snapshot.Snapshot) []map[string]any {
	return make([]map[string]any, 0)
}

var trueSymptom = symptoms.NewSymptom("true", "", "", trueSelector)
var falseSymptom = symptoms.NewSymptom("false", "", "", falseSelector)

func OrTestTree() symptoms.SymptomTreeNode {
	root := symptoms.NewSymptomTreeNode("or", trueSymptom, symptoms.OrConditionPropagateFirst)
	trueNode := symptoms.NewSymptomTreeLeaf("true", trueSymptom)
	falseNode := symptoms.NewSymptomTreeLeaf("false", falseSymptom)
	root.AddChild(trueNode)
	root.AddChild(falseNode)
	return root
}

func AndTestTree() symptoms.SymptomTreeNode {
	root := symptoms.NewSymptomTreeNode("and", trueSymptom, symptoms.AndCondition)
	trueNode := symptoms.NewSymptomTreeLeaf("true", trueSymptom)
	falseNode := symptoms.NewSymptomTreeLeaf("false", falseSymptom)
	root.AddChild(trueNode)
	root.AddChild(falseNode)
	return root
}
