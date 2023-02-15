package controllers

import (
	"errors"
	"io/ioutil"
	"main/lib"
	"main/models"
	"main/services"
	"net/http"
	"net/url"
	"os"
	"strings"

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
func (u TestController) GetOneWorkspace(c *gin.Context) {
	paramID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	Test, err := u.service.GetOneWorkspace(objID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(200, Test)

}

// @Summary Get all test
// @Tags get tests
// @Description Get all the Tests
// @Accept */*
// @Security ApiKeyAuth
// @Router /api/test [get]
func (u TestController) GetWorkspaces(c *gin.Context) {
	Tests, err := u.service.GetAllWorkspaces()
	if err != nil || Tests == nil {
		c.JSON(http.StatusNotFound, err)
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
func (u TestController) CreateWorkspace(c *gin.Context) {
	Workspace := models.Workspace{}
	//trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&Workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	id, err := u.service.CreateWorkspace(Workspace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"workspace_id": id,
	})
}

func (u TestController) ValidateCardLabel(n models.Node) error {
	allowedNames := []string{
		"PV",
		"VM",
		"pod",
		"ingress",
		"service",
		"storage",
		"claim",
		"endpoints",
		"PVC",
		"rules",
	}
	for _, label := range allowedNames {
		if strings.EqualFold(n.CardLabel, label) {
			return nil
		}
	}
	return errors.New("invalid CardLabel")
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
func (u TestController) AddNode(c *gin.Context) {
	paramID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	Node := models.Node{}
	if err := c.ShouldBindJSON(&Node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := u.ValidateCardLabel(Node); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := u.service.AddNode(Node, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, "Workspace updated")
}

func (u TestController) AddEdge(c *gin.Context) {
	paramID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	Edge := models.Edge{}
	if err := c.ShouldBindJSON(&Edge); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := u.service.AddEdge(Edge, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, "Workspace updated")

}

// @Summary delete test
// @Tags Delete
// @Description delete test
// @ID delete-test
// @Param id path int true "Test id"
// @Produce json
// @Security ApiKeyAuth
// @Router /api/test/{id} [delete]
func (u TestController) DeleteWorkspace(c *gin.Context) {
	paramID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	if err := u.service.DeleteWorkspace(objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, "Test deleted")
}

func (u TestController) DeleteNode(c *gin.Context) {
	workspace_id := c.Param("id")
	node_id := c.Param("node_id")
	objID, err := primitive.ObjectIDFromHex(workspace_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	if err := u.service.DeleteNode(objID, node_id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, "Node deleted")
}

func (u TestController) UpdateNode(c *gin.Context) {
	Node := models.Node{}
	workspace_id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(workspace_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	if err := c.ShouldBindJSON(&Node); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := u.ValidateCardLabel(Node); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := u.service.UpdateNode(objID, Node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, "node updated")

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
		c.String(500, err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(500, err.Error())
	}
	resp.Body.Close()
	c.JSON(200, strings.Split(strings.Split(string(body), ",")[0], ":")[1])
	u.logger.Info(resp)
}
