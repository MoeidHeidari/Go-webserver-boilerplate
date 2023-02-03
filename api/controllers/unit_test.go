package controllers

// import (
// 	"main/lib"
// 	"main/models"
// 	"main/services"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGetTest(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	gin.SetMode(gin.TestMode)
// 	_, r := gin.CreateTestContext(w)
// 	s := NewTestController(services.TestService{}, lib.Logger{})
// 	r.GET("/test", s.GetTest)
// 	request, err := http.NewRequest("GET", "/test", nil)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	resp := httptest.NewRecorder()

// 	r.ServeHTTP(resp, request)

// 	assert.Equal(t, 200, resp.Code)
// 	assert.JSONEq(t, models.Test, resp.Body.String())

// }
