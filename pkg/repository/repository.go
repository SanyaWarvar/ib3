package repository

import (
	"math/big"

	"github.com/SanyaWarvar/ib3/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type IUserRepo interface {
	CreateUser(user models.User) error
	GetUser(username, password string) (models.User, error)
	GetUserByUP(username, hashedPassword string) (models.User, error)
	GetUserById(userId uuid.UUID) (models.User, error)
	GetUserByU(username string) (models.User, error)
	HashPassword(password string) (string, error)
	ComparePassword(password, hashedPassword string) bool
}

type ICacheRepo interface {
	SaveSecret(secretId string, secret *big.Int) error
	GetSecret(secretId string) (*big.Int, error)
}

type IJwtManagerRepo interface {
	GenerateAccessToken(userId, refreshId uuid.UUID) (string, error)
	GenerateRefreshToken(userId uuid.UUID) (string, error)
	GeneratePairToken(userId uuid.UUID) (string, string, uuid.UUID, error)
	CompareTokens(hashedToken, token string) bool
	HashToken(refreshToken string) (string, error)
	SaveRefreshToken(hashedToken string, tokenId, userId uuid.UUID) error
	DeleteRefreshTokenById(tokenId uuid.UUID) error
	GetRefreshTokenById(tokenId uuid.UUID) (string, error)
	ParseToken(accessToken string) (*models.AccessTokenClaims, error)
	CheckRefreshTokenExp(tokenId uuid.UUID) bool
}

type IReviewRepo interface {
	CreateReview(review models.Review, userId uuid.UUID) error
	EditReview(new models.ReviewInput, userId, reviewId uuid.UUID) error
	DeleteReview(reviewId, userId uuid.UUID) error

	GetAllReviewByF(filmTitle string) ([]models.ReviewOutput, error)
	GetAllReview() ([]models.Film, error)
}

type Repository struct {
	IUserRepo
	ICacheRepo
	IJwtManagerRepo
	IReviewRepo
}

func NewRepository(db *sqlx.DB, cache *redis.Client, cfg *JwtManagerCfg) *Repository {
	return &Repository{
		IUserRepo:       NewUserPostgres(db),
		ICacheRepo:      NewCacheRedis(cache),
		IJwtManagerRepo: NewJwtManagerPostgres(db, cfg),
		IReviewRepo:     NewReviewPostgres(db),
	}
}
