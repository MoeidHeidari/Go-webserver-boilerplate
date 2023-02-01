package kubes

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"main/lib"
	"net/http"
	"time"

	"github.com/gin-gonic/gin" // swagger embed files
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeRequest struct {
	logger              lib.Logger
	clientset           *kubernetes.Clientset
	settings            *cli.EnvSettings
	actionConfiguration *action.Configuration
}

func NewKubeRequest(logger lib.Logger) KubeRequest {

	settings := cli.New()
	actionConfiguration := new(action.Configuration)

	clientset, err := clientsetInit()

	if err != nil {
		panic(err.Error())
	}

	return KubeRequest{
		logger:              logger,
		clientset:           clientset,
		settings:            settings,
		actionConfiguration: actionConfiguration,
	}
}

func clientsetInit() (*kubernetes.Clientset, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return clientset, err
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

	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	pod, err := u.CreatePod(body)

	if err != nil {
		u.logger.Panic(err.Error())

		return
	}

	c.JSON(200, gin.H{
		"message": pod.Name + " is created",
	})
}

// @Summary Get pod info
// @Tags kubernetes
// @Accept json
// @Produce json
// @Description Post request
// @Security ApiKeyAuth
// @Router /api/kube_get [post]
func (u KubeRequest) GetPodInfoRequest(c *gin.Context) {

	nodelist, err := u.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		u.logger.Panic(err.Error())
	}
	for _, n := range nodelist.Items {
		pods, err := u.clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			u.logger.Panic(err.Error())
		}

		var names []string = make([]string, len(pods.Items))

		for i := 0; i < len(pods.Items); i++ {
			names[i] = pods.Items[i].Name
		}

		c.JSON(200, gin.H{
			"node name":  n.Name,
			"pods count": len(pods.Items),
			"pods names": names,
		})

		time.Sleep(10 * time.Second)
	}

}

func (u KubeRequest) GetCurrentPodStatusRequest(pod_name string) []byte {

	pod, err := u.clientset.CoreV1().Pods("default").Get(context.Background(), pod_name, metav1.GetOptions{})
	if err != nil {
		u.logger.Panic(err.Error())
	}
	status, err := json.Marshal(pod.Status)

	if err != nil {
		u.logger.Panic(err.Error())
	}
	return status
}

func (u KubeRequest) DeletePodRequest(c *gin.Context) {
	pod_name := c.Param("pod_name")
	namespace := c.Param("namespace")

	err := u.clientset.CoreV1().Pods(namespace).Delete(context.TODO(), pod_name, metav1.DeleteOptions{})

	if err != nil {
		c.JSON(404, "pod not found")
		return
	}

	c.JSON(200, gin.H{
		"message": pod_name + " is deleted",
	})
}

func (u KubeRequest) CreateOrUpdateConfigMapRequest(c *gin.Context) {
	body := ConfigMapBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}
	configmap, err := u.CreateOrUpdateConfigMap(body)

	if err != nil {
		u.logger.Panic(err.Error())
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"created/updated": configmap.Name,
	})
}

func (u KubeRequest) CreateOrUpdateSecretRequest(c *gin.Context) {
	body := SecretBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}
	secret, err := u.CreateOrUpdateSecret(body)

	if err != nil {
		u.logger.Panic(err.Error())
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"created/updated": secret.Name,
	})
}

func (u KubeRequest) CreateNamespaceRequest(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		u.logger.Panic(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	namespace, err := u.CreateNamespace(string(body))

	if err != nil {
		u.logger.Panic(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"created": namespace.Name,
	})
}

func (u KubeRequest) CreatePersistentVolumeRequest(c *gin.Context) {
	body := PV{}
	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}
	pv, err := u.CreatePersistentVolume(body)

	if err != nil {
		u.logger.Panic(err.Error())
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"created": pv.Name,
	})
}

func (u KubeRequest) CreatePersistentVolumeClaimRequest(c *gin.Context) {
	body := PVC{}
	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}
	pvc, err := u.CreatePersistentVolumeClaim(body)

	if err != nil {
		u.logger.Panic(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"created": pvc.Name,
	})
}

func (u KubeRequest) CreateVolumesPodRequest(c *gin.Context) {
	body := VolumesPod{}
	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}
	pod, err := u.CreateVolumesPod(body)

	if err != nil {
		u.logger.Panic(err.Error())
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"created": pod.Name,
	})
}
