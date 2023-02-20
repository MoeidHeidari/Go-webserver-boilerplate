package kubescontrollers

import (
	"io/ioutil"
	"main/lib"
	"main/models"
	"main/services"
	"net/http"

	"github.com/gin-gonic/gin" // swagger embed files
)

type KubeController struct {
	Service services.KubernetesService
	logger  lib.Logger
}

func NewKubeController(kube_Service services.KubernetesService, logger lib.Logger) KubeController {
	return KubeController{
		Service: kube_Service,
		logger:  logger,
	}

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
func (u KubeController) CreatePodRequest(c *gin.Context) {
	body := models.PodBody{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	pod, err := u.Service.CreatePod(body)
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
func (u KubeController) GetPodList(c *gin.Context) {
	namespace := c.Param("namespace")
	names, err := u.Service.GetPodsList(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(200, gin.H{
		"pods": names,
	})

}

func (u KubeController) DeletePodRequest(c *gin.Context) {
	pod_name := c.Param("pod_name")
	namespace := c.Param("namespace")

	err := u.Service.DeletePod(pod_name, namespace)

	if err != nil {
		c.JSON(404, "pod not found")
	}

	c.JSON(200, gin.H{
		"message": pod_name + " is deleted",
	})
}

func (u KubeController) CreateOrUpdateConfigMapRequest(c *gin.Context) {
	body := models.ConfigMapBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	configmap, err := u.Service.CreateOrUpdateConfigMap(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, gin.H{
		"created/updated": configmap.Name,
	})
}

func (u KubeController) CreateOrUpdateSecretRequest(c *gin.Context) {
	body := models.SecretBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	secret, err := u.Service.CreateOrUpdateSecret(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, gin.H{
		"created/updated": secret.Name,
	})
}

func (u KubeController) CreateNamespaceRequest(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	namespace, err := u.Service.CreateNamespace(string(body))

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, namespace)
}

func (u KubeController) CreatePersistentVolumeRequest(c *gin.Context) {
	body := models.PV{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	pv, err := u.Service.CreatePersistentVolume(body)

	if err != nil || pv == nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, pv)
}

func (u KubeController) CreatePersistentVolumeClaimRequest(c *gin.Context) {
	body := models.PVC{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	pvc, err := u.Service.CreatePersistentVolumeClaim(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, pvc)
}

func (u KubeController) CreateNodePortRequest(c *gin.Context) {
	body := models.Nodeport{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	Service, err := u.Service.CreateNodePort(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, Service)
}

func (u KubeController) CreateRoleRequest(c *gin.Context) {
	body := models.Role{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	role, err := u.Service.CreateRole(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, role)
}

func (u KubeController) CreateRoleBindingRequest(c *gin.Context) {
	body := models.RoleBinding{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	rolebinding, err := u.Service.CreateRoleBinding(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(200, rolebinding)
}

func (u KubeController) CreateServiceAccountRequest(c *gin.Context) {
	body := models.ServiceAccount{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	rolebinding, err := u.Service.CreateServiceAccount(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(200, rolebinding)
}
