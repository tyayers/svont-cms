package data

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GCSProvider struct {
	BucketName string
	BucketPath string
}

// Initialize loads persisted data structures from storage, if available
func (provider *GCSProvider) Initialize() {

	log.Printf("Initializing Google Cloud Storage data provider.")

	provider.BucketName = os.Getenv("BUCKET_NAME")
	provider.BucketPath = os.Getenv("BUCKET_PATH")
}

// Finalize writes the data structures to storage
func (provider *GCSProvider) Finalize(persistMode PersistMode, index PostIndex) {
	// Nothing special here
}

// Create a dir
func (provider *GCSProvider) CreateDir(dirName string) error {
	// No special logic needed for GCS, dirs created automatically
	return nil
}

// Uploads a file
func (provider *GCSProvider) UploadFile(fileName string, content []byte) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	buf := bytes.NewBuffer(content)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(provider.BucketName).Object(provider.BucketPath + fileName).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.

	if _, err = io.Copy(wc, buf); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	// Data can continue to be added to the file until the writer is closed.
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

// Downloads a file
func (provider *GCSProvider) DownloadFile(fileName string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(provider.BucketName).Object(provider.BucketPath + fileName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", provider.BucketPath+fileName, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return data, nil
}

// Deletes a file
func (provider *GCSProvider) DeleteFile(fileName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	b := client.Bucket(provider.BucketName)

	query := &storage.Query{Prefix: provider.BucketPath + fileName}
	var names []string
	it := b.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, attrs.Name)
		fmt.Printf("found blob " + attrs.Name)

		o := client.Bucket(provider.BucketName).Object(attrs.Name)
		if err := o.Delete(ctx); err != nil {
			return fmt.Errorf("Object(%q).Delete: %v", fileName, err)
		}

		fmt.Printf("Blob %v deleted.\n", fileName)
	}

	return nil
}
