package handler

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/SanyaWarvar/ib3/pkg/models"
	"github.com/SanyaWarvar/ib3/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type signUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid json")
		return
	}
	secretAny, ok := c.Get("userSecret")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"details": "something with secret went wrong"})
		return
	}

	secret := secretAny.(*big.Int)
	utils.DecryptStruct(&input, secret.Bytes())
	fmt.Println(input)
	user := models.User{
		Id:       uuid.New(),
		Username: input.Username,
		Password: input.Password,
	}
	err := h.services.IUserService.CreateUser(user)
	if err != nil {
		errorMessage := ""
		if strings.Contains(err.Error(), "username") {
			errorMessage = "This username already exist"
		}
		c.JSON(http.StatusConflict, errorMessage)
		fmt.Println(err.Error())
		return
	}

	fmt.Println(user)

	token, refresh, _, err := h.services.IJwtManagerService.GeneratePairToken(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	output := RefreshTokenInput{
		AccessToken:  token,
		RefreshToken: refresh,
	}

	utils.EncryptStruct(&output, secret.Bytes())

	c.JSON(http.StatusCreated, output)
}

func (h *Handler) signIn(c *gin.Context) {
	var input signUpInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid json")
		return
	}
	secretAny, ok := c.Get("userSecret")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"details": "something with secret went wrong"})
		return
	}

	secret := secretAny.(*big.Int)
	fmt.Println(input.Username)
	utils.DecryptStruct(&input, secret.Bytes())

	fmt.Println(input.Username)

	user := models.User{
		Id:       uuid.New(), // значение поля не имеет разницы в контексте этого хендлера
		Username: input.Username,
		Password: input.Password,
	}

	target, err := h.services.IUserService.GetUserByUP(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	token, refresh, _, err := h.services.IJwtManagerService.GeneratePairToken(target.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	output := RefreshTokenInput{
		AccessToken:  token,
		RefreshToken: refresh,
	}

	utils.EncryptStruct(&output, secret.Bytes())

	c.JSON(http.StatusCreated, output)
}

type RefreshTokenInput struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) refreshToken(c *gin.Context) {
	var input RefreshTokenInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid json")
		return
	}
	secretAny, ok := c.Get("userSecret")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"details": "something with secret went wrong"})
		return
	}

	secret := secretAny.(*big.Int)
	utils.DecryptStruct(&input, secret.Bytes())

	accessToken, err := h.services.IJwtManagerService.ParseToken(input.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad access token")
		return
	}
	expStatus := h.services.IJwtManagerService.CheckRefreshTokenExp(accessToken.RefreshId)
	if !expStatus {
		c.JSON(http.StatusUnauthorized, "refresh token is expired or not found")
		return
	}

	compareStatus := h.services.IJwtManagerService.CompareTokens(accessToken.RefreshId, input.RefreshToken)
	if !compareStatus {
		c.JSON(http.StatusBadRequest, "invalid refresh token")
		return
	}

	err = h.services.IJwtManagerService.DeleteRefreshTokenById(accessToken.RefreshId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	token, refresh, _, err := h.services.IJwtManagerService.GeneratePairToken(accessToken.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	output := RefreshTokenInput{
		AccessToken:  token,
		RefreshToken: refresh,
	}

	utils.EncryptStruct(&output, secret.Bytes())

	c.JSON(http.StatusCreated, output)

}
