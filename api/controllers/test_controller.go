package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/constants"
	"main/lib"
	"main/models"
	"main/services"
	"net/http"
	"strconv"
)

// TestController data type
type TestController struct {
	service services.TestService
	logger  lib.Logger
}

// NewTestController creates new Test controller
func NewTestController(TestService services.TestService, logger lib.Logger) TestController {
	return TestController{
		service: TestService,
		logger:  logger,
	}
}

// GetOneTest gets one Test
func (u TestController) GetOneTest(c *gin.Context) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	Test, err := u.service.GetOneTest(uint(id))

	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": Test,
	})

}

// GetTest gets the Test
func (u TestController) GetTest(c *gin.Context) {
	Tests, err := u.service.GetAllTest()
	if err != nil {
		u.logger.Error(err)
	}
	c.JSON(200, gin.H{"data": Tests})
}

// SaveTest saves the Test
func (u TestController) SaveTest(c *gin.Context) {
	Test := models.Test{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&Test); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := u.service.WithTrx(trxHandle).CreateTest(Test); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "Test created"})
}

// UpdateTest updates Test
func (u TestController) UpdateTest(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Test updated"})
}

// DeleteTest deletes Test
func (u TestController) DeleteTest(c *gin.Context) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := u.service.DeleteTest(uint(id)); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "Test deleted"})
}
