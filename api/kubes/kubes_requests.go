package kubes

import (
	"context"
	"fmt"
	"log"
	"main/lib"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin" // swagger embed files
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeRequest struct {
	logger    lib.Logger
	clientset *kubernetes.Clientset
}

type ChartBody struct {
	ChartPath   string `json:"chart_path"`
	Namespace   string `json:"namespace"`
	ReleaseName string `json:"release_name"`
}

func NewKubeRequest(logger lib.Logger) KubeRequest {

	if err := helmInit(); err != nil {
		panic(err.Error())
	}

	clientset, err := clientsetInit()

	if err != nil {
		panic(err.Error())
	}

	return KubeRequest{
		logger:    logger,
		clientset: clientset,
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

func helmInit() error {

	return nil
}

// Gets release in Helm
func HGetRelease(release_name string) (release.Release, error) {
	panic("Not implemented exception")
}

func (u KubeRequest) HCreateReleaseRequest(c *gin.Context) {
	body := ChartBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	release, err := u.HCreateRelease(body.ChartPath, body.Namespace, body.ReleaseName)

	if err != nil {
		u.logger.Panic(err.Error())

		return
	}

	c.JSON(200, release.Name+"is created")
}

func (u KubeRequest) HCreateRelease(chartPath string, namespace string, releaseName string) (*release.Release, error) {

	settings := cli.New()
	actionConfiguration := new(action.Configuration)

	if err := actionConfiguration.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}

	client := action.NewInstall(actionConfiguration)
	client.Namespace = namespace
	client.ReleaseName = releaseName

	locatedChart, err := client.LocateChart(chartPath, settings)
	if err != nil {
		return nil, err
	}

	newChart, err := loader.Load(locatedChart)
	if err != nil {
		return nil, err
	}

	release, err := client.Run(newChart, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return release, nil
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

	pod_name := c.GetHeader("pod_name")
	container_name := c.GetHeader("container_name")
	image := c.GetHeader("image")
	namespace := c.GetHeader("namespace")
	command := c.GetHeader("command")
	command_value := c.GetHeader("command_value")

	newpod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod_name,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    container_name,
					Image:   image,
					Command: []string{command, command_value},
				},
			},
		},
	}

	_, err := u.clientset.CoreV1().Pods(namespace).Create(context.TODO(), newpod, metav1.CreateOptions{})

	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, pod_name+" was added.")
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
		panic(err.Error())
	}
	for _, n := range nodelist.Items {
		fmt.Println(n)
		pods, err := u.clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		_, err = u.clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod example-xxxxx not found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
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

func (u KubeRequest) DeletePodRequest(c *gin.Context) {
	pod_name := c.Param("pod_name")
	namespace := c.Param("namespace")

	fmt.Println(pod_name)

	err := u.clientset.CoreV1().Pods(namespace).Delete(context.TODO(), pod_name, metav1.DeleteOptions{})

	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, pod_name+" was deleted.")
}
