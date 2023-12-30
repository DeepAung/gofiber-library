package models

import "gorm.io/gorm"

func InitModel(db *gorm.DB) error {
	return db.SetupJoinTable(&User{}, "FavBooks", &UserFavbooks{})
}
