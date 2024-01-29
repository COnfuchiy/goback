package services

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"goback/api/response"
	"goback/domain/entity"
	"goback/repository"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type IFileService interface {
	Create(file *entity.File) error
	GetFromContext(contextFileID string) (*entity.File, error)
	CheckExisting(filename string) (bool, error)
	GetFileTag(file *entity.File, fileObject *multipart.FileHeader) error
}

type FileService struct {
	fileRepository repository.IFileRepository
}

func NewFileService(fileRepository repository.IFileRepository) IFileService {
	return &FileService{fileRepository: fileRepository}
}

func (s FileService) Create(file *entity.File) error {
	return s.fileRepository.Create(file)
}

func (s FileService) GetFromContext(contextFileID string) (*entity.File, error) {
	fileID, err := uuid.Parse(contextFileID)
	if err != nil {
		return nil, err
	}
	file, err := s.fileRepository.FindByID(fileID)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s FileService) CheckExisting(filename string) (bool, error) {
	return s.fileRepository.CheckExisting(filename)
}

func (s FileService) GetFileTag(file *entity.File, fileObject *multipart.FileHeader) error {
	if ext := filepath.Ext(fileObject.Filename); ext == "docx" || ext == "pdf" {

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		fileHandler, err := fileObject.Open()
		if err != nil {
			return err
		}

		fileWriter, err := writer.CreateFormFile("fileHandler", fileObject.Filename)
		if err != nil {
			return err
		}

		_, err = io.Copy(fileWriter, fileHandler)
		if err != nil {
			return err
		}

		_ = writer.Close()

		req, err := http.NewRequest("POST", "http://62.84.119.63/uploadfile/", &requestBody)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		content, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		textExtractorResponse := response.TextExtractorResponse{}
		err = json.Unmarshal(content, &textExtractorResponse)
		if err != nil {
			return err
		}

		if textExtractorResponse.Success {
			tagForm := url.Values{}
			tagForm.Add("file_text", textExtractorResponse.FileText)
			newRequest, err := http.NewRequest("POST", "http://localhost:8234/tag/", strings.NewReader(tagForm.Encode()))
			if err != nil {
				return err
			}
			newRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			res, err = client.Do(req)
			if err != nil {
				return err
			}
			content, err = io.ReadAll(res.Body)
			if err != nil {
				return err
			}

			taggingResponse := response.TaggingResponse{}
			err = json.Unmarshal(content, &taggingResponse)
			if err != nil {
				return err
			}

			if taggingResponse.Success {
				file.Tag = taggingResponse.FileTag
			}
		}
	}
	return nil
}
