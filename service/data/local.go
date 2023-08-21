package data

import (
	"fmt"
	"log"
	"os"
)

type LocalProvider struct {
	RootDirectory string
}

// Initialize loads persisted data structures from storage, if available
func (provider *LocalProvider) Initialize() {

	log.Printf("Initializing Google Cloud Storage data provider.")

	provider.RootDirectory = os.Getenv("ROOT_DIR")
	if provider.RootDirectory == "" {
		provider.RootDirectory = "./localdata/"
	}

	os.Mkdir(provider.RootDirectory, os.ModePerm)
	os.Mkdir(provider.RootDirectory+"data/", os.ModePerm)
}

// Finalize writes the data structures to storage
func (provider *LocalProvider) Finalize(persistMode PersistMode, index PostIndex) {
	// do nothing
}

// Create a dir
func (provider *LocalProvider) CreateDir(dirName string) error {
	if err := os.Mkdir(provider.RootDirectory+dirName, os.ModePerm); err != nil {
		return err
	} else {
		return nil
	}
}

// Uploads a file
func (provider *LocalProvider) UploadFile(fileName string, content []byte) error {

	return os.WriteFile(provider.RootDirectory+fileName, content, 0644)
}

// Downloads a file
func (provider *LocalProvider) DownloadFile(fileName string) ([]byte, error) {
	return os.ReadFile(provider.RootDirectory + fileName)
}

// Deletes a file
func (provider *LocalProvider) DeleteFile(fileName string) error {
	err := os.RemoveAll(provider.RootDirectory + fileName)

	if err != nil {
		fmt.Printf("Error deleting post: %s", err)
		return err
	} else {
		return nil
	}
}
