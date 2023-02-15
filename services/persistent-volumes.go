package services

import (
	"context"
	"main/models"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (u KubernetesService) CreatePersistentVolume(persistentVolume models.PV) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   persistentVolume.Name,
			Labels: map[string]string{"type": "local"},
		},
		Spec: corev1.PersistentVolumeSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(persistentVolume.Storage),
			},
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
			StorageClassName:              "",
			PersistentVolumeSource: corev1.PersistentVolumeSource{HostPath: &corev1.HostPathVolumeSource{
				Path: persistentVolume.Path, //"/tmp/mongodb"
			}},
		},
	}
	Pv, err := u.Repository.Clientset.CoreV1().PersistentVolumes().Create(context.Background(), pv, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return Pv, nil

}

func (u KubernetesService) CreatePersistentVolumeClaim(persistentVolumeClaim models.PVC) (*corev1.PersistentVolumeClaim, error) {
	storageclass := ""
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      persistentVolumeClaim.Name,
			Namespace: persistentVolumeClaim.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			StorageClassName: &storageclass,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(persistentVolumeClaim.Storage),
				},
			},
		},
	}
	pvc, err := u.Repository.Clientset.CoreV1().PersistentVolumeClaims(persistentVolumeClaim.Namespace).Create(context.Background(), pvc, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return pvc, nil

}
