package service

import (
	"math/big"

	"github.com/SanyaWarvar/ib3/pkg/repository"
)

type CacheService struct {
	cache repository.ICacheRepo
}

func NewCacheService(cache repository.ICacheRepo) *CacheService {
	return &CacheService{cache: cache}
}

func (s *CacheService) SaveSecret(secretId string, secret *big.Int) error {
	return s.cache.SaveSecret(secretId, secret)
}
func (s *CacheService) GetSecret(secretId string) (*big.Int, error) {
	return s.cache.GetSecret(secretId)
}
