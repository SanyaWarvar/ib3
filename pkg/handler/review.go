package handler

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SanyaWarvar/ib3/pkg/models"
	"github.com/SanyaWarvar/ib3/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) createReview(c *gin.Context) {
	var input models.ReviewInputForHandler
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	secret, err := getSecret(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	err = utils.DecryptStruct(&input, secret.Bytes())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	ratingInt, err := strconv.Atoi(input.Rating)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}
	review := models.Review{
		Id:        uuid.New(),
		Body:      input.Body,
		Rating:    ratingInt,
		FilmTitle: input.FilmTitle,
	}

	err = h.services.CreateReview(review, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}
	encId, err := utils.Encrypt([]byte(review.Id.String()), secret.Bytes())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}
	encIdStr := hex.EncodeToString(encId)
	c.JSON(http.StatusCreated, gin.H{"review_id": encIdStr})
}

func (h *Handler) editReview(c *gin.Context) {
	var input models.ReviewInputForHandler
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	secret, err := getSecret(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	err = utils.DecryptStruct(&input, secret.Bytes())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	review := models.ReviewInput{
		Body:   &input.Body,
		Rating: new(int),
	}

	if input.Body == "" {
		review.Body = nil
	}
	fmt.Print(input.Rating == "")
	if input.Rating == "" {
		review.Rating = nil
	} else {
		ratingInt, err := strconv.Atoi(input.Rating)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
			return
		}
		review.Rating = &ratingInt
	}

	reviewIdStr := c.Param("id")
	reviewId, err := uuid.Parse(reviewIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	err = h.services.IReviewService.EditReview(review, userId, reviewId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) deleteReview(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	reviewIdStr := c.Param("id")
	reviewId, err := uuid.Parse(reviewIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	err = h.services.IReviewService.DeleteReview(reviewId, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) getAllReviews(c *gin.Context) {

	output, err := h.services.IReviewService.GetAllReview()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	for ind, item := range output {
		output[ind].CoverUrl = c.Request.Host + "/covers/" + item.CoverUrl
	}

	c.JSON(http.StatusOK, output)
}

func (h *Handler) getFilmReview(c *gin.Context) {
	title := c.Param("title")

	output, err := h.services.IReviewService.GetAllReviewByF(title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
