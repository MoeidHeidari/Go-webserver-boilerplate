package kubescontrollers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u KubeController) HCreateReleaseRequest(c *gin.Context) {
	body := models.ChartBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	release, err := u.Service.HCreateRelease(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, release)
}

func (u KubeController) HGetReleaseRequest(c *gin.Context) {
	results, err := u.Service.HGetRelease()
	if err != nil || results == nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(200, results)
}

func (u KubeController) HCreateRepoRequest(c *gin.Context) {
	body := models.RepositoryBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	err := u.Service.HelmRepoAdd(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(200, "repo created")
}
