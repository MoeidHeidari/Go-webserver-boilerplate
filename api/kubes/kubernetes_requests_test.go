package kubes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"main/api/kubes"
	"main/lib"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/action"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ReqTest = test{
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

func DeleteAllReq() {
	u := kubes.NewKubeRequest(lib.Logger{})
	err := u.Clientset.CoreV1().Pods("default").Delete(context.TODO(), ReqTest.PodName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Clientset.CoreV1().PersistentVolumeClaims("default").Delete(context.TODO(), ReqTest.PVCName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Clientset.CoreV1().PersistentVolumes().Delete(context.TODO(), ReqTest.PVName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Clientset.CoreV1().Services("default").Delete(context.Background(), ReqTest.NodeportName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Clientset.CoreV1().Secrets("default").Delete(context.Background(), ReqTest.SecretName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Clientset.CoreV1().Namespaces().Delete(context.Background(), ReqTest.Namespace, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Clientset.CoreV1().ConfigMaps("default").Delete(context.Background(), ReqTest.ConfigmapName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Clientset.RbacV1().Roles("default").Delete(context.Background(), ReqTest.RoleName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Clientset.RbacV1().RoleBindings("default").Delete(context.Background(), ReqTest.RoleBindingName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	err = u.Clientset.CoreV1().ServiceAccounts("default").Delete(context.Background(), ReqTest.ServiceAccountName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	pvlist, err := u.Clientset.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, i := range pvlist.Items {
		if i.Spec.ClaimRef.Name == ("data-" + ReqTest.ChartName + "-postgresql-0") {
			err = u.Clientset.CoreV1().PersistentVolumeClaims("default").Delete(context.TODO(), i.Spec.ClaimRef.Name, metav1.DeleteOptions{})
			if err != nil {
				panic(err.Error())
			}
			err = u.Clientset.CoreV1().PersistentVolumes().Delete(context.TODO(), i.Name, metav1.DeleteOptions{})
			if err != nil {
				panic(err.Error())
			}

		}
	}
	u.ActionConfiguration.Init(u.Settings.RESTClientGetter(), u.Settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf)
	client := action.NewUninstall(u.ActionConfiguration)
	client.DisableHooks = true
	client.Run(ReqTest.ChartName)
}
func TestGetPodInfo(t *testing.T) {
	router := gin.Default()
	gin.SetMode(gin.TestMode)
	k := kubes.NewKubeRequest(lib.Logger{})
	router.GET("/:namespace", k.GetNodeInfoRequest)
	req, _ := http.NewRequest("GET", "/default", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
}

func TestDeletePod(t *testing.T) {
	var podBody kubes.PodBody
	podBody.Name = faker.Word()
	podBody.Namespace = "default"
	podBody.ClaimName = faker.Word()
	podBody.VolumeName = faker.Word()
	podBody.SecretName = faker.Word()
	podBody.ContainerName = faker.Word()
	podBody.Image = faker.Word()
	podBody.Port = int32(rand.Intn(30000-20000) + 20000)
	podBody.MountPath = faker.URL()
	podBody.ConfigmapName = faker.Word()
	u := kubes.NewKubeRequest(lib.Logger{})
	pod, err := u.CreatePod(podBody)
	assert.Nil(t, err)
	assert.NotNil(t, pod)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.DELETE("/:namespace/:pod_name", u.DeletePodRequest)
	ctx.Request, _ = http.NewRequest(http.MethodDelete, "/default/"+podBody.Name, nil)
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePodRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})
	r.POST("/", u.CreatePodRequest)
	podBody := kubes.PodBody{}
	faker.FakeData(&podBody)
	podBody.Name = ReqTest.PodName
	podBody.Namespace = "default"
	podBody.ClaimName = faker.Word()
	podBody.VolumeName = faker.Word()
	podBody.SecretName = faker.Word()
	podBody.ContainerName = faker.Word()
	podBody.Image = faker.Word()
	podBody.Port = int32(rand.Intn(30000-20000) + 20000)
	podBody.MountPath = faker.URL()
	podBody.ConfigmapName = faker.Word()
	jsonbytes, err := json.Marshal(podBody)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateConfigmapRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	configmap := kubes.ConfigMapBody{}
	configmap.Name = ReqTest.ConfigmapName
	configmap.Namespace = "default"
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})
	r.POST("/", u.CreateOrUpdateConfigMapRequest)
	jsonbytes, err := json.Marshal(configmap)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateSecretRequestTest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := kubes.SecretBody{}
	secret.Name = ReqTest.SecretName
	secret.Namespace = "default"
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})
	r.POST("/", u.CreateOrUpdateSecretRequest)
	jsonbytes, err := json.Marshal(secret)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateNamespaceRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})

	r.POST("/", u.CreateNamespaceRequest)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(ReqTest.Namespace)))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePVRequest(t *testing.T) {
	pv := kubes.PV{}
	pv.Name = ReqTest.PVName
	pv.Path = faker.Word()
	pv.Storage = "1Gi"
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})
	r.POST("/", u.CreatePersistentVolumeRequest)
	jsonbytes, err := json.Marshal(pv)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePVCRequest(t *testing.T) {
	pvc := kubes.PVC{}
	pvc.Name = ReqTest.PVCName
	pvc.Namespace = "default"
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
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateNodePortRequest(t *testing.T) {
	nodeport := kubes.Nodeport{}
	nodeport.Name = ReqTest.NodeportName
	nodeport.Namespace = "default"
	nodeport.RedirectPort = int32(rand.Intn(32767-30000) + 30000)
	nodeport.Port = int32(rand.Intn(30000-20000) + 20000)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	u := kubes.NewKubeRequest(lib.Logger{})
	r.POST("/", u.CreateNodePortRequest)
	jsonbytes, err := json.Marshal(nodeport)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHCreateReleaseRequest(t *testing.T) {
	chart := kubes.ChartBody{}
	chart.Namespace = "default"
	chart.ReleaseName = ReqTest.ChartName
	chart.ChartPath = "https://charts.bitnami.com/bitnami/keycloak-13.0.2.tgz"
	u := kubes.NewKubeRequest(lib.Logger{})
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.POST("/", u.HCreateReleaseRequest)
	jsonbytes, err := json.Marshal(chart)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHGetReleaseRequest(t *testing.T) {
	router := gin.Default()
	gin.SetMode(gin.TestMode)
	k := kubes.NewKubeRequest(lib.Logger{})
	router.GET("/", k.HGetReleaseRequest)
	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
}

func TestHCreateRepoRequest(t *testing.T) {
	u := kubes.NewKubeRequest(lib.Logger{})
	repobody := kubes.RepositoryBody{}
	repobody.Name = faker.Word()
	repobody.Url = "https://charts.helm.sh/stable"
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.POST("/", u.HCreateRepoRequest)
	jsonbytes, err := json.Marshal(repobody)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateServiceAccountRequest(t *testing.T) {
	servacc := kubes.ServiceAccount{
		Name:            ReqTest.ServiceAccountName,
		Namespace:       "default",
		SecretNamespace: "default",
		SecretName:      ReqTest.SecretName,
	}
	u := kubes.NewKubeRequest(lib.Logger{})
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
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRoleRequest(t *testing.T) {
	u := kubes.NewKubeRequest(lib.Logger{})
	rolebody := kubes.Role{
		Name:      ReqTest.RoleName,
		Namespace: "default",
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
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRoleBindingRequest(t *testing.T) {
	u := kubes.NewKubeRequest(lib.Logger{})
	rolebindingbody := kubes.RoleBinding{
		Name:        ReqTest.RoleBindingName,
		Namespace:   "default",
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
	assert.Equal(t, http.StatusOK, w.Code)
}
