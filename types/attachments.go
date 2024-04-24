package types

import (
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"gorm.io/gorm"
)

type Attachment struct {
	gorm.Model
	Filename string `gorm:"-"`
	Dest     string `gorm:"unique"`
	Url      string `gorm:"-"`
	BookID   uint
}

func (a *Attachment) Fill(bucket string) {
	url, filename := utils.GetUrlAndFilename(bucket, a.Dest)
	a.Filename = filename
	a.Url = url
}
