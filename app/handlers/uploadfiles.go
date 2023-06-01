package handlers

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func generateFileName() string {
	const (
		letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		fileNameLenght = 16
	)

	namesBytes := make([]byte, fileNameLenght)
	for i := 0; i < fileNameLenght; i++ {
		namesBytes[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(namesBytes)
}

func uploadFile(maxFileUploadSize int64, fileHeader *multipart.FileHeader, pathToSave string) (string, error) {
	fileSize := fileHeader.Size
	if fileSize > maxFileUploadSize {
		return "", errors.New("file is too large")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// detect the real file type
	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
	case "image/bmp":
		break
	default:
		return "", fmt.Errorf("invalide file type of %s, real type is %s", fileHeader.Filename, detectedFileType)
	}
	// In a real-world application, we would probably do something with the file metadata, such as saving it to a database or pushing it to an external service - in any way, we would parse and manipulate metadata. Here we create a randomized new name (this would probably be a UUID in practice) and log the future filename.

	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		return "", fmt.Errorf("cannot read the file type of %s", fileHeader.Filename)
	}
	newFileName := generateFileName() + fileEndings[0]

	newFile, err := os.Create(filepath.Join(pathToSave, newFileName))
	if err != nil {
		return "", fmt.Errorf("cannot create the file %s ", newFileName)
	}
	defer newFile.Close()

	if _, err := newFile.Write(fileBytes); err != nil {
		return "", fmt.Errorf("cannot save the file %s as a file %s", fileHeader.Filename, newFileName)
	}

	return newFileName, nil
}
