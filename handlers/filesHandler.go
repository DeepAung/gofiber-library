package handlers

import (
	"strconv"

	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/DeepAung/gofiber-library/services"
	"github.com/gofiber/fiber/v2"
)

type FilesHandler struct {
	filesSvc *services.FilesService
}

func NewFilesHandler(filesSvc *services.FilesService) *FilesHandler {
	return &FilesHandler{
		filesSvc: filesSvc,
	}
}

func (h *FilesHandler) UploadFiles(c *fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("bookId"))
	if err != nil {
		return utils.RenderErrorText(c, "bookId is not an integer")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return utils.RenderErrorText(c, "c.MultipartForm: "+err.Error())
	}
	files, ok := form.File["files"]
	if !ok {
		return utils.RenderErrorText(c, "\"files\" field not found")
	}

	attachments, err := h.filesSvc.UploadFiles(files, bookId)
	if err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	return c.Render("components/images", &fiber.Map{"Attachments": attachments})
}

func (h *FilesHandler) DeleteFile(c *fiber.Ctx) error {
	dest := c.FormValue("dest")
	if err := h.filesSvc.DeleteFiles([]string{dest}); err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
