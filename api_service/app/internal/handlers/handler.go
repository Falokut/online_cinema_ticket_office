package handlers

import (
	"net/http"

	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Init(mainRouter *gin.Engine)
}

type APIHandler struct {
	handlers []Handler
}

func NewHandler(handlers []Handler) *APIHandler {
	return &APIHandler{handlers: handlers}
}

func (h *APIHandler) InitRouters(jwtHelp *jwt.Helper) http.Handler {
	router := gin.New()

	api := router.Group("/api")
	{
		api.GET("/liveness-probe", h.livenessProbe)
		api.GET("/readiness-probe", h.readinessProbe)
	}

	for _, handler := range h.handlers {
		handler.Init(router)
	}

	return router
}

func (h *APIHandler) livenessProbe(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}
func (h *APIHandler) readinessProbe(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}
