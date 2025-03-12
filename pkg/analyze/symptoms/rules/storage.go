package rules

import (
	"errors"
	"fmt"
	"github.com/scylladb/scylla-operator/pkg/analyze/selectors"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	scyllav1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"regexp"
	"strings"
)

const (
	csiDriverContainerName                     = "local-csi-driver"
	nodeConfigConditionPattern                 = "NodeSetup(?P<node>.*)Degraded"
	nodeConfigNonexistentVolumeMessageFragment = "resolve RAID device"
)

var StorageSymptoms = symptoms.NewSymptomSet("storage", []*symptoms.SymptomSet{
	buildLocalCsiDriverMissingSymptoms(),
	buildStorageClassMissingSymptoms(),
	buildNodeConfigSymptoms(),
})

func buildLocalCsiDriverMissingSymptoms() *symptoms.SymptomSet {
	// Scenario #2: local-csi-driver CSIDriver, referenced by scylladb-local-xfs StorageClass, is missing
	csiDriverMissing := symptoms.NewSymptom("CSIDriver is missing",
		"%[csi-driver.Name]% CSIDriver, referenced by %[storage-class.Name]% StorageClass, is missing",
		"deploy %[csi-driver.Name]% provisioner (or change StorageClass)",
		selectors.
			Select("scylla-cluster", selectors.Type[*scyllav1.ScyllaCluster]()).
			Select("scylla-pod", selectors.Type[*v1.Pod]()).
			Select("storage-class", selectors.Type[*storagev1.StorageClass]()).
			Select("csi-driver", selectors.Type[*storagev1.CSIDriver]()).
			Filter("scylla-cluster", func(c *scyllav1.ScyllaCluster) bool {
				if c == nil {
					return false
				}

				statefulSetControllerProgressing := false
				progressing := false
				storageClassXfs := false
				for _, rack := range c.Spec.Datacenter.Racks {
					if *rack.Storage.StorageClassName == "scylladb-local-xfs" {
						storageClassXfs = true
					}
				}
				for _, cond := range c.Status.Conditions {
					if cond.Type == "StatefulSetControllerProgressing" && cond.Status == "True" {
						statefulSetControllerProgressing = true
					} else if cond.Type == "Progressing" && cond.Status == "True" {
						progressing = true
					}
				}
				return storageClassXfs && (statefulSetControllerProgressing || progressing)
			}).
			Filter("scylla-pod", func(p *v1.Pod) bool {
				return p != nil && p.Status.Phase == v1.PodPending
			}).
			Filter("storage-class", func(s *storagev1.StorageClass) bool {
				return s != nil
			}).
			Assert("csi-driver", func(d *storagev1.CSIDriver) bool {
				return d == nil
			}).
			Relate("scylla-cluster", "scylla-pod", func(c *scyllav1.ScyllaCluster, p *v1.Pod) bool {
				return c.Name == p.Labels["scylla/cluster"]
			}).
			Relate("scylla-cluster", "storage-class", func(c *scyllav1.ScyllaCluster, sc *storagev1.StorageClass) bool {
				for _, rack := range c.Spec.Datacenter.Racks {
					if *rack.Storage.StorageClassName == sc.Name {
						return true
					}
				}
				return false
			}).
			Relate("storage-class", "csi-driver", func(sc *storagev1.StorageClass, d *storagev1.CSIDriver) bool {
				return sc.Provisioner == d.Name
			}).
			Collect(symptoms.DefaultLimit))

	csiDriverMissingSymptoms := symptoms.NewEmptySymptomSet("csi-driver-missing")
	err := csiDriverMissingSymptoms.Add(&csiDriverMissing)
	if err != nil {
		panic(errors.New("failed to create csiDriverMissing symptom" + err.Error()))
	}
	return &csiDriverMissingSymptoms
}

func buildStorageClassMissingSymptoms() *symptoms.SymptomSet {
	// Scenario #1: scylladb-local-xfs StorageClass used by a ScyllaCluster is missing
	notDeployedStorageClass := symptoms.NewSymptom("StorageClass is missing",
		"%[storage-class.Name]% StorageClass used by a ScyllaCluster is missing",
		"deploy %[storage-class.Name]% StorageClass (or change StorageClass)",
		selectors.
			Select("scylla-cluster", selectors.Type[*scyllav1.ScyllaCluster]()).
			Select("storage-class", selectors.Type[*storagev1.StorageClass]()).
			Select("scylla-pod", selectors.Type[*v1.Pod]()).
			Select("pod-pvc", selectors.Type[*v1.PersistentVolumeClaim]()).
			Select("csi-driver", selectors.Type[*storagev1.CSIDriver]()).
			Filter("scylla-cluster", func(c *scyllav1.ScyllaCluster) bool {
				if c == nil {
					return false
				}

				statefulSetControllerProgressing := false
				progressing := false
				storageClassXfs := false
				for _, rack := range c.Spec.Datacenter.Racks {
					if *rack.Storage.StorageClassName == "scylladb-local-xfs" {
						storageClassXfs = true
					}
				}
				for _, cond := range c.Status.Conditions {
					if cond.Type == "StatefulSetControllerProgressing" && cond.Status == "True" {
						statefulSetControllerProgressing = true
					} else if cond.Type == "Progressing" && cond.Status == "True" {
						progressing = true
					}
				}
				return storageClassXfs && (statefulSetControllerProgressing || progressing)
			}).
			Filter("scylla-pod", func(p *v1.Pod) bool {
				if p == nil {
					return false
				}

				podScheduled := true
				unboundPVC := false
				for _, cond := range p.Status.Conditions {
					if cond.Type == "PodScheduled" {
						if strings.Contains(cond.Message, "pod has unbound immediate PersistentVolumeClaims") {
							unboundPVC = true
						}
						if cond.Status == "False" {
							podScheduled = false
						}
					}
				}
				return !podScheduled && unboundPVC
			}).
			Assert("storage-class", func(s *storagev1.StorageClass) bool {
				return s == nil
			}).
			Filter("csi-driver", func(d *storagev1.CSIDriver) bool {
				return d != nil
			}).
			Filter("pod-pvc", func(pvc *v1.PersistentVolumeClaim) bool {
				return pvc != nil && pvc.Status.Phase == v1.ClaimPending
			}).
			Relate("scylla-cluster", "scylla-pod", func(c *scyllav1.ScyllaCluster, p *v1.Pod) bool {
				return c.Name == p.Labels["scylla/cluster"]
			}).
			Relate("scylla-cluster", "storage-class", func(c *scyllav1.ScyllaCluster, sc *storagev1.StorageClass) bool {
				for _, rack := range c.Spec.Datacenter.Racks {
					if *rack.Storage.StorageClassName == sc.Name {
						return true
					}
				}
				return false
			}).
			Relate("storage-class", "csi-driver", func(sc *storagev1.StorageClass, d *storagev1.CSIDriver) bool {
				return sc.Provisioner == d.Name
			}).
			Relate("pod-pvc", "scylla-pod", func(pvc *v1.PersistentVolumeClaim, p *v1.Pod) bool {
				for _, volume := range p.Spec.Volumes {
					vPvc := volume.PersistentVolumeClaim
					if vPvc != nil && (*vPvc).ClaimName == pvc.Name {
						return true
					}
				}
				return false
			}).
			Relate("pod-pvc", "storage-class", func(pvc *v1.PersistentVolumeClaim, sc *storagev1.StorageClass) bool {
				return pvc.Spec.StorageClassName != nil && *pvc.Spec.StorageClassName == sc.Name
			}).
			Collect(symptoms.DefaultLimit))

	storageClassMissingSymptoms := symptoms.NewEmptySymptomSet("StorageClass missing")
	err := storageClassMissingSymptoms.Add(&notDeployedStorageClass)
	if err != nil {
		panic(errors.New("failed to create storageClassMissingSymptoms symptom" + err.Error()))
	}
	return &storageClassMissingSymptoms
}

func buildNodeConfigSymptoms() *symptoms.SymptomSet {
	// Scenario #4: E2Es for misconfigured NodeConfigs
	nodeConfigClusterWideNonexistentVolume := symptoms.NewSymptomWithManyDiagSug(
		"NodeConfig doesn't provision storage for local-csi-driver's volumes-dir",
		[]string{
			"driver is unable to provide storage due to NodeConfig being misconfigured",
			"this may be false-positive in clusters which don't use NodeConfig to mount `volumes-dir` of local-csi-drivers' Pods",
		},
		[]string{
			"fix the NodeConfig (see node-config's #TODO: node config path# error conditions), then rolling restart local-csi-driver",
		},
		selectors.
			Select("scylla-cluster", selectors.Type[*scyllav1.ScyllaCluster]()).
			Select("scylla-pod", selectors.Type[*v1.Pod]()).
			Select("pod-pvc", selectors.Type[*v1.PersistentVolumeClaim]()).
			// TODO: 	sprawdzic wszystkie storage class'y
			Select("storage-class", selectors.Type[*storagev1.StorageClass]()).
			Select("csi-driver", selectors.Type[*storagev1.CSIDriver]()).
			Select("scylla-node", selectors.Type[*v1.Node]()).
			Select("node-config", selectors.Type[*scyllav1alpha1.NodeConfig]()).
			Select("csi-pod", selectors.Type[*v1.Pod]()).
			Filter("scylla-cluster", func(c *scyllav1.ScyllaCluster) bool {
				return c != nil && !symptoms.MeetsCondition(c, "Available", "True")
			}).
			Filter("scylla-pod", func(p *v1.Pod) bool {
				return p != nil && p.Status.Phase != v1.PodSucceeded
			}).
			Filter("pod-pvc", func(pvc *v1.PersistentVolumeClaim) bool {
				return pvc != nil && pvc.Status.Phase == v1.ClaimPending && pvc.Spec.StorageClassName != nil
			}).
			Filter("storage-class", func(s *storagev1.StorageClass) bool {
				return s != nil
			}).
			Filter("csi-driver", func(d *storagev1.CSIDriver) bool {
				return d != nil
			}).
			Filter("scylla-node", func(n *v1.Node) bool {
				return n != nil
			}).
			Filter("node-config", func(nc *scyllav1alpha1.NodeConfig) bool {
				return nc != nil
			}).
			Filter("csi-pod", func(d *v1.Pod) bool {
				if d == nil {
					return false
				}
				for _, cont := range d.Spec.Containers {
					if cont.Name == csiDriverContainerName {
						return true
					}
				}
				return false
			}).
			Relate("scylla-cluster", "scylla-pod", func(c *scyllav1.ScyllaCluster, p *v1.Pod) bool {
				return c.Name == p.Labels["scylla/cluster"]
			}).
			Relate("scylla-pod", "pod-pvc", func(p *v1.Pod, pvc *v1.PersistentVolumeClaim) bool {
				for _, volume := range p.Spec.Volumes {
					if volume.PersistentVolumeClaim != nil {
						volumePvc := *volume.PersistentVolumeClaim
						if volumePvc.ClaimName == pvc.Name {
							return true
						}
					}
				}
				return false
			}).
			Relate("pod-pvc", "storage-class", func(pvc *v1.PersistentVolumeClaim, s *storagev1.StorageClass) bool {
				return *pvc.Spec.StorageClassName == s.Name
			}).
			Relate("storage-class", "csi-driver", func(s *storagev1.StorageClass, csi *storagev1.CSIDriver) bool {
				return s.Provisioner == csi.Name
			}).
			Relate("scylla-node", "scylla-cluster", func(n *v1.Node, c *scyllav1.ScyllaCluster) bool {
				// TODO: relate the node to the cluster (with node placement rules?)
				return true
			}).
			Relate("node-config", "scylla-node", func(nc *scyllav1alpha1.NodeConfig, n *v1.Node) bool {
				return symptoms.MeetsNodeSelectorPlacementRules(n, nc.Spec.Placement.NodeSelector)
			}).
			// node-config tunes the node on which csi-pod is placed
			Relate("csi-pod", "node-config", func(csiPod *v1.Pod, nodeConfig *scyllav1alpha1.NodeConfig) bool {
				matchesPod := false
				dirsMounted := false
				for _, status := range nodeConfig.Status.NodeStatuses {
					if status.Name == csiPod.Spec.NodeName {
						matchesPod = true
					}
				}
				if !matchesPod {
					return false
				}
				volumesDir := ""
				volumesMountPoint := ""
				for _, container := range csiPod.Spec.Containers {
					csiContainer := false
					for _, arg := range container.Args {
						if strings.Contains(arg, "--volumes-dir=") {
							volumesDir = strings.Split(arg, "--volumes-dir=")[1]
							volumesDir = strings.TrimSpace(volumesDir)
							csiContainer = true
						}
					}
					if csiContainer {
						for _, mnt := range container.VolumeMounts {
							if mnt.MountPath == volumesDir {
								volumesMountPoint = mnt.Name
							}
						}
					}
				}

				if volumesMountPoint != "" {
					for _, vol := range csiPod.Spec.Volumes {
						if vol.Name == volumesMountPoint {
							for _, mp := range nodeConfig.Spec.LocalDiskSetup.Mounts {
								if (*vol.HostPath).Path == mp.MountPoint {
									dirsMounted = true
								}
							}
						}
					}
				}

				return !dirsMounted
			}).
			Collect(symptoms.DefaultLimit))

	// Scenario #4': Detects the non-existent device just by condition logs
	nodeConfigNonexistentDevice := symptoms.NewSymptom(
		"NodeConfig non-existent device condition",
		"NodeConfig configured with a non-existent device",
		"fix NodeConfig",
		selectors.
			Select("node-config", selectors.Type[*scyllav1alpha1.NodeConfig]()).
			Filter("node-config", func(nc *scyllav1alpha1.NodeConfig) bool {
				if nc == nil {
					return false
				}
				r := regexp.MustCompile(nodeConfigConditionPattern)
				for _, cond := range nc.Status.Conditions {
					t := fmt.Sprintf("%v", cond)
					match := r.FindStringSubmatch(t)
					if len(match) == 0 {
						continue
					}
					if cond.Status != v1.ConditionTrue {
						continue
					}
					if strings.Contains(cond.Message, nodeConfigNonexistentVolumeMessageFragment) {
						return true
					}
				}
				return false
			}).
			Collect(symptoms.DefaultLimit))

	symptomSet := symptoms.NewEmptySymptomSet("NodeConfig errors")
	symptomSet.MustAdd(&nodeConfigClusterWideNonexistentVolume)
	symptomSet.MustAdd(&nodeConfigNonexistentDevice)
	return &symptomSet
}
