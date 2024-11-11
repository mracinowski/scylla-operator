package analyze

import (
	"context"
	"fmt"
	scyllaversioned "github.com/scylladb/scylla-operator/pkg/client/scylla/clientset/versioned"
	scyllav1listers "github.com/scylladb/scylla-operator/pkg/client/scylla/listers/scylla/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/pager"
)

type listingFn = func(metav1.ListOptions) (runtime.Object, error)

type DataSource struct {
	ServiceLister        corev1listers.ServiceLister
	SecretLister         corev1listers.SecretLister
	ConfigMapLister      corev1listers.ConfigMapLister
	ServiceAccountLister corev1listers.ServiceAccountLister
	ScyllaClusterLister  scyllav1listers.ScyllaClusterLister
}

func NewDataSource(
	serviceLister corev1listers.ServiceLister,
	secretLister corev1listers.SecretLister,
	configMapLister corev1listers.ConfigMapLister,
	accountLister corev1listers.ServiceAccountLister,
	clusterLister scyllav1listers.ScyllaClusterLister,
) DataSource {
	return DataSource{
		ServiceLister:        serviceLister,
		SecretLister:         secretLister,
		ConfigMapLister:      configMapLister,
		ServiceAccountLister: accountLister,
		ScyllaClusterLister:  clusterLister,
	}
}

func makeIndexer(ctx context.Context, lister listingFn) (cache.Indexer, error) {
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{"NamespaceIndex": cache.MetaNamespaceIndexFunc})
	objPager := pager.New(pager.SimplePageFunc(lister))
	err := objPager.EachListItem(ctx, metav1.ListOptions{}, func(obj runtime.Object) error {
		err := indexer.Add(obj)
		if err != nil {
			return fmt.Errorf("unable to add item to indexer: %v", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("makeIndexer error: %v", err)
	}

	return indexer, nil
}

func listWrapper[T runtime.Object](ctx context.Context, clientFunc func(context.Context, metav1.ListOptions) (T, error)) listingFn {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return clientFunc(ctx, opts)
	}
}

func NewDataSourceFromClients(
	ctx context.Context,
	client *kubernetes.Clientset,
	scyllaClient *scyllaversioned.Clientset,
) (DataSource, error) {
	refCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	servicesIndexer, err := makeIndexer(refCtx, listWrapper(refCtx, client.CoreV1().Services(metav1.NamespaceAll).List))
	if err != nil {
		return DataSource{}, err
	}

	secretsIndexer, err := makeIndexer(refCtx, listWrapper(refCtx, client.CoreV1().Secrets(metav1.NamespaceAll).List))
	if err != nil {
		return DataSource{}, err
	}

	configMapsIndexer, err := makeIndexer(refCtx, listWrapper(refCtx, client.CoreV1().ConfigMaps(metav1.NamespaceAll).List))
	if err != nil {
		return DataSource{}, err
	}

	serviceAccountsIndexer, err := makeIndexer(refCtx, listWrapper(refCtx, client.CoreV1().ServiceAccounts(metav1.NamespaceAll).List))
	if err != nil {
		return DataSource{}, err
	}

	scyllaClustersIndexer, err := makeIndexer(refCtx, listWrapper(refCtx, scyllaClient.ScyllaV1().ScyllaClusters(metav1.NamespaceAll).List))
	if err != nil {
		return DataSource{}, err
	}

	return NewDataSource(
		corev1listers.NewServiceLister(servicesIndexer),
		corev1listers.NewSecretLister(secretsIndexer),
		corev1listers.NewConfigMapLister(configMapsIndexer),
		corev1listers.NewServiceAccountLister(serviceAccountsIndexer),
		scyllav1listers.NewScyllaClusterLister(scyllaClustersIndexer),
	), nil
}
