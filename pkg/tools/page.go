package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Page - pagination params
// @Description Pagination params
type Page struct {
	// Page - page number (0,1,2...)
	Page int `json:"page" form:"page" example:"0" validate:"required"`
	// PerPage - count of objects per page (default 25)
	PerPage int `json:"per_page" form:"per_page" example:"25"`
}

// Pagination - pagination response
// @Description Pagination response
type Pagination struct {
	Page
	// Total - total count of objects
	Total      int64       `json:"total"`
	ObjectList interface{} `json:"object_list"`
}

func GetPage(g *gin.Context) Page {
	page := Page{}
	if g.BindQuery(&page) != nil || page.PerPage == 0 {
		page.Page = 0
		page.PerPage = 50
	}
	log.Info().Interface("page", page).Msg("By page")
	return page
}
