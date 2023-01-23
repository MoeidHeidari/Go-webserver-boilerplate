package controllers

import (
	"main/constants"
	"main/lib"
	"main/models"
	"main/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// @Summary Gets one test
// @Tags get tests
// @Description Get one test by id
// @Param id path int true "Test id"
// @Security ApiKeyAuth
// @Router /api/test/{id} [get]
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

// @Summary Get all test
// @Tags get tests
// @Description Get all the Tests
// @Security ApiKeyAuth
// @Router /api/test [get]
func (u TestController) GetTest(c *gin.Context) {
	Tests, err := u.service.GetAllTest()
	if err != nil {
		u.logger.Error(err)
	}
	c.JSON(200, gin.H{"data": Tests})
}

// @Summary Create GetTests
// @Tags create test
// @Description Create new test
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body string true "test data"
// @Router /api/test [post]
func (u TestController) CreateTest(c *gin.Context) {
	Test := models.Test{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&Test); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	Test.CreatedAt = time.Now()
	Test.UpdatedAt = time.Now()

	if err := u.service.WithTrx(trxHandle).CreateTest(Test); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "Test created"})
}

// @Summary Update test
// @Tags update test
// @Description Update an old test
// @Param id path int true "Test id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body string true "test data"
// @Router /api/test/{id} [post]
func (u TestController) UpdateTest(c *gin.Context) {

	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)

	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	Test, _ := u.service.GetOneTest(uint(id))
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&Test); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	Test.UpdatedAt = time.Now()

	if err := u.service.WithTrx(trxHandle).UpdateTest(Test); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"data": "Test updated"})
}

// @Summary delete test
// @Tags Delete
// @Description delete test
// @ID delete-test
// @Param id path int true "Test id"
// @Produce json
// @Security ApiKeyAuth
// @Router /api/test/{id} [delete]
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
