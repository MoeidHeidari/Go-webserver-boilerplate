package kubes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (u KubeRequest) CreatePod(podBody PodBody) (*corev1.Pod, error) {

	newpod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podBody.Name,
			Namespace: podBody.Namespace,
			Labels: map[string]string{
				"kubernetes.io/hostname": "minikube",
			},
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: podBody.VolumeName,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: podBody.ClaimName,
						},
					},
				},
			},
			NodeSelector: map[string]string{
				"kubernetes.io/hostname": "minikube",
			},
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
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      podBody.VolumeName,
							MountPath: podBody.MountPath, //"/usr/share/mongo"
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

	pod, err := u.Clientset.CoreV1().Pods(podBody.Namespace).Create(context.TODO(), newpod, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (u KubeRequest) CreateOrUpdateConfigMap(Map ConfigMapBody) (corev1.ConfigMap, error) {
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Map.Name,
			Namespace: Map.Namespace,
		},
		Data: Map.Data,
	}
	if _, err := u.Clientset.CoreV1().ConfigMaps(Map.Namespace).Get(context.TODO(), Map.Name, metav1.GetOptions{}); err != nil {
		_, err := u.Clientset.CoreV1().ConfigMaps(Map.Namespace).Create(context.Background(), &cm, metav1.CreateOptions{})
		if err != nil {
			return cm, err
		}
	} else {
		u.Clientset.CoreV1().ConfigMaps(Map.Namespace).Update(context.Background(), &cm, metav1.UpdateOptions{})
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
	}
	if _, err := u.Clientset.CoreV1().Secrets(s.Namespace).Get(context.TODO(), s.Name, metav1.GetOptions{}); err != nil {
		_, err := u.Clientset.CoreV1().Secrets(s.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
		if err != nil {
			return secret, err
		}
	} else {
		u.Clientset.CoreV1().Secrets(s.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
	}
	return secret, nil
}

func (u KubeRequest) CreateNamespace(name string) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	created_namespace, err := u.Clientset.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	} else {
		return created_namespace, nil
	}
}

func (u KubeRequest) CreateNodePort(nodeport Nodeport) (*corev1.Service, error) {
	nport := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nodeport.Name,
			Namespace: nodeport.Namespace,
			Labels: map[string]string{
				"kubernetes.io/hostname": "minikube",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: nodeport.Port, NodePort: nodeport.RedirectPort,
				},
			},
			Type: corev1.ServiceTypeNodePort,
			Selector: map[string]string{
				"kubernetes.io/hostname": "minikube",
			},
		},
	}
	service, err := u.Clientset.CoreV1().Services(nodeport.Namespace).Create(context.TODO(), nport, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (u KubeRequest) CreateServiceAccount(serviceAccount ServiceAccount) (*corev1.ServiceAccount, error) {

	serviceacc := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccount.Name,
			Namespace: serviceAccount.Namespace,
		},
		Secrets: []corev1.ObjectReference{
			{Kind: "secret",
				Namespace: serviceAccount.SecretNamespace,
				Name:      serviceAccount.SecretName,
			},
		},
	}
	service_account, err := u.Clientset.CoreV1().ServiceAccounts("default").Create(context.Background(), serviceacc, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return service_account, nil

}

func (u KubeRequest) CreateRole(rolebody Role) (*rbacv1.Role, error) {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rolebody.Name,
			Namespace: rolebody.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{Verbs: rolebody.Verbs,
				Resources: rolebody.Resources,
				APIGroups: []string{
					"",
				},
			},
		},
	}
	r, err := u.Clientset.RbacV1().Roles(rolebody.Namespace).Create(context.Background(), role, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (u KubeRequest) CreateRoleBinding(rolebinding RoleBinding) (*rbacv1.RoleBinding, error) {
	rbind := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rolebinding.Name,
			Namespace: rolebinding.Namespace,
		},
		Subjects: []rbacv1.Subject{
			{Kind: "ServiceAccount",
				Name: rolebinding.AccountName},
		},
		RoleRef: rbacv1.RoleRef{
			Name: rolebinding.RoleName,
			Kind: "Role",
		},
	}
	rb, err := u.Clientset.RbacV1().RoleBindings(rolebinding.Namespace).Create(context.Background(), rbind, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}
	return rb, nil
}
