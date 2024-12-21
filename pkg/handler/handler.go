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
	router.Static("/covers", "./covers")
	crp := router.Group("/api", h.userSecretIdentity)
	{
		crp.POST("/encrypt", h.encrypt)
		crp.POST("/decrypt", h.decrypt)
	}

	auth := router.Group("/auth", h.userSecretIdentity)
	{
		auth.POST("/sign_up", h.signUp)
		auth.POST("/sign_in", h.signIn)
		auth.POST("/refresh", h.refreshToken)
	}

	router.GET("api/review/all", h.getAllReviews)
	router.GET("api/review/film/:title", h.getFilmReview)
	review := crp.Group("/review", h.userIdentity)
	{
		review.POST("/", h.createReview)
		review.PUT("/:id", h.editReview)
		review.DELETE("/:id", h.deleteReview)
	}
	return router
}
