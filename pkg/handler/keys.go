package handler

import (
	"math/big"
	"net/http"

	dh "github.com/SanyaWarvar/ib3/pkg/dh"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type KeysInput struct {
	P               string `json:"p" binding:"required"`
	G               string `json:"g" binding:"required"`
	ClientPublicKey string `json:"public_key" binding:"required"`
}

func (h *Handler) keys(c *gin.Context) {
	var input KeysInput
	c.BindJSON(&input)
	bigP, ok := new(big.Int).SetString(input.P, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"details": "incorrect p value"})
		return
	}
	bigG, ok := new(big.Int).SetString(input.G, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"details": "incorrect g value"})
		return
	}
	BigClientPublicKey, ok := new(big.Int).SetString(input.ClientPublicKey, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"details": "incorrect public_key value"})
		return
	}

	private, public, err := dh.GenerateKeys(bigP, bigG)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}
	sharedSecret := dh.ComputeSharedSecret(BigClientPublicKey, private, bigP)
	secretId := uuid.New()
	err = h.services.ICacheService.SaveSecret(secretId.String(), sharedSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"public_key": public.String(), "secret_id": secretId})
}
