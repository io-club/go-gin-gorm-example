package util

import (
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
)

func CreateFileName(fileHeader *multipart.FileHeader) string {
	// Generate a new UUID for the filename
	id := uuid.New()
	// Get the extension of the uploaded file
	extension := filepath.Ext(fileHeader.Filename)
	// Concatenate the UUID and extension to create the filename
	return id.String() + extension
}
