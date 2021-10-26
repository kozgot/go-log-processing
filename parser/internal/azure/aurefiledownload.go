package azure

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"

	azblob "github.com/Azure/azure-storage-blob-go/azblob"
)

func handleErrors(err error) {
	if err != nil {
		// todo error handling
		/*
			if serr, ok := err.(azblob.StorageError); ok { // This error is a Service-specific
				switch serr.ServiceCode() { // Compare serviceCode to ServiceCodeXxx constants
				case azblob.ServiceCodeContainerAlreadyExists:
					fmt.Println("Received 409. Container already exists")
					return
				}
			}
		*/
		log.Fatal(err)
	}
}

func Cucc() {
	fmt.Printf("Azure Blob storage quick start sample\n")

	// From the Azure portal, get your storage account name and key and set environment variables.
	// accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")

	// todo: add these to environment variables
	accountName := "kozgotstorage"
	accountKey := "ET1z3fA9QVNK5sbZ/aH7cootN3f8R4qnUQyfSsAIaBLl7NBjffXEYJ4dIN7r76PFSKxaQ5Vew2YEpu6EdEU9Cw=="
	myContainerName := "testcontainer"

	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create the container
	ctx := context.Background() // This example uses a never-expiring context

	fileNames := []string{}

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL2, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, myContainerName))

	containerURL2 := azblob.NewContainerURL(*URL2, p)

	// List the container that we have created above
	fmt.Println("Listing the blobs in the container:")
	for marker := (azblob.Marker{}); marker.NotDone(); {
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL2.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		handleErrors(err)

		// ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			fmt.Print("	Blob name: " + blobInfo.Name + "\n")
			fileNames = append(fileNames, blobInfo.Name)
		}
	}

	if len(fileNames) == 0 {
		fmt.Println("no files found")
	}

	fileToDownload := fileNames[0]
	blobURL := containerURL2.NewBlockBlobURL(fileToDownload)

	// Here's how to download the blob
	downloadResponse, err := blobURL.Download(
		ctx,
		0,
		azblob.CountToEnd,
		azblob.BlobAccessConditions{},
		false,
		azblob.ClientProvidedKeyOptions{})

	handleErrors(err)

	// NOTE: automatically retries are performed if the connection fails
	maxRetries := 20
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: maxRetries})

	// todo
	scanner := bufio.NewScanner(bodyStream)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	/*
		// read the body into a buffer
		downloadedData := bytes.Buffer{}
		_, err = downloadedData.ReadFrom(bodyStream)
		handleErrors(err)

		// The downloaded blob data is in downloadData's buffer. :Let's print it
		fmt.Printf("Downloaded the blob: " + downloadedData.String())
	*/

	/*
		// Cleaning up the quick start by deleting the container and the file created locally
		fmt.Printf("Press enter key to delete the sample files, example container, and exit the application.\n")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Printf("Cleaning up.\n")

			containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
			file.Close()
			os.Remove(fileName)
	*/
}
