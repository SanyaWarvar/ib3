package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	value, err := h.services.ICacheService.GetSecret(secretHeader)
	if err != nil {
		if header == "" {
			c.JSON(http.StatusBadRequest, gin.H{"details": "bad or exp secret header"})
			return
		}
	}
	c.Set(userSecret, value)
}

/*
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid auth header")
		return
	}

	accessToken, err := h.services.IJwtManagerService.ParseToken(headerParts[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.Set(userCtx, accessToken.UserId)
}

func getUserId(c *gin.Context, withMessage bool) (uuid.UUID, error) {
	userId, ok := c.Get(userCtx)
	if !ok && withMessage {
		c.JSON(http.StatusBadRequest, gin.H{"details": "user id not found"})
		return [16]byte{}, errors.New("user id not found")
	}

	idInt, ok := userId.(uuid.UUID)
	if !ok && withMessage {
		c.JSON(http.StatusBadRequest, gin.H{"details": "user id is invalid"})
		return [16]byte{}, errors.New("user id is invalid")
	}

	return idInt, nil
}
*/
