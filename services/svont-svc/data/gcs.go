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
	Index           map[string]PostOverview
	IndexWriteMutex sync.Mutex
}

func (provider *GCSProvider) Initialize() {

	log.Printf("Initializing Google Cloud Storage data provider.")

	postBytes, err := downloadFileIntoMemory("posts/index.json")

	if err == nil {
		json.Unmarshal(postBytes, &provider.Index)
	}

	if provider.Index == nil {
		provider.Index = map[string]PostOverview{}
	}

	provider.IndexWriteMutex = sync.Mutex{}
}

func (provider *GCSProvider) Finalize() {
	//_, err := os.Create("./localdata/index.json")

	jsonData, _ := json.Marshal(provider.Index)
	err := os.WriteFile("./localdata/index.json", jsonData, 0644)

	if err != nil {
		fmt.Printf("Could not write index: %s", err)
	} else {
		fmt.Printf("Successfully wrote index.")
	}
}

// Returns the map index of post overviews
func (provider *GCSProvider) GetIndex() map[string]PostOverview {
	return provider.Index
}

// Returns paginated posts array
func (provider *GCSProvider) GetPosts(start int, limit int) []PostOverview {
	return []PostOverview{}
}

// Returns the post specified by postId.
func (provider *GCSProvider) GetPost(postId string) Post {

	dat, _ := os.ReadFile("./localdata/" + postId + "/post.json")

	var post Post
	json.Unmarshal(dat, &post)

	post.Header = provider.Index[postId]

	return post
}

// Creates a new post.
func (provider *GCSProvider) CreatePost(newPost Post, fileAttachments map[string][]byte) error {
	if err := os.Mkdir("./localdata/"+newPost.Header.Id, os.ModePerm); err != nil {
		return err
	}

	provider.IndexWriteMutex.Lock()
	provider.Index[newPost.Header.Id] = newPost.Header
	provider.IndexWriteMutex.Unlock()

	jsonData, _ := json.Marshal(newPost)
	err := os.WriteFile("./localdata/"+newPost.Header.Id+"/post.json", jsonData, 0644)

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Updates a post
func (provider *GCSProvider) UpdatePost(post Post) {

}

// Increases the upvote count by 1 for a post identified by postId
func (provider *GCSProvider) UpvotePost(postId string, userEmail string) (*PostOverview, error) {

	post, ok := provider.Index[postId]

	if ok {
		provider.IndexWriteMutex.Lock()
		post.Upvotes++
		provider.IndexWriteMutex.Unlock()

		provider.Index[postId] = post
		return &post, nil
	} else {
		return nil, errors.New("Post not found")
	}
}

// Adds a comment to the post identified by postId, and optionally nested under a parent comment
// identified by parentCommentId
func (provider *GCSProvider) AddComment(postId string, parentCommentId string, content string) (*[]PostComment, error) {

	var postComments []PostComment = nil
	var newComment *PostComment = nil

	post, ok := provider.Index[postId]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Post %s not found!", postId))
	}

	newComment = new(PostComment)
	newComment.Id = time.Now().Format("20060102_150405.99_") + RandomString(12)
	newComment.Children = []PostComment{}
	newComment.Content = content

	dat, err := os.ReadFile("./localdata/" + postId + "/comments.json")

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
			return nil, errors.New(fmt.Sprintf("Could not find post '%s' to add comment to.", parentCommentId))
		}
	}

	jsonData, _ := json.Marshal(postComments)
	err = os.WriteFile("./localdata/"+postId+"/comments.json", jsonData, 0644)

	if err != nil {
		return nil, err
	} else {
		provider.IndexWriteMutex.Lock()
		post = provider.Index[postId]
		post.CommentCount++
		provider.Index[postId] = post
		provider.IndexWriteMutex.Unlock()

		return &postComments, nil
	}
}

// Gets the comments for a post identified by postId
func (provider *GCSProvider) GetComments(postId string) (*[]PostComment, error) {
	var postComments []PostComment = nil

	dat, err := os.ReadFile("./localdata/" + postId + "/comments.json")

	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(dat, &postComments)

		if err != nil {
			return nil, err
		} else {
			return &postComments, nil
		}
	}
}

// Deletes the post identified by postId
func (provider *GCSProvider) DeletePost(postId string) error {

	provider.IndexWriteMutex.Lock()
	delete(provider.Index, postId)
	provider.IndexWriteMutex.Unlock()

	err := os.RemoveAll("./localdata/" + postId)

	if err != nil {
		fmt.Printf("Error deleting post: %s", err)
		return err
	} else {
		return nil
	}
}

func streamFileUpload(name string, content []byte) error {
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
	wc := client.Bucket("cms_tg6qp4dq8").Object(name).NewWriter(ctx)
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

func downloadFileIntoMemory(object string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket("cms_tg6qp4dq8").Object(object).NewReader(ctx)
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

func deleteObject(name string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	b := client.Bucket("cms_tg6qp4dq8")

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

		o := client.Bucket("cms_tg6qp4dq8").Object(attrs.Name)
		if err := o.Delete(ctx); err != nil {
			return fmt.Errorf("Object(%q).Delete: %v", name, err)
		}

		fmt.Printf("Blob %v deleted.\n", name)
	}

	// blobs := b.Objects(name)
	// o := client.Bucket("cms_tg6qp4dq8").Object(name)

	// if err := o.Delete(ctx); err != nil {
	// 	return fmt.Errorf("Object(%q).Delete: %v", name, err)
	// }
	// fmt.Fprintf(w, "Blob %v deleted.\n", name)

	return nil
}
