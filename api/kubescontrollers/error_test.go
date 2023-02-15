package kubescontrollers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"main/api/kubescontrollers"
	"main/lib"
	"main/models"
	"main/repository"
	"main/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/action"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConfigMapUpdate(t *testing.T) {
	configmapbody := models.ConfigMapBody{}
	err := faker.FakeData(&configmapbody.Data)
	if err != nil {
		panic(err.Error())
	}
	configmapbody.Name = faker.Word()
	configmapbody.Namespace = "default"
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	cm, err := u.Service.CreateOrUpdateConfigMap(configmapbody)
	assert.Nil(t, err)
	assert.NotNil(t, cm)
	err = faker.FakeData(&configmapbody.Data)
	if err != nil {
		panic(err.Error())
	}
	cm, err = u.Service.CreateOrUpdateConfigMap(configmapbody)
	if err != nil {
		panic(err.Error())
	}
	assert.Nil(t, err)
	get_cm, err := u.Service.Repository.Clientset.CoreV1().ConfigMaps("default").Get(context.Background(), configmapbody.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, get_cm.Labels, cm.Labels)
	err = u.Service.Repository.Clientset.CoreV1().ConfigMaps("default").Delete(context.Background(), configmapbody.Name, metav1.DeleteOptions{})
	assert.Nil(t, err)

}

func TestSecretUpdate(t *testing.T) {
	secret := models.SecretBody{}
	err := faker.FakeData(&secret.Data)
	if err != nil {
		panic(err.Error())
	}
	secret.Name = faker.Word()
	secret.Namespace = "default"
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	cm, err := u.Service.CreateOrUpdateSecret(secret)
	assert.Nil(t, err)
	assert.NotNil(t, cm)
	err = faker.FakeData(&secret.Data)
	if err != nil {
		panic(err.Error())
	}
	cm, err = u.Service.CreateOrUpdateSecret(secret)
	if err != nil {
		panic(err.Error())
	}
	assert.Nil(t, err)
	get_cm, err := u.Service.Repository.Clientset.CoreV1().Secrets("default").Get(context.Background(), secret.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, get_cm.Labels, cm.Labels)
	err = u.Service.Repository.Clientset.CoreV1().Secrets("default").Delete(context.Background(), secret.Name, metav1.DeleteOptions{})
	assert.Nil(t, err)
}

func TestCreateNodePortError(t *testing.T) {
	nodeportbody := models.Nodeport{}
	nodeportbody.Name = faker.Word()
	nodeportbody.Namespace = "random"
	nodeportbody.RedirectPort = 10
	nodeportbody.Port = 20
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	service, err := u.Service.CreateNodePort(nodeportbody)
	assert.NotNil(t, err)
	assert.Nil(t, service)
}

func TestPodCreateError(t *testing.T) {
	podBody := models.PodBody{}
	podBody.Name = faker.FirstName()
	podBody.Namespace = faker.Word()
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	pod, err := u.Service.CreatePod(podBody)
	assert.Nil(t, pod)
	assert.NotNil(t, err)

}

func TestCreatePVError(t *testing.T) {
	pvbody := models.PV{}
	pvbody.Name = faker.Name()
	pvbody.Path = faker.Word()
	pvbody.Storage = "1Gi"
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	pv, err := u.Service.CreatePersistentVolume(pvbody)
	assert.Nil(t, pv)
	assert.NotNil(t, err)
}

func TestCreatePVCError(t *testing.T) {
	pvcbody := models.PVC{}
	pvcbody.Name = faker.Word()
	pvcbody.Namespace = faker.Word()
	pvcbody.Storage = "1Gi"
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	pvc, err := u.Service.CreatePersistentVolumeClaim(pvcbody)
	assert.Nil(t, pvc)
	assert.NotNil(t, err)
}

func TestHgetReleaseError(t *testing.T) {
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	u.Service.Repository.Settings.SetNamespace("random")
	results, _ := u.Service.HGetRelease()
	assert.Nil(t, results)
}

func TestHCreateReleaseRequestError(t *testing.T) {
	chart := models.ChartBody{}
	chart.Namespace = faker.Word()
	chart.ReleaseName = faker.Word()
	chart.ChartPath = faker.Word()
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.POST("/", u.HCreateReleaseRequest)
	jsonbytes, err := json.Marshal(chart)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHCreateReleaseError(t *testing.T) {
	chart := models.ChartBody{}
	chart.Namespace = faker.Word()
	chart.ReleaseName = faker.Word()
	chart.ChartPath = "https://charts.bitnami.com/bitnami/keycloak-13.0.2.tgz"
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	release, err := u.Service.HCreateRelease(chart)
	assert.NotNil(t, err)
	assert.Nil(t, release)
	client := action.NewUninstall(u.Service.Repository.ActionConfiguration)
	client.Run(chart.ReleaseName)
	chart.ChartPath = faker.Word()
	release, err = u.Service.HCreateRelease(chart)
	assert.NotNil(t, err)
	assert.Nil(t, release)
}

func TestHGetReleaseRequestError(t *testing.T) {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	gin.SetMode(gin.TestMode)
	k := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	k.Service.Repository.Settings.SetNamespace("random")
	router.GET("/", k.HGetReleaseRequest)
	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestGetPodInfoRequestError(t *testing.T) {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	gin.SetMode(gin.TestMode)
	k := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	router.GET("/:namespace", k.GetPodList)
	req := httptest.NewRequest("GET", "/nil", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestDeletePodRequestError(t *testing.T) {
	Podname := faker.Word()
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.DELETE("/:namespace/:pod_name", u.DeletePodRequest)
	ctx.Request, _ = http.NewRequest(http.MethodDelete, "/default/"+Podname, nil)
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, 404, w.Code)
}

func TestGetCurrentPodStatusError(t *testing.T) {
	podname := faker.Word()
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	err := u.Service.GetCurrentPodStatusRequest(podname)
	assert.Nil(t, err)
}

func TestCreatePodRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	r.POST("/", u.CreatePodRequest)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateNodePortRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	r.POST("/", u.CreateNodePortRequest)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateConfigmapRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	r.POST("/", u.CreateOrUpdateConfigMapRequest)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	configmap := models.ConfigMapBody{}
	err := faker.FakeData(&configmap.Data)
	if err != nil {
		panic(err.Error())
	}
	configmap.Name = faker.Word()
	configmap.Namespace = faker.Word()
	jsonbytes, err := json.Marshal(configmap)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateSecretRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	r.POST("/", u.CreateOrUpdateSecretRequest)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	secret := models.SecretBody{}
	secret.Name = faker.Word()
	secret.Namespace = faker.Word()
	jsonbytes, err := json.Marshal(secret)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateNamespaceRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})

	r.POST("/", u.CreateNamespaceRequest)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("default")))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// ctx.Request, _ = http.NewRequest(http.MethodPost, "/", nil)
	// r.ServeHTTP(w, ctx.Request)
	// assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreatePVRequestError(t *testing.T) {
	pvbody := models.PV{}
	pvbody.Name = faker.Name()
	pvbody.Storage = "1Gi"
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	r.POST("/", u.CreatePersistentVolumeRequest)
	jsonbytes, err := json.Marshal(pvbody)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreatePVCRequestError(t *testing.T) {
	pvc := models.PVC{}
	pvc.Storage = "1Gi"
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	r.POST("/", u.CreatePersistentVolumeClaimRequest)
	jsonbytes, err := json.Marshal(pvc)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHelmCreateRepoRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	r.POST("/", u.HCreateRepoRequest)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	repobody := models.RepositoryBody{}
	repobody.Name = faker.Word()
	repobody.Url = faker.URL()
	jsonbytes, err := json.Marshal(repobody)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateServiceAccountRequestError(t *testing.T) {
	servacc := models.ServiceAccount{
		Name:            ReqTest.ServiceAccountName,
		Namespace:       faker.Word(),
		SecretNamespace: "default",
		SecretName:      ReqTest.SecretName,
	}
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.POST("/", u.CreateServiceAccountRequest)
	jsonbytes, err := json.Marshal(servacc)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateRoleRequestError(t *testing.T) {
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	rolebody := models.Role{
		Name:      ReqTest.RoleName,
		Namespace: faker.Word(),
		Verbs: []string{
			"get", "list", "watch",
		},
		Resources: []string{
			"pods",
		},
	}
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.POST("/", u.CreateRoleRequest)
	jsonbytes, err := json.Marshal(rolebody)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateRoleBindingRequestError(t *testing.T) {
	u := kubescontrollers.NewKubeController(services.NewKubernetesService(lib.Logger{}, repository.NewKubernetesRepository(lib.NewKubernetesClient(lib.Logger{}), lib.Logger{})), lib.Logger{})
	rolebindingbody := models.RoleBinding{
		Name:        ReqTest.RoleBindingName,
		Namespace:   faker.Word(),
		AccountName: ReqTest.ServiceAccountName,
		RoleName:    ReqTest.RoleName,
	}
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.POST("/", u.CreateRoleBindingRequest)
	jsonbytes, err := json.Marshal(rolebindingbody)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
