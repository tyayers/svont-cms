package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GCSProvider struct {
	IndexWriteMutex sync.Mutex
	BucketName      string
	BucketPath      string
}

// Initialize loads persisted data structures from storage, if available
func (provider *GCSProvider) Initialize() (map[string]PostOverview, []string, map[int][]string, map[string]map[int]string) {

	log.Printf("Initializing Google Cloud Storage data provider.")

	provider.BucketName = os.Getenv("BUCKET_NAME")
	provider.BucketPath = os.Getenv("BUCKET_PATH")

	var index_main map[string]PostOverview
	var index_time []string
	var index_popularity map[int][]string
	var index_tags map[string]map[int]string

	postBytes, err := downloadFileIntoMemory(provider.BucketName, provider.BucketPath+"index.json")

	if err == nil {
		json.Unmarshal(postBytes, &index_main)
	}

	if index_main == nil {
		index_main = map[string]PostOverview{}
	}

	postBytes, err = downloadFileIntoMemory(provider.BucketName, provider.BucketPath+"index_time.json")

	if err == nil {
		json.Unmarshal(postBytes, &index_time)
	}

	if index_time == nil {
		index_time = []string{}
	}

	postBytes, err = downloadFileIntoMemory(provider.BucketName, provider.BucketPath+"index_popularity.json")

	if err == nil {
		json.Unmarshal(postBytes, &index_popularity)
	}

	if index_popularity == nil {
		index_popularity = map[int][]string{}
		index_popularity[0] = []string{}
	}

	postBytes, err = downloadFileIntoMemory(provider.BucketName, provider.BucketPath+"index_tags.json")

	if err == nil {
		json.Unmarshal(postBytes, &index_tags)
	}

	if index_tags == nil {
		index_tags = map[string]map[int]string{}
	}

	provider.IndexWriteMutex = sync.Mutex{}

	return index_main, index_time, index_popularity, index_tags

}

// Finalize writes the data structures to storage
func (provider *GCSProvider) Finalize(index_main map[string]PostOverview, index_time []string, index_popularity map[int][]string, index_tags map[string]map[int]string) {

	jsonData, err := json.Marshal(index_main)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	streamFileUpload(provider.BucketName, provider.BucketPath+"index.json", jsonData)

	jsonData, err = json.Marshal(index_time)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	streamFileUpload(provider.BucketName, provider.BucketPath+"index_time.json", jsonData)

	jsonData, err = json.Marshal(index_popularity)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	streamFileUpload(provider.BucketName, provider.BucketPath+"index_popularity.json", jsonData)

	jsonData, err = json.Marshal(index_tags)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	streamFileUpload(provider.BucketName, provider.BucketPath+"index_tags.json", jsonData)
}

// Returns the post specified by postId.
func (provider *GCSProvider) GetPost(postId string) *Post {

	dat, _ := downloadFileIntoMemory(provider.BucketName, provider.BucketPath+postId+"/post.json")

	var post Post
	json.Unmarshal(dat, &post)

	return &post
}

// Creates a new post.
func (provider *GCSProvider) CreatePost(newPost Post, fileAttachments map[string][]byte) error {

	jsonData, _ := json.Marshal(newPost)
	err := streamFileUpload(provider.BucketName, provider.BucketPath+newPost.Header.Id+"/post.json", jsonData)

	for k, v := range fileAttachments {
		err = streamFileUpload(provider.BucketName, provider.BucketPath+newPost.Header.Id+"/"+k, v)
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Updates a post
func (provider *GCSProvider) UpdatePost(post Post, fileAttachments map[string][]byte) error {

	jsonData, _ := json.Marshal(post)
	err := streamFileUpload(provider.BucketName, provider.BucketPath+post.Header.Id+"/post.json", jsonData)

	for k, v := range fileAttachments {
		err = streamFileUpload(provider.BucketName, provider.BucketPath+post.Header.Id+"/"+k, v)
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Adds a comment to the post identified by postId, and optionally nested under a parent comment
// identified by parentCommentId
func (provider *GCSProvider) CreateComment(postId string, parentCommentId string, newComment *PostComment) (*[]PostComment, error) {

	var postComments []PostComment = nil

	dat, err := downloadFileIntoMemory(provider.BucketName, provider.BucketPath+postId+"/comments.json")

	if err != nil {
		postComments = *new([]PostComment)
	} else {
		err = json.Unmarshal(dat, &postComments)

		if err != nil {
			return nil, err
		}
	}

	if parentCommentId == "" {
		postComments = append(postComments, *newComment)
	} else {
		fmt.Println("Adding comment to parent " + parentCommentId)
		res := AddCommentToParent(&postComments, parentCommentId, newComment)

		if !res {
			return nil, errors.New(fmt.Sprintf("Could not find comment '%s' to add comment to.", parentCommentId))
		}
	}

	jsonData, _ := json.Marshal(postComments)
	err = streamFileUpload(provider.BucketName, provider.BucketPath+postId+"/comments.json", jsonData)

	if err != nil {
		return nil, err
	} else {
		return &postComments, nil
	}
}

// Gets the comments for a post identified by postId
func (provider *GCSProvider) GetComments(postId string) (*[]PostComment, error) {
	var postComments []PostComment = nil

	dat, err := downloadFileIntoMemory(provider.BucketName, provider.BucketPath+postId+"/comments.json")

	if err != nil {
		postComments = *new([]PostComment)
		return &postComments, nil
	} else {
		err = json.Unmarshal(dat, &postComments)

		if err != nil {
			return nil, err
		} else {
			return &postComments, nil
		}
	}
}

// Increases the upvote count by 1 for a post identified by postId
func (provider *GCSProvider) UpvoteComment(postId string, commentId string, userEmail string) (*PostComment, error) {

	var postComments []PostComment = nil

	dat, err := downloadFileIntoMemory(provider.BucketName, provider.BucketPath+postId+"/comments.json")

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Comments for post %s not found!", postId))
	} else {
		err = json.Unmarshal(dat, &postComments)

		if err != nil {
			return nil, err
		}
	}

	upvotedComment := UpvoteComment(&postComments, commentId)

	jsonData, _ := json.Marshal(postComments)
	err = streamFileUpload(provider.BucketName, provider.BucketPath+postId+"/comments.json", jsonData)

	if err != nil {
		return nil, err
	} else {
		return upvotedComment, nil
	}
}

// Gets the file
func (provider *GCSProvider) GetFile(postId string, fileName string) ([]byte, error) {

	dat, err := downloadFileIntoMemory(provider.BucketName, provider.BucketPath+postId+"/"+fileName)

	if err != nil {
		return nil, err
	} else {
		return dat, nil
	}
}

// Deletes the post identified by postId
func (provider *GCSProvider) DeletePost(postId string) error {

	err := deleteObject(provider.BucketName, provider.BucketPath+postId)
	if err != nil {
		fmt.Printf("could not delete post %s: %s\n", postId, err)
		return err
	} else {
		return nil
	}
}

func streamFileUpload(bucketName string, name string, content []byte) error {
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
	wc := client.Bucket(bucketName).Object(name).NewWriter(ctx)
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

func downloadFileIntoMemory(bucketName string, object string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucketName).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return data, nil
}

func deleteObject(bucketName string, name string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	b := client.Bucket(bucketName)

	query := &storage.Query{Prefix: name}
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

		o := client.Bucket(bucketName).Object(attrs.Name)
		if err := o.Delete(ctx); err != nil {
			return fmt.Errorf("Object(%q).Delete: %v", name, err)
		}

		fmt.Printf("Blob %v deleted.\n", name)
	}

	return nil
}
