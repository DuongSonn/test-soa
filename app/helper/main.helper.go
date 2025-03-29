package helper

import (
	"sondth-test_soa/app/repository"
	"sondth-test_soa/config"
)

type HelperCollections struct {
	ProductHelper  IProductHelper
	CategoryHelper ICategoryHelper
	OAuthHelper    IOAuthHelper
	UserHelper     IUserHelper
}

func RegisterHelpers(
	postgresRepo repository.RepositoryCollections,
	config config.Configuration,
) HelperCollections {
	return HelperCollections{
		ProductHelper:  NewProductHelper(postgresRepo),
		CategoryHelper: NewCategoryHelper(postgresRepo),
		OAuthHelper:    NewOAuthHelper(config),
		UserHelper:     NewUserHelper(postgresRepo),
	}
}
