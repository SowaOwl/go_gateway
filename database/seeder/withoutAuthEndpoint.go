package seeder

import (
	"gateway/database/model"
	"gorm.io/gorm"
)

func WithoutAuthEndpointSeed(db *gorm.DB) error {
	endpoints := []string{
		"api/auth/refresh",
		"api/auth/login",
		"api/auth/temp-users/register",
		"api/auth/reset-password-check",
		"api/auth/user/reset-password",
		"api/numerator/number-generation",
	}

	for _, endpoint := range endpoints {
		seedValue := model.WithoutAuthEndpoint{
			Value: endpoint,
		}

		if err := db.Unscoped().FirstOrCreate(&seedValue, model.WithoutAuthEndpoint{Value: seedValue.Value}).Error; err != nil {
			return err
		}
	}

	return nil
}
