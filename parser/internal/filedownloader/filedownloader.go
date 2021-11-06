package filedownloader

import (
	"io"
)

// FileDownloader interface describes the methods needed to list and download files to process.
type FileDownloader interface {
	ListFileNames() []string
	DownloadFile(fileName string) io.ReadCloser
}
