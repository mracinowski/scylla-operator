package analyze

import (
	"encoding/json"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
)

func CheckOperatorPodsExist(ds *DataSource) (bool, error) {
	sel, err := labels.Parse("app.kubernetes.io/name == scylla-operator")
	if err != nil {
		return false, err
	}

	list, err := ds.PodLister.List(sel)
	if err != nil {
		return false, err
	}

	return len(list) > 0, nil
}

func CheckStorageClassMissing(ds *DataSource) (bool, error) {
	list, err := ds.ScyllaClusterLister.List(labels.Everything())
	if err != nil {
		return false, err
	}

	for _, cluster := range list {
		for _, condition := range cluster.Status.Conditions {
			if condition.Status != "True" {
				continue
			}

			if condition.Type != "StatefulSetControllerProgressing" {
				continue
			}

			if condition.Reason != "StatefulSetControllerProgressing_WaitingForStatefulSetRollout" {
				continue
			}

			s, err := json.MarshalIndent(condition, "", "\t")
			if err != nil {
				return false, err
			}

			klog.Infof("%s", s)

			return true, nil
		}

	}

	return false, nil
}
