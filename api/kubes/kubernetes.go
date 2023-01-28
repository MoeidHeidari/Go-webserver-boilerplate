package kubes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (u KubeRequest) CreatePod(podBody PodBody) (*corev1.Pod, error) {
	newpod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podBody.Name,
			Namespace: podBody.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    podBody.ContainerName,
					Image:   podBody.Image,
					Command: podBody.Command,
				},
			},
		},
	}

	pod, err := u.clientset.CoreV1().Pods(podBody.Namespace).Create(context.TODO(), newpod, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return pod, nil
}
