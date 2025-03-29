package mocks

import (
	"testing"

	"sondth-test_soa/app/repository"
)

func InitMockRepository(t *testing.T) repository.RepositoryCollections {
	return repository.RepositoryCollections{
		CategoryRepo: NewICategoryRepository(t),
		ProductRepo:  NewIProductRepository(t),
		ReviewRepo:   NewIReviewRepository(t),
		WishlistRepo: NewIWishlistRepository(t),
		UserRepo:     NewIUserRepository(t),
	}
}
