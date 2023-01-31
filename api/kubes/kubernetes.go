package kubes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// func GetEnvPort(podbody PodBody) int32 {
// 	var Intport int32
// 	port, _ := strconv.ParseInt(os.Getenv("KUBE_PORT"), 10, 32)
// 	Intport = int32(port)
// 	for i := 0; i < len(podbody.Envs); i++ {
// 		if podbody.Envs[i].Name == "PORT" {
// 			port, err := strconv.ParseInt(podbody.Envs[i].Value, 10, 32)
// 			if err != nil {
// 				panic(err.Error())
// 			}
// 			Intport = int32(port)
// 		}
// 	}
// 	return Intport
// }

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
					Ports: []corev1.ContainerPort{
						{
							HostPort:      podBody.Port,
							ContainerPort: podBody.Port,
							Name:          "http",
						},
					},
					EnvFrom: []corev1.EnvFromSource{
						{
							ConfigMapRef: &corev1.ConfigMapEnvSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: podBody.ConfigmapName,
								},
							},
						},
						{
							SecretRef: &corev1.SecretEnvSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: podBody.SecretName,
								},
							},
						},
					},
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

func (u KubeRequest) CreateOrUpdateConfigMap(Map ConfigMapBody) (corev1.ConfigMap, error) {
	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Map.Name,
			Namespace: Map.Namespace,
		},
		Data: Map.Data,
	}
	if _, err := u.clientset.CoreV1().ConfigMaps(Map.Namespace).Get(context.TODO(), Map.Name, metav1.GetOptions{}); err != nil {
		u.clientset.CoreV1().ConfigMaps(Map.Namespace).Create(context.Background(), &cm, metav1.CreateOptions{})
	} else {
		u.clientset.CoreV1().ConfigMaps(Map.Namespace).Update(context.Background(), &cm, metav1.UpdateOptions{})
	}
	return cm, nil
}

func (u KubeRequest) CreateOrUpdateSecret(s SecretBody) (*corev1.Secret, error) {

	data := make(map[string][]byte)
	for key, value := range s.Data {
		data[key] = []byte(value)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Data: data,
		// StringData: map[string][]byte{
		// 	"SecretKey": []byte("Secret key is python"),
		// },
	}
	if _, err := u.clientset.CoreV1().Secrets(s.Namespace).Get(context.TODO(), s.Name, metav1.GetOptions{}); err != nil {
		u.clientset.CoreV1().Secrets(s.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	} else {
		u.clientset.CoreV1().Secrets(s.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
	}
	return secret, nil
}
