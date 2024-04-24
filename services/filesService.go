package services

import (
	"mime/multipart"
	"strconv"

	"github.com/DeepAung/gofiber-library/pkg/storer"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/DeepAung/gofiber-library/types"
	"gorm.io/gorm"
)

type FilesService struct {
	storer storer.Storer
	db     *gorm.DB
}

func NewFilesService(db *gorm.DB, storer storer.Storer) *FilesService {
	return &FilesService{
		db:     db,
		storer: storer,
	}
}

func (h *FilesService) UploadFiles(
	files []*multipart.FileHeader,
	bookId int,
) ([]*types.Attachment, error) {
	dir := utils.Join("books", strconv.Itoa(bookId))
	results, err := h.storer.UploadFiles(files, dir)
	if err != nil {
		return []*types.Attachment{}, err
	}

	attachments := make([]*types.Attachment, len(results))
	for i, result := range results {
		attachments[i] = &types.Attachment{
			Filename: result.Filename,
			Dest:     utils.Join(dir, result.Filename),
			Url:      result.Url,
			BookID:   uint(bookId),
		}
	}

	return attachments, h.db.Model(&types.Attachment{}).Create(attachments).Error
}

func (h *FilesService) DeleteFiles(dests []string) error {
	if err := h.storer.DeleteFiles(dests); err != nil {
		return err
	}

	// TODO: optimize this
	return h.db.Transaction(func(tx *gorm.DB) error {
		for _, dest := range dests {
			if err := tx.Unscoped().
				Model(&types.Attachment{}).
				Where("dest = ?", dest).
				Delete(&types.Attachment{}).
				Error; err != nil {
				return err
			}
		}
		return nil
	})
}
