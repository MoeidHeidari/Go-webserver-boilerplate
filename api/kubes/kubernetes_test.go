package kubes

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"main/lib"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetPodInfo(t *testing.T) {
	router := gin.Default()
	k := NewKubeRequest(lib.Logger{})
	router.GET("/test", k.GetPodInfoRequest)
	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
}

func TestConfigMapCreate(t *testing.T) {
	var configmapbody ConfigMapBody
	err := faker.FakeData(&configmapbody.Data)
	if err != nil {
		panic(err.Error())
	}
	configmapbody.Name = faker.Word()
	configmapbody.Namespace = "default"
	u := NewKubeRequest(lib.Logger{})
	cm, err := u.CreateOrUpdateConfigMap(configmapbody)
	assert.Nil(t, err)
	assert.NotNil(t, cm)
	get_cm, err := u.clientset.CoreV1().ConfigMaps("default").Get(context.Background(), configmapbody.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, get_cm.Labels, cm.Labels)
}

func TestNamespace(t *testing.T) {
	namespace_name := faker.Word()
	u := NewKubeRequest(lib.Logger{})
	namespace, err := u.CreateNamespace(namespace_name)
	assert.Nil(t, err)
	assert.Equal(t, namespace_name, namespace.Name)
}

func TestSercetCreate(t *testing.T) {
	var secretbody SecretBody
	err := faker.FakeData(&secretbody.Data)
	if err != nil {
		panic(err.Error())
	}
	secretbody.Name = faker.Word()
	secretbody.Namespace = "default"
	u := NewKubeRequest(lib.Logger{})
	s, err := u.CreateOrUpdateSecret(secretbody)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	get_secretbody, err := u.clientset.CoreV1().Secrets("default").Get(context.Background(), secretbody.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, get_secretbody.Labels, s.Labels)

}

func TestCreateNodePort(t *testing.T) {
	var nodeportbody Nodeport
	nodeportbody.Name = faker.Word()
	nodeportbody.Namespace = "default"
	nodeportbody.RedirectPort = int32(rand.Intn(32767-30000) + 30000)
	nodeportbody.Port = int32(rand.Intn(30000-20000) + 20000)
	u := NewKubeRequest(lib.Logger{})
	service, err := u.CreateNodePort(nodeportbody)
	assert.Nil(t, err)
	assert.NotNil(t, service)
	checkservice, err := u.clientset.CoreV1().Services("default").Get(context.Background(), nodeportbody.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, checkservice.UID, service.UID)
}

func TestCreatePod(t *testing.T) {
	var podBody PodBody
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
	u := NewKubeRequest(lib.Logger{})
	pod, err := u.CreatePod(podBody)
	assert.Nil(t, err)
	assert.NotNil(t, pod)
	checkpod, err := u.clientset.CoreV1().Pods("default").Get(context.TODO(), podBody.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, checkpod.UID, pod.UID)
}

func TestCreatePV(t *testing.T) {
	var pv PV
	pv.Name = faker.Word()
	pv.Path = faker.Word()
	pv.Storage = "1Gi"
	u := NewKubeRequest(lib.Logger{})
	persistent_volume, err := u.CreatePersistentVolume(pv)
	assert.Nil(t, err)
	assert.NotNil(t, persistent_volume)
	pvcheck, err := u.clientset.CoreV1().PersistentVolumes().Get(context.TODO(), pv.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, pvcheck.UID, persistent_volume.UID)
}

func TestCreatePVC(t *testing.T) {
	var pvc PVC
	pvc.Name = faker.Word()
	pvc.Namespace = "default"
	pvc.Storage = "1Gi"
	u := NewKubeRequest(lib.Logger{})
	persistent_volume_claim, err := u.CreatePersistentVolumeClaim(pvc)
	assert.Nil(t, err)
	assert.NotNil(t, persistent_volume_claim)
	pvccheck, err := u.clientset.CoreV1().PersistentVolumeClaims("default").Get(context.TODO(), pvc.Name, metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, pvccheck.UID, persistent_volume_claim.UID)
}

func TestCreatePodRequest(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJsonPost(ctx, map[string]interface{}{

		"command":        "[]",
		"configmap_name": "conff",
		"port":           27017,
		"secret_name":    "secrr",
		"name":           "mongos",
		"namespace":      "default",
		"container_name": "mongo-container",
		"claim_name":     "mongo-pvc",
		"image":          "mongo",
		"volume_name":    "mongodb-data",
		"mountpath":      "/usr/share/mongo",
	})
	u := NewKubeRequest(lib.Logger{})
	u.CreatePodRequest(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func MockJsonPost(c *gin.Context, content interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")
	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err.Error())
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}
