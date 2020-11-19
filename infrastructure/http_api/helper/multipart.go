package helper

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/market-place/usecase/usecase_error"
)

const BUCKET string = "ecommerce_s2l_assets"
const STORAGE_URL = "https://storage.googleapis.com"

type multipartFile struct {
	file     multipart.File
	files    []*multipart.FileHeader
	gStorage *storage.Client
}

func NewMultiPart(gStorage *storage.Client) *multipartFile {
	return &multipartFile{
		gStorage: gStorage,
	}
}

func (m *multipartFile) ReadAvatar(r *http.Request) error {
	fieldName := "avatar"
	file, _, err := r.FormFile(fieldName)
	if err != nil {
		return err
	}
	m.file = file

	return nil
}

func (m *multipartFile) ReadPhotos(r *http.Request) error {
	fieldName := "photos"
	files := r.MultipartForm.File[fieldName]
	for _, file := range files {
		m.files = append(m.files, file)
	}
	return nil
}

func (m *multipartFile) IsImageAvatar() (bool, error) {
	fileHeader := make([]byte, 512)
	if _, err := m.file.Read(fileHeader); err != nil {
		return false, err
	}
	if _, err := m.file.Seek(0, 0); err != nil {
		return false, err
	}

	re := regexp.MustCompile("^image/.*$")
	if isImage := re.MatchString(http.DetectContentType(fileHeader)); !isImage {
		return false, nil
	}

	return true, nil
}

func (m *multipartFile) IsImagePhotos() (bool, error) {
	for _, fileHeader := range m.files {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Printf("[MULTIPART] : Open Photos %#v \n", err.Error)
			return false, usecase_error.ErrInternalServerError
		}

		fileHeader := make([]byte, 512)
		if _, err := file.Read(fileHeader); err != nil {
			return false, err
		}
		if _, err := file.Seek(0, 0); err != nil {
			return false, err
		}

		re := regexp.MustCompile("^image/.*$")
		if isImage := re.MatchString(http.DetectContentType(fileHeader)); !isImage {
			return false, nil
		}
	}

	return true, nil
}

func (m *multipartFile) StorePhoto(ctx context.Context) (string, error) {

	fileName := fmt.Sprintf("assets-%d.png", time.Now().UnixNano())

	sw := m.gStorage.Bucket(BUCKET).Object(fileName).NewWriter(ctx)

	if _, err := io.Copy(sw, m.file); err != nil {
		fmt.Printf("[MultiPartHelper] : Copy file to storage writer %#v \n", err)
		return "", usecase_error.ErrInternalServerError
	}
	if err := sw.Close(); err != nil {
		fmt.Printf("[MultiPartHelper] : Closestorage writer %#v \n", err)
		return "", usecase_error.ErrInternalServerError
	}

	urlImage := fmt.Sprintf("%s/%s/%s", STORAGE_URL, BUCKET, fileName)

	return urlImage, nil
}

func (m *multipartFile) StorePhotos(ctx context.Context) ([]string, error) {
	var files []string
	var wgUploadFiles sync.WaitGroup

	type response struct {
		Url string
		Err error
	}

	responses := []response{}
	for _, header := range m.files {
		wgUploadFiles.Add(1)
		go func(h *multipart.FileHeader) {
			defer wgUploadFiles.Done()

			file, err := h.Open()
			if err != nil {
				fmt.Printf("[MULTIPART] : Open Photos %s \n", err.Error)
				responses = append(responses, response{
					Err: err,
				})
				return
			}

			fileName := fmt.Sprintf("assets-%d.png", time.Now().UnixNano())

			sw := m.gStorage.Bucket(BUCKET).Object(fileName).NewWriter(ctx)
			if _, err := io.Copy(sw, file); err != nil {
				fmt.Printf("[MultiPartHelper] : Copy file to storage writer %s \n", err)
				responses = append(responses, response{
					Err: err,
				})
				return
			}

			if err := sw.Close(); err != nil {
				fmt.Printf("[MultiPartHelper] : Closestorage writer %s \n", err)
				responses = append(responses, response{
					Err: err,
				})
				return
			}

			urlImage := fmt.Sprintf("%s/%s/%s", STORAGE_URL, BUCKET, fileName)
			responses = append(responses, response{
				Url: urlImage,
			})
		}(header)
	}

	wgUploadFiles.Wait()
	for _, response := range responses {
		if response.Err != nil {
			return files, response.Err
		}
		files = append(files, response.Url)
	}

	return files, nil
}
