package handler

import (
	"math/big"
	"net/http"

	"github.com/SanyaWarvar/ib3/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) encrypt(c *gin.Context) {
	var input signUpInput
	err := c.BindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}
	secretAny, ok := c.Get("userSecret")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"details": "something with secret went wrong"})
		return
	}

	secret := secretAny.(*big.Int)
	utils.EncryptStruct(&input, secret.Bytes())
	c.JSON(200, input)
}

func (h *Handler) decrypt(c *gin.Context) {
	var input signUpInput
	err := c.BindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}
	secretAny, ok := c.Get("userSecret")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"details": "something with secret went wrong"})
		return
	}

	secret := secretAny.(*big.Int)

	utils.DecryptStruct(&input, secret.Bytes())
	c.JSON(200, input)

}
