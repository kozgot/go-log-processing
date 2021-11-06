package filedownloader

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"

	azblob "github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/kozgot/go-log-processing/parser/internal/utils"
)

// AzureDownloader contains data needed to list or dowload blobs from azure.
type AzureDownloader struct {
	Credential        *azblob.SharedKeyCredential
	StorageAccountURL *url.URL
	ContainerURL      azblob.ContainerURL
}

// NewAzureDownloader creates and returns data a DownloaderData.
func NewAzureDownloader(accountName string, accountKey string, containerName string) *AzureDownloader {
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	utils.FailOnError(err, "Invalid credentials for azure")

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your storage account blob service storageAccountURL endpoint.
	storageAccountURL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*storageAccountURL, pipeline)

	downloader := AzureDownloader{
		Credential:        credential,
		ContainerURL:      containerURL,
		StorageAccountURL: storageAccountURL,
	}

	return &downloader
}

// ListFileNames lists the blobs in the azure container.
func (downloader *AzureDownloader) ListFileNames() []string {
	fileNames := []string{}
	ctx := context.Background()

	// List the container that we have created above
	log.Println("Listing the blobs in the container:")
	for marker := (azblob.Marker{}); marker.NotDone(); {
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := downloader.ContainerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		utils.FailOnError(err, "Could not list blobs")

		// ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			log.Print("	Blob name: " + blobInfo.Name + "\n")
			fileNames = append(fileNames, blobInfo.Name)
		}
	}

	if len(fileNames) == 0 {
		log.Println("No files found in Azure blob storage container.")
	}

	return fileNames
}

// DownloadFile downloads the blob with the given name from azure.
func (downloader *AzureDownloader) DownloadFile(fileName string) io.ReadCloser {
	blobURL := downloader.ContainerURL.NewBlockBlobURL(fileName)
	ctx := context.Background()

	// Download the blob
	downloadResponse, err := blobURL.Download(
		ctx,
		0,
		azblob.CountToEnd,
		azblob.BlobAccessConditions{},
		false,
		azblob.ClientProvidedKeyOptions{})

	utils.FailOnError(err, "Could not download blob")

	// Automatic retries are performed if the connection fails
	maxRetries := 20
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: maxRetries})

	return bodyStream
}
