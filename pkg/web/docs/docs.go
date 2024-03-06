// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/auth/signIn": {
            "post": {
                "description": "Return session.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Sign-in user.",
                "parameters": [
                    {
                        "description": "query params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web_admin.SignInAdmin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web_admin.Session"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            }
        },
        "/admin/user/{user_id}": {
            "get": {
                "description": "Return users.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Get users.",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Page number",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 25,
                        "description": "Object count in page",
                        "name": "per_page",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.User"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            },
            "put": {
                "description": "Return user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Update user.",
                "parameters": [
                    {
                        "description": "query params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web_admin.UpdateUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Return nothing.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Delete user.",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            }
        },
        "/auth/signIn": {
            "post": {
                "description": "Return session.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign-in user.",
                "parameters": [
                    {
                        "description": "query params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web_client.SignInUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web_client.Session"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            }
        },
        "/auth/signOut": {
            "get": {
                "description": "Return Status OK.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Remove user session",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            }
        },
        "/auth/signUp": {
            "post": {
                "description": "Return session.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign-up new user.",
                "parameters": [
                    {
                        "description": "query params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web_client.SignUpUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web_client.Session"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    },
                    "402": {
                        "description": "Payment Required",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            }
        },
        "/user/": {
            "get": {
                "description": "Return user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            },
            "put": {
                "description": "Return user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Update user.",
                "parameters": [
                    {
                        "description": "query params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web_client.UpdateUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            }
        },
        "/user/activate": {
            "get": {
                "description": "Return token.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Activate user.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.UserActivate"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            },
            "post": {
                "description": "Return token.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "ActivateCode user.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/testing_dating_api_pkg_web_common.WebError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Admin": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "domain.User": {
            "type": "object",
            "properties": {
                "activated": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        },
        "domain.UserActivate": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "testing_dating_api_pkg_web_common.WebError": {
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                }
            }
        },
        "web_admin.Session": {
            "type": "object",
            "properties": {
                "admin": {
                    "$ref": "#/definitions/domain.Admin"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "web_admin.SignInAdmin": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "maxLength": 120,
                    "minLength": 3,
                    "example": "admin"
                },
                "password": {
                    "type": "string",
                    "maxLength": 40,
                    "minLength": 3,
                    "example": "password"
                }
            }
        },
        "web_admin.UpdateUser": {
            "type": "object",
            "required": [
                "name",
                "phone"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 120,
                    "minLength": 3,
                    "example": "John Doe"
                },
                "phone": {
                    "type": "string",
                    "maxLength": 12,
                    "minLength": 11,
                    "example": "79123456789"
                }
            }
        },
        "web_client.Session": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/domain.User"
                }
            }
        },
        "web_client.SignInUser": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 120,
                    "minLength": 3,
                    "example": "example@example.com"
                },
                "password": {
                    "type": "string",
                    "maxLength": 40,
                    "minLength": 3,
                    "example": "password"
                }
            }
        },
        "web_client.SignUpUser": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password",
                "phone"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 120,
                    "minLength": 3,
                    "example": "example@example.com"
                },
                "name": {
                    "type": "string",
                    "maxLength": 120,
                    "minLength": 3,
                    "example": "John Doe"
                },
                "password": {
                    "type": "string",
                    "maxLength": 40,
                    "minLength": 3,
                    "example": "password"
                },
                "phone": {
                    "type": "string",
                    "maxLength": 12,
                    "minLength": 11,
                    "example": "79123456789"
                }
            }
        },
        "web_client.UpdateUser": {
            "type": "object",
            "required": [
                "name",
                "phone"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 120,
                    "minLength": 3,
                    "example": "John Doe"
                },
                "phone": {
                    "type": "string",
                    "maxLength": 12,
                    "minLength": 11,
                    "example": "79123456789"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8084",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Dating API",
	Description:      "Dating api backend server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
