package seeder

import "gorm.io/gorm"

func Seed(db *gorm.DB) error {
	err := WithoutAuthEndpointSeed(db)
	if err != nil {
		return err
	}

	return nil
}
