package snapshot

import (
	"context"
	"github.com/google/go-cmp/cmp"
	scyllav1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1"
	"github.com/scylladb/scylla-operator/pkg/client/scylla/clientset/versioned/fake"
	"github.com/scylladb/scylla-operator/pkg/helpers/slices"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubefake "k8s.io/client-go/kubernetes/fake"
	"reflect"
	"sort"
	"testing"
)

func TestNewDataSourceFromClients_SingleLister(t *testing.T) {
	t.Parallel()
	newDataSourceTests := []struct {
		name              string
		kubernetesObjects []runtime.Object
		scyllaObjects     []runtime.Object
		checkedType       reflect.Type
		expectedObjects   []any
		expectedErr       error
	}{
		{
			name: "empty pod list",
			kubernetesObjects: func() []runtime.Object {
				return slices.FilterOut(newBasicKubernetesObjects(), func(obj runtime.Object) bool {
					_, ok := obj.(*corev1.Pod)
					return ok
				})
			}(),
			scyllaObjects:   newBasicScyllaObjects(),
			expectedObjects: []any{},
			expectedErr:     nil,
		},
		{
			name:              "nonempty pod list",
			kubernetesObjects: newBasicKubernetesObjects(),
			scyllaObjects:     newBasicScyllaObjects(),
			expectedObjects: []any{
				&corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod1",
						Namespace: "test",
					},
				},
				&corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod2",
						Namespace: "test",
					},
				},
			},
			checkedType: reflect.TypeOf(&corev1.Pod{}),
			expectedErr: nil,
		},
		{
			name: "empty service list",
			kubernetesObjects: func() []runtime.Object {
				return slices.FilterOut(newBasicKubernetesObjects(), func(obj runtime.Object) bool {
					_, ok := obj.(*corev1.Service)
					return ok
				})
			}(),
			scyllaObjects:   newBasicScyllaObjects(),
			expectedObjects: []any{},
			expectedErr:     nil,
		},
		{
			name:              "nonempty service list",
			kubernetesObjects: newBasicKubernetesObjects(),
			scyllaObjects:     newBasicScyllaObjects(),
			expectedObjects: []any{
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service1",
						Namespace: "test",
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service2",
						Namespace: "test",
					},
				},
			},
			checkedType: reflect.TypeOf(&corev1.Service{}),
			expectedErr: nil,
		},
		{
			name: "empty secret list",
			kubernetesObjects: func() []runtime.Object {
				return slices.FilterOut(newBasicKubernetesObjects(), func(obj runtime.Object) bool {
					_, ok := obj.(*corev1.Secret)
					return ok
				})
			}(),
			scyllaObjects:   newBasicScyllaObjects(),
			expectedObjects: []any{},
			expectedErr:     nil,
		},
		{
			name:              "nonempty secret list",
			kubernetesObjects: newBasicKubernetesObjects(),
			scyllaObjects:     newBasicScyllaObjects(),
			expectedObjects: []any{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret1",
						Namespace: "test",
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret2",
						Namespace: "test",
					},
				},
			},
			checkedType: reflect.TypeOf(&corev1.Secret{}),
			expectedErr: nil,
		},
		{
			name: "empty config map list",
			kubernetesObjects: func() []runtime.Object {
				return slices.FilterOut(newBasicKubernetesObjects(), func(obj runtime.Object) bool {
					_, ok := obj.(*corev1.ConfigMap)
					return ok
				})
			}(),
			scyllaObjects:   newBasicScyllaObjects(),
			expectedObjects: []any{},
			expectedErr:     nil,
		},
		{
			name:              "nonempty config map list",
			kubernetesObjects: newBasicKubernetesObjects(),
			scyllaObjects:     newBasicScyllaObjects(),
			expectedObjects: []any{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "configmap1",
						Namespace: "test",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "configmap2",
						Namespace: "test",
					},
				},
			},
			checkedType: reflect.TypeOf(&corev1.ConfigMap{}),
			expectedErr: nil,
		},
		{
			name: "empty service account list",
			kubernetesObjects: func() []runtime.Object {
				return slices.FilterOut(newBasicKubernetesObjects(), func(obj runtime.Object) bool {
					_, ok := obj.(*corev1.ServiceAccount)
					return ok
				})
			}(),
			scyllaObjects:   newBasicScyllaObjects(),
			expectedObjects: []any{},
			expectedErr:     nil,
		},
		{
			name:              "nonempty service account list",
			kubernetesObjects: newBasicKubernetesObjects(),
			scyllaObjects:     newBasicScyllaObjects(),
			expectedObjects: []any{
				&corev1.ServiceAccount{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "serviceaccount1",
						Namespace: "test",
					},
				},
				&corev1.ServiceAccount{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "serviceaccount2",
						Namespace: "test",
					},
				},
			},
			checkedType: reflect.TypeOf(&corev1.ServiceAccount{}),
			expectedErr: nil,
		},
		{
			name:              "empty scylla cluster list",
			kubernetesObjects: []runtime.Object{},
			scyllaObjects: func() []runtime.Object {
				return slices.FilterOut(newBasicScyllaObjects(), func(obj runtime.Object) bool {
					_, ok := obj.(*scyllav1.ScyllaCluster)
					return ok
				})
			}(),
			expectedObjects: []any{},
			expectedErr:     nil,
		},
		{
			name:              "nonempty scylla cluster list",
			kubernetesObjects: newBasicKubernetesObjects(),
			scyllaObjects:     newBasicScyllaObjects(),
			expectedObjects: []any{
				&scyllav1.ScyllaCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "scyllacluster1",
						Namespace: "test",
					},
				},
				&scyllav1.ScyllaCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "scyllacluster2",
						Namespace: "test",
					},
				},
			},
			checkedType: reflect.TypeOf(&scyllav1.ScyllaCluster{}),
			expectedErr: nil,
		},
	}

	for _, tc := range newDataSourceTests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			fakeClient := kubefake.NewSimpleClientset(tc.kubernetesObjects...)
			scyllaFakeClient := fake.NewSimpleClientset(tc.scyllaObjects...)

			snapshot, err := NewDataSourceFromClients(ctx, fakeClient, scyllaFakeClient)

			objects := snapshot.List(tc.checkedType)

			sort.Slice(objects, func(i, j int) bool {
				return compareRuntimeObjects(objects[i], objects[j])
			})
			sort.Slice(tc.expectedObjects, func(i, j int) bool {
				return compareRuntimeObjects(tc.expectedObjects[i], tc.expectedObjects[j])
			})
			if !reflect.DeepEqual(err, tc.expectedErr) {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, err)
			}

			if !equality.Semantic.DeepEqual(objects, tc.expectedObjects) {
				t.Errorf("expected and actual objects differ: %s", cmp.Diff(tc.expectedObjects, objects))
			}
		})
	}
}

func compareRuntimeObjects(a interface{}, b interface{}) bool {
	aObj, aOK := a.(metav1.Object)
	bObj, bOK := b.(metav1.Object)
	if !aOK || !bOK {
		panic("can't cast indexer object to metav1.Object")
	}
	valueA := reflect.ValueOf(aObj).Elem().FieldByName("ObjectMeta")
	valueB := reflect.ValueOf(bObj).Elem().FieldByName("ObjectMeta")
	metaA := valueA.Interface().(metav1.ObjectMeta)
	metaB := valueB.Interface().(metav1.ObjectMeta)
	return metaA.Namespace+metaA.Name < metaB.Namespace+metaB.Name
}

func newBasicKubernetesObjects() []runtime.Object {
	return []runtime.Object{
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod1",
				Namespace: "test",
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod2",
				Namespace: "test",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service1",
				Namespace: "test",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service2",
				Namespace: "test",
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret1",
				Namespace: "test",
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret2",
				Namespace: "test",
			},
		},
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "configmap1",
				Namespace: "test",
			},
		},
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "configmap2",
				Namespace: "test",
			},
		},
		&corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "serviceaccount1",
				Namespace: "test",
			},
		},
		&corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "serviceaccount2",
				Namespace: "test",
			},
		},
	}
}

func newBasicScyllaObjects() []runtime.Object {
	return []runtime.Object{
		&scyllav1.ScyllaCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "scyllacluster1",
				Namespace: "test",
			},
		},
		&scyllav1.ScyllaCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "scyllacluster2",
				Namespace: "test",
			},
		},
	}
}
