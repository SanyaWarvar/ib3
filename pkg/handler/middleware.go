package handler

import (
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	secretHeader        = "Secret"
	userCtx             = "userId"
	userSecret          = "userSecret"
)

func (h *Handler) userSecretIdentity(c *gin.Context) {
	header := c.GetHeader(secretHeader)

	if header == "" {
		c.JSON(http.StatusBadRequest, gin.H{"details": "empty secret header"})
		return
	}

	value, err := h.services.ICacheService.GetSecret(header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	c.Set(userSecret, value)
}

func getSecret(c *gin.Context) (*big.Int, error) {
	secretAny, ok := c.Get("userSecret")
	if !ok {
		return nil, errors.New("user secret not found")
	}

	return secretAny.(*big.Int), nil
}

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		c.JSON(http.StatusUnauthorized, "Empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		c.JSON(http.StatusUnauthorized, "Invalid auth header")
		return
	}

	accessToken, err := h.services.IJwtManagerService.ParseToken(headerParts[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.Set(userCtx, accessToken.UserId)
}

func getUserId(c *gin.Context) (uuid.UUID, error) {
	userId, ok := c.Get(userCtx)
	if !ok {
		return [16]byte{}, errors.New("user id not found")
	}

	idInt, ok := userId.(uuid.UUID)
	if !ok {
		return [16]byte{}, errors.New("user id is invalid")
	}

	return idInt, nil
}
