package uploads

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const uploadDir = "static/uploads"

func Save(file multipart.File, header *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create upload directory: %w", err)
	}

	ext := filepath.Ext(header.Filename)
	// Utilizes time.Now... call to specify a unique upload timestamp
	// this timestamp is used to ensure the uniqueness of filename. 
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dest := filepath.Join(uploadDir, filename)

	out, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("could not create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("could not save file: %w", err)
	}

	return "/static/uploads/" + filename, nil
}