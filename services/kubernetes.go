package services

import (
	"context"
	"encoding/json"
	"errors"
	"main/lib"
	"main/models"
	"main/repository"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubernetesService struct {
	logger     lib.Logger
	Repository repository.KubernetesRepository
}

func NewKubernetesService(logger lib.Logger, Repository repository.KubernetesRepository) KubernetesService {
	return KubernetesService{
		logger:     logger,
		Repository: Repository,
	}
}

func (u KubernetesService) CreatePod(podBody models.PodBody) (*corev1.Pod, error) {

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
	pod, err := u.Repository.Clientset.CoreV1().Pods(podBody.Namespace).Create(context.TODO(), newpod, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (u KubernetesService) CreateOrUpdateConfigMap(Map models.ConfigMapBody) (corev1.ConfigMap, error) {
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Map.Name,
			Namespace: Map.Namespace,
		},
		Data: Map.Data,
	}
	if _, err := u.Repository.Clientset.CoreV1().ConfigMaps(Map.Namespace).Get(context.TODO(), Map.Name, metav1.GetOptions{}); err != nil {
		_, err := u.Repository.Clientset.CoreV1().ConfigMaps(Map.Namespace).Create(context.Background(), &cm, metav1.CreateOptions{})
		if err != nil {
			return cm, err
		}
	} else {
		u.Repository.Clientset.CoreV1().ConfigMaps(Map.Namespace).Update(context.Background(), &cm, metav1.UpdateOptions{})
	}
	return cm, nil
}

func (u KubernetesService) CreateOrUpdateSecret(s models.SecretBody) (*corev1.Secret, error) {

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
	if _, err := u.Repository.Clientset.CoreV1().Secrets(s.Namespace).Get(context.TODO(), s.Name, metav1.GetOptions{}); err != nil {
		_, err := u.Repository.Clientset.CoreV1().Secrets(s.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
		if err != nil {
			return secret, err
		}
	} else {
		u.Repository.Clientset.CoreV1().Secrets(s.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
	}
	return secret, nil
}

func (u KubernetesService) CreateNamespace(name string) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	created_namespace, err := u.Repository.Clientset.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	} else {
		return created_namespace, nil
	}
}

func (u KubernetesService) CreateNodePort(nodeport models.Nodeport) (*corev1.Service, error) {
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
	service, err := u.Repository.Clientset.CoreV1().Services(nodeport.Namespace).Create(context.TODO(), nport, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (u KubernetesService) CreateServiceAccount(serviceAccount models.ServiceAccount) (*corev1.ServiceAccount, error) {

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
	service_account, err := u.Repository.Clientset.CoreV1().ServiceAccounts("default").Create(context.Background(), serviceacc, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return service_account, nil

}

func (u KubernetesService) CreateRole(rolebody models.Role) (*rbacv1.Role, error) {
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
	r, err := u.Repository.Clientset.RbacV1().Roles(rolebody.Namespace).Create(context.Background(), role, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (u KubernetesService) CreateRoleBinding(rolebinding models.RoleBinding) (*rbacv1.RoleBinding, error) {
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
	rb, err := u.Repository.Clientset.RbacV1().RoleBindings(rolebinding.Namespace).Create(context.Background(), rbind, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}
	return rb, nil
}

func (u KubernetesService) GetCurrentPodStatusRequest(pod_name string) []byte {

	pod, err := u.Repository.Clientset.CoreV1().Pods("default").Get(context.Background(), pod_name, metav1.GetOptions{})
	if err != nil {
		return nil
	}
	status, err := json.Marshal(pod.Status)

	if err != nil {
		return nil
	}
	return status
}

func (u KubernetesService) DeletePod(name, namespace string) error {
	err := u.Repository.Clientset.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (u KubernetesService) GetPodsList(namespace string) ([]string, error) {
	pods, err := u.Repository.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	var names []string
	if err != nil || len(pods.Items) == 0 {
		return nil, errors.New("Not found")
	} else {
		names = make([]string, len(pods.Items))

		for i := 0; i < len(pods.Items); i++ {
			names[i] = pods.Items[i].Name
		}
	}
	return names, nil
}
