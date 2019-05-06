package service

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"log"
	"net/url"
	"time"
)

type AzureStorageService struct {
	AzContainerURL *azblob.ContainerURL
	AccountName    string
	AccountKey     string
	ContainerName  string
	Credential     *azblob.SharedKeyCredential
}

func NewAzureStorageService(accountName string, accountKey string, containerName string, ctx context.Context) *AzureStorageService {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	return &AzureStorageService{
		AccountName:    accountName,
		AccountKey:     accountKey,
		ContainerName:  containerName,
		AzContainerURL: createContainerURL(accountName, accountKey, containerName, ctx),
		Credential:     createCredential(accountName, accountKey),
	}
}

func createContainerURL(accountName string, accountKey string, containerName string, ctx context.Context) *azblob.ContainerURL {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	credential, _ := azblob.NewSharedKeyCredential(accountName, accountKey)
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*URL, p)
	containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)

	return &containerURL
}

func createCredential(accountName string, accountKey string) *azblob.SharedKeyCredential {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	res, _ := azblob.NewSharedKeyCredential(accountName, accountKey)
	return res
}

func (azStorageService *AzureStorageService) Upload(content []byte, ctx context.Context, pin string, fileName string) (azblob.CommonResponse, error, string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	blobName := fmt.Sprintf("%s-%s", pin, fileName)
	blobURL := azStorageService.AzContainerURL.NewBlockBlobURL(blobName)
	uploadOptions := azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16,
	}

	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		StartTime:     time.Now().UTC(),
		ExpiryTime:    time.Now().UTC().Add(168 * time.Hour),
		ContainerName: azStorageService.ContainerName,
		BlobName:      blobName,
		Permissions:   azblob.BlobSASPermissions{Add: false, Read: true, Write: false}.String(),
	}.NewSASQueryParameters(azStorageService.Credential)
	qp := sasQueryParams.Encode()
	publicUrl := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
		azStorageService.AccountName, azStorageService.ContainerName, blobName, qp)

	res, err := azblob.UploadBufferToBlockBlob(ctx, content, blobURL, uploadOptions)

	return res, err, publicUrl
}

func (azStorageService *AzureStorageService) Delete(ctx context.Context, fileName string) (*azblob.BlobDeleteResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	blobUrl := azStorageService.AzContainerURL.NewBlockBlobURL(fileName)
	log.Println(blobUrl)
	return blobUrl.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
}
