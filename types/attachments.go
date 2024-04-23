package types

import (
	"gorm.io/gorm"
)

// TODO: we don't need to store filename and url
type Attachment struct {
	gorm.Model
	Filename string `gorm:"unique"`
	Dest     string `gorm:"unique"`
	Url      string `gorm:"unique"`
	BookID   uint
}
