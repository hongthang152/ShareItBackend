package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"context"
	"github.com/hongthang152/ShareItBackend/models"
	"github.com/hongthang152/ShareItBackend/repository"
	file "github.com/hongthang152/ShareItBackend/repository/file"
	service "github.com/hongthang152/ShareItBackend/service"
	"github.com/hongthang152/ShareItBackend/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileHandler struct {
	repo      repository.FileRepo
	azService *service.AzureStorageService
}

func NewFileHandler(client *mongo.Client) *FileHandler {
	return &FileHandler{
		repo:      file.NewMongoRepository(client),
		azService: service.NewAzureStorageService(os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY"), os.Getenv("CONTAINER_NAME"), context.Background()),
	}
}

func (fileHandler *FileHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	utils.RespondwithJSON(w, http.StatusOK, "Things are working 😃")
}

func (fileHandler *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	codeParam, _ := r.URL.Query()["code"]
	code := codeParam[0]
	file, _ := fileHandler.repo.FetchByPIN(r.Context(), code)

	if file == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Cannot find file with this ID")
		return
	}

	utils.RespondwithJSON(w, http.StatusOK, *file)
}

func (fileHandler *FileHandler) Create(w http.ResponseWriter, r *http.Request) {

	codeParam, _ := r.URL.Query()["code"]
	pin := codeParam[0]

	var Buf bytes.Buffer
	file, header, _ := r.FormFile("file")
	defer file.Close()
	io.Copy(&Buf, file)

	res, _ := fileHandler.repo.FetchByPIN(r.Context(), pin)
	if res != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Somebody has taken this ID. Please try a different ID.")
		return
	}
	_, _, url := fileHandler.azService.Upload(Buf.Bytes(), r.Context(), pin, header.Filename)

	fileHandler.repo.Create(r.Context(), models.File{Name: fmt.Sprintf("%s-%s", pin, header.Filename), Pin: pin, Url: url})
	utils.RespondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
}

func (fileHandler *FileHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	files, _ := fileHandler.repo.FetchAll(r.Context())
	utils.RespondwithJSON(w, http.StatusCreated, files)
}
