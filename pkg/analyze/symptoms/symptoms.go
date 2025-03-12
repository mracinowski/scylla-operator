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

//type AndSymptom interface {
//	Symptom
//	SubSymptoms() []*Symptom
//}
//
//type multiSymptom struct {
//	name     string
//	symptoms []*Symptom
//	selector func(*snapshot.DataSource) (bool, error)
//}
//
//func NewMultiSymptom(name string, symptoms []*Symptom) AndSymptom {
//	return &multiSymptom{
//		name:     name,
//		symptoms: symptoms,
//		selector: func(_ *snapshot.DataSource) (bool, error) { panic("not implemented :(") },
//	}
//}
//
//func (m *multiSymptom) Name() string {
//	return m.name
//}
//
//func (m *multiSymptom) Diagnoses() []string {
//	diagnoses := make([]string, 0)
//	for _, sym := range m.symptoms {
//		diagnoses = append(diagnoses, (*sym).Diagnoses()...)
//	}
//	return diagnoses
//}
//
//func (m *multiSymptom) Suggestions() []string {
//	suggestions := make([]string, 0)
//	for _, sym := range m.symptoms {
//		suggestions = append(suggestions, (*sym).Suggestions()...)
//	}
//	return suggestions
//}
//
//func (m *multiSymptom) Match(ds *snapshot.DataSource) ([]front.Diagnosis, error) {
//	match, err := m.selector(ds)
//	if err != nil {
//		return nil, err
//	}
//	if match {
//		// TODO: construct diagnosis
//		return make([]front.Diagnosis, 0), nil
//	}
//	return nil, nil
//}
//
//func (m *multiSymptom) SubSymptoms() []*Symptom {
//	return m.symptoms
//}

type conditionHandler func(*MatchWorkerPool, Symptom, int, chan JobStatus, chan JobStatus)

type SymptomTreeNode interface {
	Name() string
	Symptom() Symptom
	Parent() *SymptomTreeNode
	SetParent(*SymptomTreeNode)
	Handler() conditionHandler
	IsLeaf() bool

	Children() map[string]SymptomTreeNode
	AddChild(SymptomTreeNode) error
}

type symptomTreeNode struct {
	name     string
	parent   *SymptomTreeNode
	symptom  Symptom
	leaf     bool
	children map[string]SymptomTreeNode
	handler  conditionHandler
}

func NewSymptomTreeLeaf(name string, symptom Symptom) SymptomTreeNode {
	return &symptomTreeNode{
		name:     name,
		symptom:  symptom,
		parent:   nil,
		children: nil,
		handler:  nil,
		leaf:     true,
	}
}

func NewSymptomTreeNode(name string, symptom Symptom, handler conditionHandler) SymptomTreeNode {
	return &symptomTreeNode{
		name:     name,
		symptom:  symptom,
		parent:   nil,
		children: make(map[string]SymptomTreeNode),
		handler:  handler,
		leaf:     false,
	}
}

func NewSymptomTreeNodeWithChildren(name string, symptom Symptom, handler conditionHandler, children ...SymptomTreeNode) SymptomTreeNode {
	node := symptomTreeNode{
		name:     name,
		symptom:  symptom,
		parent:   nil,
		children: make(map[string]SymptomTreeNode),
		handler:  handler,
		leaf:     false,
	}

	for _, c := range children {
		err := node.AddChild(c)
		if err != nil {
			klog.Warningf("can't add child symptoms for set %s: %v, name, err")
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

func (s *symptomTreeNode) Children() map[string]SymptomTreeNode {
	return s.children
}

func (s *symptomTreeNode) Parent() *SymptomTreeNode {
	return s.parent
}

func (s *symptomTreeNode) SetParent(parent *SymptomTreeNode) {
	s.parent = parent
}

func (s *symptomTreeNode) Handler() conditionHandler {
	return s.handler
}

func (s *symptomTreeNode) IsLeaf() bool {
	return s.leaf
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

	var thisAsInterface SymptomTreeNode = s
	c.SetParent(&thisAsInterface)
	return nil
}

// Chyba useless
func TrueCondition(w *MatchWorkerPool, symptom Symptom, children int, recv chan JobStatus, send chan JobStatus) {
	w.EnqueueNode(symptom, send, nil)
	for range children {
		_ = <-recv
	}
	close(recv)
}

func OrConditionPropagateFirst(w *MatchWorkerPool, symptom Symptom, children int, recv chan JobStatus, send chan JobStatus) {
	enqueued := false
	for i := 0; i < children; i++ {
		jobStatus := <-recv
		if jobStatus.matched() && !enqueued {
			w.EnqueueNode(symptom, send, jobStatus.Issues)
			enqueued = true
		}
	}
	if !enqueued {
		send <- JobStatus{
			Job:       nil,
			Error:     nil,
			Issues:    make([]Issue, 0),
			SubIssues: make([]Issue, 0),
		}
	}

	close(recv)
}

func OrConditionPropagateAll(w *MatchWorkerPool, symptom Symptom, children int, recv chan JobStatus, send chan JobStatus) {
	matched := false
	subIssues := make([]Issue, 0)
	for i := 0; i < children; i++ {
		jobStatus := <-recv
		subIssues = append(subIssues, jobStatus.Issues...)
		if jobStatus.matched() {
			matched = true
		}
	}
	if matched {
		w.EnqueueNode(symptom, send, subIssues)
	} else {
		send <- JobStatus{
			Job:       nil,
			Error:     nil,
			Issues:    make([]Issue, 0),
			SubIssues: make([]Issue, 0),
		}
	}

	close(recv)
}

func AndCondition(w *MatchWorkerPool, symptom Symptom, children int, recv chan JobStatus, send chan JobStatus) {
	msgSend := false
	subIssues := make([]Issue, 0)
	for i := 0; i < children; i++ {
		jobStatus := <-recv
		subIssues = append(subIssues, jobStatus.Issues...)
		if !jobStatus.matched() && !msgSend {
			jobStatus.SubIssues = append(jobStatus.Issues, jobStatus.SubIssues...)
			jobStatus.Issues = make([]Issue, 0)
			send <- jobStatus
			msgSend = true
		}
	}
	if !msgSend {
		w.EnqueueNode(symptom, send, subIssues)
	}

	close(recv)
}
