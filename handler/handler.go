package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/proggcreator/WbWorkDb/repository"
)

type Handler struct {
	repository *repository.Repository
}

func NewHandler(repository *repository.Repository) *Handler {
	return &Handler{repository: repository}

}
func (h *Handler) InitRoutes() *gin.Engine { //эндпоинты

	router := gin.New()
	res := router.Group("/api")
	{
		res.POST("/do", h.do_something)
		res.GET("/:id", h.request_error_code)
	}
	return router
}
