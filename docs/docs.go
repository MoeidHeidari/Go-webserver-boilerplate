// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/currency": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get one test by id",
                "tags": [
                    "get tests"
                ],
                "summary": "Gets currencies",
                "responses": {}
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Post request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "get tests"
                ],
                "summary": "Gets post responce",
                "parameters": [
                    {
                        "description": "Post form",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/test": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all the Tests",
                "consumes": [
                    "*/*"
                ],
                "tags": [
                    "get tests"
                ],
                "summary": "Get all test",
                "responses": {}
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create new test",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "create test"
                ],
                "summary": "Create GetTests",
                "parameters": [
                    {
                        "description": "test data",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/test/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get one test by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "get tests"
                ],
                "summary": "Gets one test",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Test id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update an old test",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "update test"
                ],
                "summary": "Update test",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Test id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "test data",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {}
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "delete test",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Delete"
                ],
                "summary": "delete test",
                "operationId": "delete-test",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Test id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "localhost:3000",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "SkyFarm",
	Description:      "The BEST API you have ever seen",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
