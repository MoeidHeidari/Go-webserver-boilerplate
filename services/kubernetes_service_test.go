package services_test

import (
	"context"
	"log"
	"main/lib"
	"main/models"
	"main/repository"
	"main/services"
	"math/rand"
	"os"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/action"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type test struct {
	ConfigmapName      string
	SecretName         string
	Namespace          string
	NodeportName       string
	PodName            string
	PVName             string
	PVCName            string
	ChartName          string
	RoleName           string
	RoleBindingName    string
	ServiceAccountName string
}

var Test = test{
	ConfigmapName:      faker.Word(),
	SecretName:         faker.Word(),
	Namespace:          faker.Word(),
	NodeportName:       faker.Word(),
	PodName:            faker.Word(),
	PVName:             faker.Word(),
	PVCName:            faker.Word(),
	ChartName:          faker.Word(),
	RoleName:           faker.Word(),
	RoleBindingName:    faker.Word(),
	ServiceAccountName: faker.Word(),
}

func DeleteAll() {
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	err := u.Repository.Clientset.CoreV1().Pods("default").Delete(context.TODO(), Test.PodName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Repository.Clientset.CoreV1().PersistentVolumeClaims("default").Delete(context.TODO(), Test.PVCName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Repository.Clientset.CoreV1().PersistentVolumes().Delete(context.TODO(), Test.PVName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Repository.Clientset.CoreV1().Services("default").Delete(context.Background(), Test.NodeportName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Repository.Clientset.CoreV1().Secrets("default").Delete(context.Background(), Test.SecretName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Repository.Clientset.CoreV1().Namespaces().Delete(context.Background(), Test.Namespace, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Repository.Clientset.CoreV1().ConfigMaps("default").Delete(context.Background(), Test.ConfigmapName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Repository.Clientset.RbacV1().Roles("default").Delete(context.Background(), Test.RoleName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Repository.Clientset.RbacV1().RoleBindings("default").Delete(context.Background(), Test.RoleBindingName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Repository.Clientset.CoreV1().ServiceAccounts("default").Delete(context.Background(), Test.ServiceAccountName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	pvlist, err := u.Repository.Clientset.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, i := range pvlist.Items {
		if i.Spec.ClaimRef.Name == ("data-" + Test.ChartName + "-postgresql-0") {
			err = u.Repository.Clientset.CoreV1().PersistentVolumeClaims("default").Delete(context.TODO(), i.Spec.ClaimRef.Name, metav1.DeleteOptions{})
			if err != nil {
				panic(err.Error())
			}
			err = u.Repository.Clientset.CoreV1().PersistentVolumes().Delete(context.TODO(), i.Name, metav1.DeleteOptions{})
			if err != nil {
				panic(err.Error())
			}

		}
	}
	u.Repository.ActionConfiguration.Init(u.Repository.Settings.RESTClientGetter(), u.Repository.Settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf)
	client := action.NewUninstall(u.Repository.ActionConfiguration)
	client.DisableHooks = true
	client.Run(Test.ChartName)
}

func TestConfigMapCreate(t *testing.T) {
	configmapbody := models.ConfigMapBody{}
	err := faker.FakeData(&configmapbody.Data)
	if err != nil {
		panic(err.Error())
	}
	configmapbody.Name = Test.ConfigmapName
	configmapbody.Namespace = "default"
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	cm, err := u.CreateOrUpdateConfigMap(configmapbody)
	assert.Nil(t, err)
	assert.NotNil(t, cm)
	get_cm, err := u.Repository.Clientset.CoreV1().ConfigMaps("default").Get(context.Background(), configmapbody.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, get_cm.Labels, cm.Labels)

}

func TestNamespaceCreate(t *testing.T) {
	namespace_name := Test.Namespace
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	namespace, err := u.CreateNamespace(namespace_name)
	assert.Nil(t, err)
	assert.Equal(t, namespace_name, namespace.Name)
}

func TestSercetCreate(t *testing.T) {
	secretbody := models.SecretBody{}
	err := faker.FakeData(&secretbody.Data)

	if err != nil {
		panic(err.Error())
	}

	secretbody.Name = Test.SecretName
	secretbody.Namespace = "default"
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	s, err := u.CreateOrUpdateSecret(secretbody)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	get_secretbody, err := u.Repository.Clientset.CoreV1().Secrets("default").Get(context.Background(), secretbody.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, get_secretbody.Labels, s.Labels)
}

func TestNodePortCreate(t *testing.T) {
	nodeportbody := models.Nodeport{}
	nodeportbody.Name = Test.NodeportName
	nodeportbody.Namespace = "default"
	nodeportbody.RedirectPort = int32(rand.Intn(32767-30000) + 30000)
	nodeportbody.Port = int32(rand.Intn(30000-20000) + 20000)
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	service, err := u.CreateNodePort(nodeportbody)
	assert.Nil(t, err)
	assert.NotNil(t, service)
	checkservice, err := u.Repository.Clientset.CoreV1().Services("default").Get(context.Background(), nodeportbody.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, checkservice.UID, service.UID)

}

func TestPodCreate(t *testing.T) {
	podBody := models.PodBody{}
	podBody.Name = Test.PodName
	podBody.Namespace = "default"
	podBody.ClaimName = faker.Word()
	podBody.VolumeName = faker.Word()
	podBody.SecretName = faker.Word()
	podBody.ContainerName = faker.Word()
	podBody.Image = faker.Word()
	podBody.Port = int32(rand.Intn(30000-20000) + 20000)
	podBody.MountPath = faker.URL()
	podBody.ConfigmapName = faker.Word()
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	pod, err := u.CreatePod(podBody)
	assert.Nil(t, err)
	assert.NotNil(t, pod)
	checkpod, err := u.Repository.Clientset.CoreV1().Pods("default").Get(context.TODO(), podBody.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, checkpod.UID, pod.UID)

}

func TestPVCreate(t *testing.T) {
	pv := models.PV{}
	pv.Name = Test.PVName
	pv.Path = faker.Word()
	pv.Storage = "1Gi"
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	persistent_volume, err := u.CreatePersistentVolume(pv)
	assert.Nil(t, err)
	pvcheck, err := u.Repository.Clientset.CoreV1().PersistentVolumes().Get(context.TODO(), pv.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, pvcheck.UID, persistent_volume.UID)
}

func TestPVCCreate(t *testing.T) {
	pvc := models.PVC{}
	pvc.Name = Test.PVCName
	pvc.Namespace = "default"
	pvc.Storage = "1Gi"
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	persistent_volume_claim, err := u.CreatePersistentVolumeClaim(pvc)
	assert.Nil(t, err)
	assert.NotNil(t, persistent_volume_claim)
	pvccheck, err := u.Repository.Clientset.CoreV1().PersistentVolumeClaims("default").Get(context.TODO(), pvc.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, pvccheck.UID, persistent_volume_claim.UID)
}

func TestGetEvents(t *testing.T) {
	namespace := "default"
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	events, err := u.GetEvents(namespace)
	assert.NotNil(t, events)
	assert.Nil(t, err)
}

func TestGetCurrentPodStatus(t *testing.T) {
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	resp := u.GetCurrentPodStatusRequest("mongo")
	assert.NotNil(t, resp)
}

func TestHCreateRelease(t *testing.T) {
	chart := models.ChartBody{}
	chart.Namespace = "default"
	chart.ReleaseName = Test.ChartName
	chart.ChartPath = "https://charts.bitnami.com/bitnami/keycloak-13.0.2.tgz"
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	release, err := u.HCreateRelease(chart)
	if err != nil {
		panic(err.Error())
	}
	assert.Nil(t, err)
	assert.NotNil(t, release)

}

func TestHGetRelease(t *testing.T) {
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	results, err := u.HGetRelease()
	assert.Nil(t, err)
	assert.NotNil(t, results)
}

func TestHRepoAdd(t *testing.T) {
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	repobody := models.RepositoryBody{}
	repobody.Name = faker.Word()
	repobody.Url = "https://charts.helm.sh/stable"
	err := u.HelmRepoAdd(repobody)
	assert.Nil(t, err)
	err = u.HelmRepoAdd(repobody)
	assert.NotNil(t, err)
}

func TestCreateServiceAccount(t *testing.T) {
	servacc := models.ServiceAccount{
		Name:            Test.ServiceAccountName,
		Namespace:       "default",
		SecretNamespace: "default",
		SecretName:      Test.SecretName,
	}
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	service_account, err := u.CreateServiceAccount(servacc)
	assert.Nil(t, err)
	assert.NotNil(t, service_account)
}

func TestCreateRole(t *testing.T) {
	rolebody := models.Role{
		Name:      Test.RoleName,
		Namespace: "default",
		Verbs: []string{
			"get", "list", "watch",
		},
		Resources: []string{
			"pods",
		},
	}
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	role, err := u.CreateRole(rolebody)
	assert.Nil(t, err)
	assert.NotNil(t, role)
}

func TestCreateRoleBinding(t *testing.T) {
	rolebindingbody := models.RoleBinding{
		Name:        Test.RoleBindingName,
		Namespace:   "default",
		AccountName: Test.ServiceAccountName,
		RoleName:    Test.RoleName,
	}
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	rolebinding, err := u.CreateRoleBinding(rolebindingbody)
	assert.Nil(t, err)
	assert.NotNil(t, rolebinding)
}

func TestCreateCRD(t *testing.T) {
	u := services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{}))
	u.CreateCRD()
}

func TestMain(m *testing.M) {
	m.Run()
	DeleteAll()
}
