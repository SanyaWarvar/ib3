package service

import (
	"math/big"

	"github.com/SanyaWarvar/ib3/pkg/models"
	"github.com/SanyaWarvar/ib3/pkg/repository"
	"github.com/google/uuid"
)

type IUserService interface {
	CreateUser(models.User) error
	HashPassword(password string) (string, error)
	GetUserByUP(user models.User) (models.User, error)
	GetUserByU(username string) (models.User, error)
}

type IJwtManagerService interface {
	GeneratePairToken(userId uuid.UUID) (string, string, uuid.UUID, error)
	CompareTokens(refreshId uuid.UUID, token string) bool
	SaveRefreshToken(hashedToken string, userId, tokenId uuid.UUID) error
	DeleteRefreshTokenById(tokenId uuid.UUID) error
	GetRefreshTokenById(tokenId uuid.UUID) (string, error)
	ParseToken(accessToken string) (*models.AccessTokenClaims, error)
	CheckRefreshTokenExp(tokenId uuid.UUID) bool
}

type ICacheService interface {
	SaveSecret(secretId string, secret *big.Int) error
	GetSecret(secretId string) (*big.Int, error)
}

type IReviewService interface {
	CreateReview(review models.Review, userId uuid.UUID) error
	EditReview(new models.ReviewInput, userId, reviewId uuid.UUID) error
	DeleteReview(reviewId, userId uuid.UUID) error

	GetAllReviewByF(filmTitle string) ([]models.ReviewOutput, error)
	GetAllReview() ([]models.Film, error)
}

type Service struct {
	IUserService
	IJwtManagerService
	ICacheService
	IReviewService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		IUserService:       NewUserService(repos.IUserRepo),
		ICacheService:      NewCacheService(repos.ICacheRepo),
		IJwtManagerService: NewJwtManagerService(repos.IJwtManagerRepo),
		IReviewService:     NewReviewService(repos.IReviewRepo),
	}
}
