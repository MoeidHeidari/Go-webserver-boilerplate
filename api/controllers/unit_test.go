package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/api/controllers"
	"main/lib"
	"main/models"
	"main/repository"
	"main/services"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type test struct {
	Collection     mongo.Collection
	Workspace_id   string
	Testcontroller controllers.TestController
}

var Test test

func TestWorkspaceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.GET("/", Test.Testcontroller.GetWorkspaces)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateWorkspaceError(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.POST("/", Test.Testcontroller.CreateWorkspace)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(faker.Word())))

	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateWorkspace(t *testing.T) {
	workspace := models.Workspace{
		Name: faker.Word(),
	}
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.POST("/", Test.Testcontroller.CreateWorkspace)
	jsonbytes, err := json.Marshal(workspace)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, err = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonbytes))
	if err != nil {
		panic(err.Error())
	}
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
	s, _ := io.ReadAll(w.Body)
	Test.Workspace_id = (strings.Split(strings.Split(string(s), ":")[1], "}")[0])
	Test.Workspace_id = strings.Trim(Test.Workspace_id, "\"")
	if err != nil {
		panic(err.Error())
	}
}

func TestAddNode(t *testing.T) {
	node := models.Node{
		CpuNumber:     2,
		MemoryNumber:  8,
		StorageNumber: 20,
		Position: models.Coordinates{
			X: 20,
			Y: 50,
		},
		NodeName:   faker.Word(),
		CardLabel:  "PV",
		LabelColor: faker.Word(),
	}
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.POST("/:id", Test.Testcontroller.AddNode)
	jsonbytes, err := json.Marshal(node)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, err = http.NewRequest(http.MethodPost, "/"+Test.Workspace_id, bytes.NewBuffer(jsonbytes))
	if err != nil {
		panic(err.Error())
	}
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	ctx.Request, err = http.NewRequest(http.MethodPost, "/"+Test.Workspace_id, bytes.NewBuffer([]byte(faker.Word())))
	if err != nil {
		panic(err.Error())
	}
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	w = httptest.NewRecorder()
	ctx.Request, err = http.NewRequest(http.MethodPost, "/"+faker.Word(), bytes.NewBuffer([]byte(faker.Word())))
	if err != nil {
		panic(err.Error())
	}
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	w = httptest.NewRecorder()
	node.CardLabel = faker.Word()
	jsonbytes, err = json.Marshal(node)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request, err = http.NewRequest(http.MethodPost, "/"+Test.Workspace_id, bytes.NewBuffer(jsonbytes))
	if err != nil {
		panic(err.Error())
	}
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateNode(t *testing.T) {
	node := models.Node{
		ID:            "1",
		CpuNumber:     4,
		MemoryNumber:  8,
		StorageNumber: 12,
		Position: models.Coordinates{
			X: 200,
			Y: 500,
		},
		NodeName:   faker.Word(),
		CardLabel:  "PVC",
		LabelColor: faker.Word(),
	}
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.POST("/:id", Test.Testcontroller.UpdateNode)
	jsonbytes, err := json.Marshal(node)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/"+Test.Workspace_id, bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	ctx.Request = httptest.NewRequest(http.MethodPost, "/"+Test.Workspace_id, bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	w = httptest.NewRecorder()
	ctx.Request = httptest.NewRequest(http.MethodPost, "/"+faker.Word(), bytes.NewBuffer([]byte(faker.Word())))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	node.CardLabel = faker.Word()
	jsonbytes, err = json.Marshal(node)
	if err != nil {
		panic(err.Error())
	}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/"+faker.Word(), bytes.NewBuffer(jsonbytes))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetOneWorkspace(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.GET("/:id", Test.Testcontroller.GetOneWorkspace)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/"+faker.Word(), nil)
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	w = httptest.NewRecorder()
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/"+Test.Workspace_id, nil)
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestGetWorkspaces(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.GET("/", Test.Testcontroller.GetWorkspaces)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteNode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.DELETE("/:id/:node_id", Test.Testcontroller.DeleteNode)
	TestAddNode(t)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/"+Test.Workspace_id+"/2", nil)
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/"+faker.Word()+"/2", nil)
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/"+faker.Word()+"/a", nil)
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	fmt.Println(w.Body)
}

func TestDeleteWorkspace(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.DELETE("/:id", Test.Testcontroller.DeleteWorkspace)
	ctx.Request, _ = http.NewRequest(http.MethodDelete, "/"+Test.Workspace_id, nil)
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMain(m *testing.M) {
	dbUrl := "mongodb+srv://user:admin@pozhiloy.qaqey6i.mongodb.net/Kubernetes?authMechanism=SCRAM-SHA-1"
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(dbUrl).
		SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err.Error())
	}
	collection := client.Database("test").Collection("test_collection")

	Test.Collection = *collection
	u := controllers.NewTestController(services.NewTestService(lib.Logger{}, repository.NewTestRepository(lib.Database{Collection: &Test.Collection}, lib.Logger{})), lib.Logger{})
	Test.Testcontroller = u
	m.Run()
}
