package repository

import (
	"math/big"

	"github.com/SanyaWarvar/ib3/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type IAuthRepo interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type ICacheRepo interface {
	SaveSecret(secretId string, secret *big.Int) error
	GetSecret(secretId string) (*big.Int, error)
}

type Repository struct {
	IAuthRepo
	ICacheRepo
}

func NewRepository(db *sqlx.DB, cache *redis.Client) *Repository {
	return &Repository{
		IAuthRepo:  NewAuthPostgres(db),
		ICacheRepo: NewCacheRedis(cache),
	}
}
