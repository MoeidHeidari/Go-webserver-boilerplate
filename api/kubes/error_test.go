package kubes_test

import (
	"context"
	"main/api/kubes"
	"main/lib"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/assert"
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
