package schedulers

import (
	"context"
	"github.com/hongthang152/ShareItBackend/repository"
	"github.com/hongthang152/ShareItBackend/repository/file"
	"github.com/hongthang152/ShareItBackend/service"
	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

type FileScheduler struct {
	repo      repository.FileRepo
	azService *service.AzureStorageService
}

func NewFileScheduler(client *mongo.Client) *FileScheduler {
	return &FileScheduler{
		repo:      file.NewMongoRepository(client),
		azService: service.NewAzureStorageService(os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY"), os.Getenv("CONTAINER_NAME"), context.Background()),
	}
}

func (fileScheduler *FileScheduler) DeleteExpiredFiles(spec string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	c := cron.New()
	c.AddFunc("*/10 * * * * *", func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		files, _ := fileScheduler.repo.FetchAllOneWeekExpired(context.Background())
		for _, file := range files {
			fileScheduler.azService.Delete(context.Background(), file.Name)
		}
		fileScheduler.repo.DeleteAllOneWeekExpired(context.Background())
	})
	c.Start()
}
