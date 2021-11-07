package logparser

import (
	"io"
	"log"
	"os"
)

type MockFileDownloader struct {
	FileNameToDownload string
}

func (mock *MockFileDownloader) ListFileNames() []string {
	return []string{mock.FileNameToDownload}
}

func (mock *MockFileDownloader) DownloadFile(fileName string) io.ReadCloser {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Could not open given log file. " + fileName)
	}

	return f
}
