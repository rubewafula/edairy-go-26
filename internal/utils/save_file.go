package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

func SaveFile(file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	os.MkdirAll("uploads/"+folder, os.ModePerm)

	name := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	path := "uploads/" + folder + "/" + name

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return path, err
}
