package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type signUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) sign_up(c *gin.Context) {
	var input signUpInput
	c.BindJSON(&input)
	fmt.Println(c.Get("userSecret"))
}
