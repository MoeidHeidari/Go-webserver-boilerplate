package kubes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u KubeRequest) HCreateReleaseRequest(c *gin.Context) {
	body := ChartBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	release, err := u.HCreateRelease(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, release)
}

func (u KubeRequest) HGetReleaseRequest(c *gin.Context) {
	results, err := u.HGetRelease()
	if err != nil || results == nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(200, results)
}

func (u KubeRequest) HCreateRepoRequest(c *gin.Context) {
	body := RepositoryBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	err := u.HelmRepoAdd(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(200, "repo created")
}
