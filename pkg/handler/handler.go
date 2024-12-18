package handler

import (
	"math/big"

	"github.com/SanyaWarvar/ib3/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	secrets  map[string]*big.Int
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services, secrets: map[string]*big.Int{}}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/keys", h.keys)

	auth := router.Group("/auth", h.userSecretIdentity)
	{
		auth.POST("/sign_up", h.sign_up)
	}
	return router
}
