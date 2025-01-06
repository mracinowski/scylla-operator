package sources

import (
	"context"
	"fmt"
	scyllaversioned "github.com/scylladb/scylla-operator/pkg/client/scylla/clientset/versioned"
	scyllav1listers "github.com/scylladb/scylla-operator/pkg/client/scylla/listers/scylla/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/pager"
)

type DataSource struct {
	PodLister                   corev1listers.PodLister
	ServiceLister               corev1listers.ServiceLister
	SecretLister                corev1listers.SecretLister
	ConfigMapLister             corev1listers.ConfigMapLister
	ServiceAccountLister        corev1listers.ServiceAccountLister
	PersistentVolumeClaimLister corev1listers.PersistentVolumeClaimLister
	ScyllaClusterLister         scyllav1listers.ScyllaClusterLister
}

func BuildListerWithOptions[T any](
	ctx context.Context,
	factory func(cache.Indexer) T,
	listFunc func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error),
	options metav1.ListOptions,
) (T, error) {
	p := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
		return listFunc(ctx, opts)
	}))

	// Prevent users from providing unwanted ones or tempering options that pager controls
	options = metav1.ListOptions{
		LabelSelector: options.LabelSelector,
		FieldSelector: options.FieldSelector,
	}

	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{"NamespaceIndex": cache.MetaNamespaceIndexFunc})
	err := p.EachListItemWithAlloc(ctx, options, func(obj runtime.Object) error {
		err := indexer.Add(obj)
		if err != nil {
			return fmt.Errorf("can't add object to indexer %v: %w", obj, err)
		}
		return nil
	})
	if err != nil {
		return *new(T), fmt.Errorf("can't iterate over list items: %w", err)
	}

	return factory(indexer), nil
}

func BuildLister[T any](ctx context.Context, factory func(cache.Indexer) T, listFunc func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error)) (T, error) {
	return BuildListerWithOptions[T](ctx, factory, listFunc, metav1.ListOptions{})
}

func NewDataSourceFromClients(
	ctx context.Context,
	kubeClient kubernetes.Interface,
	scyllaClient scyllaversioned.Interface,
) (*DataSource, error) {
	podLister, err := BuildLister(ctx, corev1listers.NewPodLister, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Pods(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build pod lister: %w", err)
	}

	serviceLister, err := BuildLister(ctx, corev1listers.NewServiceLister, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Services(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build service lister: %w", err)
	}

	secretLister, err := BuildLister(ctx, corev1listers.NewSecretLister, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Secrets(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build secret lister: %w", err)
	}

	configMapLister, err := BuildLister(ctx, corev1listers.NewConfigMapLister, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().ConfigMaps(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build config map lister: %w", err)
	}

	serviceAccountLister, err := BuildLister(ctx, corev1listers.NewServiceAccountLister, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().ServiceAccounts(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build service account lister: %w", err)
	}

	persistentVolumeClaimLister, err := BuildLister(ctx, corev1listers.NewPersistentVolumeClaimLister, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().PersistentVolumeClaims(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build persistent volume claim lister: %w", err)
	}

	scyllaClusterLister, err := BuildLister(ctx, scyllav1listers.NewScyllaClusterLister, func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
		return scyllaClient.ScyllaV1().ScyllaClusters(corev1.NamespaceAll).List(ctx, options)
	})
	if err != nil {
		return nil, fmt.Errorf("can't build scylla cluster lister: %w", err)
	}

	return &DataSource{
		PodLister:                   podLister,
		ServiceLister:               serviceLister,
		SecretLister:                secretLister,
		ConfigMapLister:             configMapLister,
		ServiceAccountLister:        serviceAccountLister,
		PersistentVolumeClaimLister: persistentVolumeClaimLister,
		ScyllaClusterLister:         scyllaClusterLister,
	}, nil
}
