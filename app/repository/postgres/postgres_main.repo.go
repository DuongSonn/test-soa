package postgres

import (
	"gorm.io/gorm"

	"sondth-test_soa/app/repository"
)

func RegisterPostgresRepositories(db *gorm.DB) repository.RepositoryCollections {
	return repository.RepositoryCollections{
		ProductRepo:  NewPostgresProductRepository(db),
		CategoryRepo: NewPostgresCategoryRepository(db),
		UserRepo:     NewPostgresUserRepository(db),
		ReviewRepo:   NewPostgresReviewRepository(db),
		WishlistRepo: NewPostgresWishlistRepository(db),
	}
}
