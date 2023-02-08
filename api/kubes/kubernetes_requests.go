package kubes

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"main/lib"
	"net/http"

	"github.com/gin-gonic/gin" // swagger embed files
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeRequest struct {
	logger              lib.Logger
	Clientset           *kubernetes.Clientset
	Settings            *cli.EnvSettings
	ActionConfiguration *action.Configuration
}

func NewKubeRequest(logger lib.Logger) KubeRequest {
	Settings := cli.New()
	ActionConfiguration := new(action.Configuration)

	Clientset, err := ClientsetInit()

	if err != nil {
		panic(err.Error())
	}

	return KubeRequest{
		logger:              logger,
		Clientset:           Clientset,
		Settings:            Settings,
		ActionConfiguration: ActionConfiguration,
	}
}

func ClientsetInit() (*kubernetes.Clientset, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()

	if err != nil {
		return nil, err
	}

	Clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return Clientset, err
}

// @Summary Create a pod
// @Tags kubernetes
// @Accept json
// @Produce json
// @Param namespace path string true "Field"
// @Param pod_name path string true "Field"
// @Param container_name path string true "Field"
// @Param image_type path string true "Field"
// @Param command query []string true "Field"
// @Description Post request
// @Security ApiKeyAuth
// @Router /api/kube_add [post]
func (u KubeRequest) CreatePodRequest(c *gin.Context) {
	body := PodBody{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	pod, err := u.CreatePod(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, pod)
}

// @Summary Get pod info/*  */
// @Tags kubernetes
// @Accept json
// @Produce json
// @Description Post request
// @Security ApiKeyAuth
// @Router /api/kube_get [post]
func (u KubeRequest) GetNodeInfoRequest(c *gin.Context) {
	namespace := c.Param("namespace")
	nodelist, err := u.Clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, n := range nodelist.Items {
		pods, err := u.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil || len(pods.Items) == 0 {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			var names []string = make([]string, len(pods.Items))

			for i := 0; i < len(pods.Items); i++ {
				names[i] = pods.Items[i].Name
			}

			c.JSON(200, gin.H{
				"node name":  n.Name,
				"pods names": names,
			})
		}
	}

}

func (u KubeRequest) GetCurrentPodStatusRequest(pod_name string) []byte {

	pod, err := u.Clientset.CoreV1().Pods("default").Get(context.Background(), pod_name, metav1.GetOptions{})
	if err != nil {
		return nil
	}
	status, err := json.Marshal(pod.Status)

	if err != nil {
		return nil
	}
	return status
}

func (u KubeRequest) DeletePodRequest(c *gin.Context) {
	pod_name := c.Param("pod_name")
	namespace := c.Param("namespace")

	err := u.Clientset.CoreV1().Pods(namespace).Delete(context.TODO(), pod_name, metav1.DeleteOptions{})

	if err != nil {
		c.JSON(404, "pod not found")
	}

	c.JSON(200, gin.H{
		"message": pod_name + " is deleted",
	})
}

func (u KubeRequest) CreateOrUpdateConfigMapRequest(c *gin.Context) {
	body := ConfigMapBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	configmap, err := u.CreateOrUpdateConfigMap(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, gin.H{
		"created/updated": configmap.Name,
	})
}

func (u KubeRequest) CreateOrUpdateSecretRequest(c *gin.Context) {
	body := SecretBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	secret, err := u.CreateOrUpdateSecret(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, gin.H{
		"created/updated": secret.Name,
	})
}

func (u KubeRequest) CreateNamespaceRequest(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	namespace, err := u.CreateNamespace(string(body))

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, namespace)
}

func (u KubeRequest) CreatePersistentVolumeRequest(c *gin.Context) {
	body := PV{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	pv, err := u.CreatePersistentVolume(body)

	if err != nil || pv == nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, pv)
}

func (u KubeRequest) CreatePersistentVolumeClaimRequest(c *gin.Context) {
	body := PVC{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	pvc, err := u.CreatePersistentVolumeClaim(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, pvc)
}

func (u KubeRequest) CreateNodePortRequest(c *gin.Context) {
	body := Nodeport{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	service, err := u.CreateNodePort(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, service)
}
