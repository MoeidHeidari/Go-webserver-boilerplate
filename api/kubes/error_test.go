package kubes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"main/api/kubes"
	"main/lib"
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
	configmapbody := kubes.ConfigMapBody{}
	err := faker.FakeData(&configmapbody.Data)
	if err != nil {
		panic(err.Error())
	}
	configmapbody.Name = faker.Word()
	configmapbody.Namespace = "default"
	u := kubes.NewKubeRequest(lib.Logger{})
	cm, err := u.CreateOrUpdateConfigMap(configmapbody)
	assert.Nil(t, err)
	assert.NotNil(t, cm)
	err = faker.FakeData(&configmapbody.Data)
	if err != nil {
		panic(err.Error())
	}
	cm, err = u.CreateOrUpdateConfigMap(configmapbody)
	if err != nil {
		panic(err.Error())
	}
	assert.Nil(t, err)
	get_cm, err := u.Clientset.CoreV1().ConfigMaps("default").Get(context.Background(), configmapbody.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, get_cm.Labels, cm.Labels)
	err = u.Clientset.CoreV1().ConfigMaps("default").Delete(context.Background(), configmapbody.Name, metav1.DeleteOptions{})
	assert.Nil(t, err)

}

func TestSecretUpdate(t *testing.T) {
	secret := kubes.SecretBody{}
	err := faker.FakeData(&secret.Data)
	if err != nil {
		panic(err.Error())
	}
	secret.Name = faker.Word()
	secret.Namespace = "default"
	u := kubes.NewKubeRequest(lib.Logger{})
	cm, err := u.CreateOrUpdateSecret(secret)
	assert.Nil(t, err)
	assert.NotNil(t, cm)
	err = faker.FakeData(&secret.Data)
	if err != nil {
		panic(err.Error())
	}
	cm, err = u.CreateOrUpdateSecret(secret)
	if err != nil {
		panic(err.Error())
	}
	assert.Nil(t, err)
	get_cm, err := u.Clientset.CoreV1().Secrets("default").Get(context.Background(), secret.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, get_cm.Labels, cm.Labels)
	err = u.Clientset.CoreV1().Secrets("default").Delete(context.Background(), secret.Name, metav1.DeleteOptions{})
	assert.Nil(t, err)
}

func TestCreateNodePortError(t *testing.T) {
	nodeportbody := kubes.Nodeport{}
	nodeportbody.Name = faker.Word()
	nodeportbody.Namespace = "random"
	nodeportbody.RedirectPort = 10
	nodeportbody.Port = 20
	u := kubes.NewKubeRequest(lib.Logger{})
	service, err := u.CreateNodePort(nodeportbody)
	assert.NotNil(t, err)
	assert.Nil(t, service)
}

func TestPodCreateError(t *testing.T) {
	podBody := kubes.PodBody{}
	podBody.Name = faker.FirstName()
	podBody.Namespace = faker.Word()
	u := kubes.NewKubeRequest(lib.Logger{})
	pod, err := u.CreatePod(podBody)
	assert.Nil(t, pod)
	assert.NotNil(t, err)

}

func TestCreatePVError(t *testing.T) {
	pvbody := kubes.PV{}
	pvbody.Name = faker.Name()
	pvbody.Path = faker.Word()
	pvbody.Storage = "1Gi"
	u := kubes.NewKubeRequest(lib.Logger{})
	pv, err := u.CreatePersistentVolume(pvbody)
	assert.Nil(t, pv)
	assert.NotNil(t, err)
}

func TestCreatePVCError(t *testing.T) {
	pvcbody := kubes.PVC{}
	pvcbody.Name = faker.Word()
	pvcbody.Namespace = faker.Word()
	pvcbody.Storage = "1Gi"
	u := kubes.NewKubeRequest(lib.Logger{})
	pvc, err := u.CreatePersistentVolumeClaim(pvcbody)
	assert.Nil(t, pvc)
	assert.NotNil(t, err)
}

func TestHgetReleaseError(t *testing.T) {
	u := kubes.NewKubeRequest(lib.Logger{})
	u.Settings.SetNamespace("random")
	results, _ := u.HGetRelease()
	assert.Nil(t, results)
}

func TestHCreateReleaseRequestError(t *testing.T) {
	chart := kubes.ChartBody{}
	chart.Namespace = faker.Word()
	chart.ReleaseName = faker.Word()
	chart.ChartPath = faker.Word()
	u := kubes.NewKubeRequest(lib.Logger{})
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
	chart := kubes.ChartBody{}
	chart.Namespace = faker.Word()
	chart.ReleaseName = faker.Word()
	chart.ChartPath = "https://charts.bitnami.com/bitnami/keycloak-13.0.2.tgz"
	u := kubes.NewKubeRequest(lib.Logger{})
	release, err := u.HCreateRelease(chart)
	assert.NotNil(t, err)
	assert.Nil(t, release)
	client := action.NewUninstall(u.ActionConfiguration)
	client.Run(chart.ReleaseName)
	chart.ChartPath = faker.Word()
	release, err = u.HCreateRelease(chart)
	assert.NotNil(t, err)
	assert.Nil(t, release)
}

func TestHGetReleaseRequestError(t *testing.T) {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	gin.SetMode(gin.TestMode)
	k := kubes.NewKubeRequest(lib.Logger{})
	k.Settings.SetNamespace("random")
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
	k := kubes.NewKubeRequest(lib.Logger{})
	router.GET("/:namespace", k.GetNodeInfoRequest)
	req := httptest.NewRequest("GET", "/nil", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestDeletePodRequestError(t *testing.T) {
	Podname := faker.Word()
	u := kubes.NewKubeRequest(lib.Logger{})
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
	u := kubes.NewKubeRequest(lib.Logger{})
	err := u.GetCurrentPodStatusRequest(podname)
	assert.Nil(t, err)
}

func TestCreatePodRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})
	r.POST("/", u.CreatePodRequest)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateNodePortRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})
	r.POST("/", u.CreateNodePortRequest)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateConfigmapRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})
	r.POST("/", u.CreateOrUpdateConfigMapRequest)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	configmap := kubes.ConfigMapBody{}
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
	u := kubes.NewKubeRequest(lib.Logger{})
	r.POST("/", u.CreateOrUpdateSecretRequest)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	secret := kubes.SecretBody{}
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
	u := kubes.NewKubeRequest(lib.Logger{})

	r.POST("/", u.CreateNamespaceRequest)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("default")))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// ctx.Request, _ = http.NewRequest(http.MethodPost, "/", nil)
	// r.ServeHTTP(w, ctx.Request)
	// assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreatePVRequestError(t *testing.T) {
	pvbody := kubes.PV{}
	pvbody.Name = faker.Name()
	pvbody.Storage = "1Gi"
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})
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
	pvc := kubes.PVC{}
	pvc.Storage = "1Gi"
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})
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
	u := kubes.NewKubeRequest(lib.Logger{})
	r.POST("/", u.HCreateRepoRequest)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	repobody := kubes.RepositoryBody{}
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
