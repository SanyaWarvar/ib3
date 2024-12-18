package service

import (
	"math/big"

	"github.com/SanyaWarvar/ib3/pkg/repository"
)

type IUserService interface {
}

type IJwtManagerService interface {
}

type ICacheService interface {
	SaveSecret(secretId string, secret *big.Int) error
	GetSecret(secretId string) (*big.Int, error)
}

type Service struct {
	IUserService
	IJwtManagerService
	ICacheService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		ICacheService: NewCacheService(repos.ICacheRepo),
	}
}
