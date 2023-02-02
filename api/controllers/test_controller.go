package controllers

import (
	"io/ioutil"
	"main/lib"
	"main/models"
	"main/services"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// @Param id path string true "Test id"
// @Produce json
// @Security ApiKeyAuth
// @Router /api/test/{id} [get]
func (u TestController) GetOneTest(c *gin.Context) {
	paramID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return
	}

	Test, err := u.service.GetOneTest(objID)

	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, Test)

}

// @Summary Get all test
// @Tags get tests
// @Description Get all the Tests
// @Accept */*
// @Security ApiKeyAuth
// @Router /api/test [get]
func (u TestController) GetTest(c *gin.Context) {
	Tests, err := u.service.GetAllTest()
	if err != nil {
		u.logger.Error(err)
	}
	c.JSON(200, Tests)
}

// @Summary Get all test fields
// @Tags get all test fields
// @Description Get all test fields
// @Param field_name path string true "Field"
// @Produce json
// @Security ApiKeyAuth
// @Router /api/test [get]
func (u TestController) GetTestField(c *gin.Context) {
	field_name := c.Param("field_name")
	Tests, err := u.service.GetAllTestField(field_name)
	if err != nil {
		u.logger.Error(err)
	}

	c.JSON(200, Tests)
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
	//trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&Test); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	Test.CreatedAt = time.Now()
	Test.UpdatedAt = time.Now()
	Test.ID = primitive.NewObjectID()

	if err := u.service.CreateTest(Test); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, "Test created")
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

	objID, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return
	}

	Test, _ := u.service.GetOneTest(objID)

	if err := c.ShouldBindJSON(&Test); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	Test.UpdatedAt = time.Now()

	if err := u.service.UpdateTest(Test); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, "Test updated")
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

	if err := u.service.DeleteTest(paramID); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, "Test deleted")
}

func (u TestController) GetCode(c *gin.Context) {
	paramcode := c.Query("code")
	Url := "http://localhost:8080/realms/master/protocol/openid-connect/token"
	url_form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {paramcode},
		"client_id":     {"skyfarm"},
		"redirect_uri":  {"http://localhost:3000/get_code"},
		"client_secret": {os.Getenv("JWT_SECRET")},
	}
	resp, err := http.PostForm(Url, url_form)
	if err != nil {
		u.logger.Error(err)
		c.String(500, err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		u.logger.Error(err)
		c.String(500, err.Error())
	}
	resp.Body.Close()
	c.JSON(200, strings.Split(strings.Split(string(body), ",")[0], ":")[1])
	u.logger.Info(resp)
}
