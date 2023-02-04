package kubes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u KubeRequest) HCreateReleaseRequest(c *gin.Context) {
	body := ChartBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	release, err := u.HCreateRelease(body)

	if err != nil {
		u.logger.Panic(err.Error())

		return
	}

	c.JSON(200, gin.H{
		"message": release.Name + " is created",
	})
}

func (u KubeRequest) HGetReleaseRequest(c *gin.Context) {
	results, err := u.HGetRelease()
	if err != nil {
		u.logger.Panic(err.Error())
	}
	c.JSON(200, results)
}
