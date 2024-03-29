// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/api/v1/getServiceInfo": {
            "get": {
                "description": "get service info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service Information"
                ],
                "summary": "get service info",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Insert your access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Get service info by GET method with token",
                        "schema": {
                            "$ref": "#/definitions/api.ServiceInfoSuccessResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ServiceInfoFailedResp"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.AuthFailedResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ServiceInfoFailedResp"
                        }
                    }
                }
            }
        },
        "/api/v1/login": {
            "post": {
                "description": "login and get token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json",
                    "text/plain"
                ],
                "tags": [
                    "AAA"
                ],
                "summary": "Login and return token after authenticate.",
                "parameters": [
                    {
                        "type": "string",
                        "default": "\u003cAdd api key here\u003e",
                        "description": "Insert your api key",
                        "name": "X-API-Key",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Account and Password",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.LoginReceiveBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.LoginSucceed"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.LoginFailed"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.AuthFailedResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.LoginFailed"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.AuthFailedResp": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "error"
                }
            }
        },
        "api.LoginFailed": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "format": "string",
                    "example": "msg"
                }
            }
        },
        "api.LoginReceiveBody": {
            "type": "object",
            "properties": {
                "Account": {
                    "type": "string",
                    "format": "string",
                    "example": "account"
                },
                "Password": {
                    "type": "string",
                    "format": "string",
                    "example": "password"
                }
            }
        },
        "api.LoginSucceed": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string",
                    "format": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50Ijoib3JnYWRtaW4iLCJyb2xlIjoiTWVtYmVyIiwiYXVkIjoib3JnYWRtaW4iLCJleHAiOjE2MjcwMzAzMjcsImp0aSI6Im9yZ2FkbWluMTYyNzAyOTEyNyIsImlhdCI6MTYyNzAyOTEyNywiaXNzIjoiSldUIiwibmJmIjoxNjI3MDI5MTI3LCJzdWIiOiJvcmdhZG1pbiJ9._HlIMv_2_vok5cLjxPNTI2qLUsIZQTtKbpW8UN-C4Tc"
                }
            }
        },
        "api.ServiceInfoFailedResp": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "msg"
                }
            }
        },
        "api.ServiceInfoSuccessResp": {
            "type": "object",
            "properties": {
                "serviceName": {
                    "type": "string",
                    "format": "string",
                    "example": "service"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
