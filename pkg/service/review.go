package service

import (
	"github.com/SanyaWarvar/ib3/pkg/models"
	"github.com/SanyaWarvar/ib3/pkg/repository"
	"github.com/google/uuid"
)

type ReviewService struct {
	repo repository.IReviewRepo
}

func NewReviewService(repo repository.IReviewRepo) *ReviewService {
	return &ReviewService{repo: repo}
}

func (r *ReviewService) CreateReview(review models.Review, userId uuid.UUID) error {
	return r.repo.CreateReview(review, userId)
}

func (r *ReviewService) EditReview(new models.ReviewInput, userId, reviewId uuid.UUID) error {
	return r.repo.EditReview(new, userId, reviewId)
}

func (r *ReviewService) DeleteReview(reviewId, userId uuid.UUID) error {
	return r.repo.DeleteReview(reviewId, userId)
}

func (r *ReviewService) GetAllReviewByF(filmTitle string) ([]models.ReviewOutput, error) {
	return r.repo.GetAllReviewByF(filmTitle)
}

func (r *ReviewService) GetAllReview() ([]models.Film, error) {
	return r.repo.GetAllReview()
}
