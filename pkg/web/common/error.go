package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing/dating/api/pkg/tools"
)

// WebError godoc
// @Summary error
type WebError struct {
	Err string `json:"err"`
}

func (e WebError) Error() string {
	return e.Err
}

func SendError(g *gin.Context, err error) {
	switch err.(type) {
	case *WebError:
		g.JSON(http.StatusBadRequest, err.Error())
	case *tools.WSError:
		var e *tools.WSError
		errors.As(err, &e)
		g.JSON(e.Code, WebError{Err: e.Err})
	default:
		g.JSON(http.StatusInternalServerError, WebError{Err: err.Error()})
	}
}
