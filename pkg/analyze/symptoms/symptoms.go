package symptoms

import (
	"errors"
	"fmt"
	"github.com/scylladb/scylla-operator/pkg/analyze/snapshot"
	"k8s.io/klog/v2"
)

const DefaultLimit = 4

type Symptom interface {
	Name() string
	Diagnoses() []string
	Suggestions() []string
	Match(snapshot.Snapshot) ([]Issue, error)
}

type symptom struct {
	name        string
	diagnoses   []string
	suggestions []string
	selector    func(snapshot.Snapshot) []map[string]any
}

func NewSymptom(name string, diag string, suggestions string, selector func(snapshot.Snapshot) []map[string]any) Symptom {
	return &symptom{
		name:        name,
		diagnoses:   []string{diag},
		suggestions: []string{suggestions},
		selector:    selector,
	}
}

func (s *symptom) Name() string {
	return s.name
}

func (s *symptom) Diagnoses() []string {
	return s.diagnoses
}

func (s *symptom) Suggestions() []string {
	return s.suggestions
}

func (s *symptom) Match(ds snapshot.Snapshot) ([]Issue, error) {
	res := s.selector(ds)
	if res != nil && len(res) > 0 {
		issues := make([]Issue, len(res))

		var sym Symptom = s
		for i, r := range res {
			issues[i] = NewIssue(&sym, r)
		}

		return issues, nil
	}
	return nil, nil
}

type conditionCallback func(SymptomTreeNode, int) bool

type SymptomTreeNode interface {
	Name() string
	Symptom() Symptom
	SetSymptom(Symptom) error
	Parent() SymptomTreeNode
	SetParent(SymptomTreeNode)
	IsLeaf() bool
	ConditionMet(int) bool
	getCallback() conditionCallback

	Children() map[string]SymptomTreeNode
	AddChild(SymptomTreeNode) error
}

type symptomTreeNode struct {
	name     string
	parent   SymptomTreeNode
	symptom  Symptom
	leaf     bool
	children map[string]SymptomTreeNode
	callback conditionCallback
}

func NewEmptySymptomNode(name string) SymptomTreeNode {
	return &symptomTreeNode{
		name:     name,
		children: make(map[string]SymptomTreeNode),
	}
}

func NewSymptomTreeLeaf(name string, symptom Symptom) SymptomTreeNode {
	return &symptomTreeNode{
		name:     name,
		symptom:  symptom,
		parent:   nil,
		children: nil,
		callback: nil,
		leaf:     true,
	}
}

func NewSymptomTreeNode(name string, symptom Symptom, callback conditionCallback) SymptomTreeNode {
	return &symptomTreeNode{
		name:     name,
		symptom:  symptom,
		parent:   nil,
		children: make(map[string]SymptomTreeNode),
		callback: callback,
		leaf:     false,
	}
}

func NewSymptomTreeNodeWithChildren(name string, symptom Symptom, callback conditionCallback, children ...SymptomTreeNode) SymptomTreeNode {
	node := symptomTreeNode{
		name:     name,
		symptom:  symptom,
		parent:   nil,
		children: make(map[string]SymptomTreeNode),
		callback: callback,
		leaf:     false,
	}

	for _, c := range children {
		err := node.AddChild(c)
		if err != nil {
			klog.Warningf("can't add child symptoms for set %s: %v", name, err)
			return nil
		}
	}
	return &node
}

func (s *symptomTreeNode) Name() string {
	return s.name
}

func (s *symptomTreeNode) Symptom() Symptom {
	return s.symptom
}

func (s *symptomTreeNode) SetSymptom(symptom Symptom) error {
	if symptom == nil {
		return errors.New("Can't set nil symtptom")
	}
	s.symptom = symptom
	return nil
}

func (s *symptomTreeNode) Children() map[string]SymptomTreeNode {
	return s.children
}

func (s *symptomTreeNode) Parent() SymptomTreeNode {
	return s.parent
}

func (s *symptomTreeNode) SetParent(parent SymptomTreeNode) {
	s.parent = parent
}

func (s *symptomTreeNode) IsLeaf() bool {
	return s.leaf
}

func (s *symptomTreeNode) getCallback() conditionCallback {
	return s.callback
}

func (s *symptomTreeNode) AddChild(c SymptomTreeNode) error {
	if c == nil {
		return errors.New("SymptomTreeNode is nil")
	}
	_, isIn := s.children[c.Name()]
	if isIn {
		return errors.New(fmt.Sprintf("symptom already exists: %v", c))
	}
	s.children[c.Name()] = c

	c.SetParent(s)
	return nil
}

func (s *symptomTreeNode) ConditionMet(matched int) bool {
	return s.callback(s, matched)
}

func OrConditionCallback(_ SymptomTreeNode, matched int) bool {
	return matched > 0
}

func AndConditionCallback(node SymptomTreeNode, matched int) bool {
	return len(node.Children()) == matched
}
