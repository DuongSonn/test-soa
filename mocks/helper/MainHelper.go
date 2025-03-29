package mocks

import (
	"testing"

	"sondth-test_soa/app/helper"
)

func InitMockHelper(t *testing.T) helper.HelperCollections {
	return helper.HelperCollections{
		CategoryHelper: NewICategoryHelper(t),
		ProductHelper:  NewIProductHelper(t),
		UserHelper:     NewIUserHelper(t),
		OAuthHelper:    NewIOAuthHelper(t),
	}
}
