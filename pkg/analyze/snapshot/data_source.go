package snapshot

import (
	"context"
	"fmt"
	scyllaversioned "github.com/scylladb/scylla-operator/pkg/client/scylla/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/pager"
	"reflect"
)

type Snapshot interface {
	List(objType reflect.Type) []interface{}
	All() map[reflect.Type][]interface{}
	Add(obj interface{})
}

type snapshot struct {
	objects map[reflect.Type][]interface{}
}

func (ds *snapshot) List(objType reflect.Type) []interface{} {
	list, exists := ds.objects[objType]
	if !exists {
		return make([]interface{}, 0)
	}
	return list
}

func (ds *snapshot) All() map[reflect.Type][]interface{} {
	return ds.objects
}

func (ds *snapshot) Add(obj interface{}) {
	t := reflect.TypeOf(obj)
	ds.objects[t] = append(ds.objects[t], obj)
}

func BuildListWithOptions(
	ctx context.Context,
	ds *Snapshot,
	listFunc func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error),
	options metav1.ListOptions,
) error {
	p := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
		return listFunc(ctx, opts)
	}))

	// Prevent users from providing unwanted ones or tempering options that pager controls
	options = metav1.ListOptions{
		LabelSelector: options.LabelSelector,
		FieldSelector: options.FieldSelector,
	}

	err := p.EachListItemWithAlloc(ctx, options, func(obj runtime.Object) error {
		(*ds).Add(obj)
		return nil
	})
	if err != nil {
		return fmt.Errorf("can't iterate over list items: %w", err)
	}

	return nil
}

func BuildList(ctx context.Context, ds *Snapshot, listFunc func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error)) error {
	return BuildListWithOptions(ctx, ds, listFunc, metav1.ListOptions{})
}

func NewDataSourceFromClients(
	ctx context.Context,
	kubeClient kubernetes.Interface,
	scyllaClient scyllaversioned.Interface,
) (Snapshot, error) {
	var ds Snapshot = &snapshot{
		objects: make(map[reflect.Type][]interface{}),
	}

	err := BuildList(ctx, &ds, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Pods(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build pod lister: %w", err)
	}

	err = BuildList(ctx, &ds, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Services(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build service lister: %w", err)
	}

	err = BuildList(ctx, &ds, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Secrets(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build secret lister: %w", err)
	}

	err = BuildList(ctx, &ds, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().ConfigMaps(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build config map lister: %w", err)
	}

	err = BuildList(ctx, &ds, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().ServiceAccounts(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build service account lister: %w", err)
	}

	err = BuildList(ctx, &ds, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return scyllaClient.ScyllaV1().ScyllaClusters(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build scylla cluster lister: %w", err)
	}

	return ds, nil
}
