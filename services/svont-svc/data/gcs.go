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
	"sort"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GCSProvider struct {
	Index           map[string]PostOverview
	IndexPopularity map[int][]string
	IndexWriteMutex sync.Mutex
}

// Initialize loads persisted data structures from storage, if available
func (provider *GCSProvider) Initialize() {

	log.Printf("Initializing Google Cloud Storage data provider.")

	postBytes, err := downloadFileIntoMemory("posts/index.json")

	if err == nil {
		json.Unmarshal(postBytes, &provider.Index)
	}

	if provider.Index == nil {
		provider.Index = map[string]PostOverview{}
	}

	postBytes, err = downloadFileIntoMemory("posts/index_popularity.json")

	if err == nil {
		json.Unmarshal(postBytes, &provider.IndexPopularity)
	}

	if provider.IndexPopularity == nil {
		provider.IndexPopularity = map[int][]string{}
		provider.IndexPopularity[0] = []string{}
	}

	provider.IndexWriteMutex = sync.Mutex{}
}

// Finalize writes the data structures to storage
func (provider *GCSProvider) Finalize() {

	jsonData, err := json.Marshal(provider.Index)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	streamFileUpload("posts/index.json", jsonData)

	if err != nil {
		fmt.Printf("Could not write index: %s", err)
	} else {
		fmt.Printf("Successfully wrote index.")
	}

	jsonData, err = json.Marshal(provider.IndexPopularity)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	streamFileUpload("posts/index_popularity.json", jsonData)

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

// Returns paginated list of most popular posts
func (provider *GCSProvider) GetPopularPosts(start int, limit int) []PostOverview {

	postsByPopularity := []PostOverview{}

	keys := make([]int, 0, len(provider.IndexPopularity))
	for k := range provider.IndexPopularity {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, v := range keys {
		for i := range provider.IndexPopularity[v] {
			postsByPopularity = append(postsByPopularity, provider.Index[provider.IndexPopularity[v][i]])

			if len(postsByPopularity) >= limit {
				break
			}
		}

		if len(postsByPopularity) >= limit {
			break
		}
	}

	return postsByPopularity
}

// Returns the post specified by postId.
func (provider *GCSProvider) GetPost(postId string) *Post {

	dat, _ := downloadFileIntoMemory("posts/" + postId + "/post.json")

	var post Post
	json.Unmarshal(dat, &post)

	post.Header = provider.Index[postId]

	return &post
}

// Returns the post overview
func (provider *GCSProvider) GetPostOverview(postId string) *PostOverview {

	post, ok := provider.Index[postId]

	if ok {
		return &post
	} else {
		return nil
	}
}

// Creates a new post.
func (provider *GCSProvider) CreatePost(newPost Post, fileAttachments map[string][]byte) error {

	provider.IndexWriteMutex.Lock()
	provider.Index[newPost.Header.Id] = newPost.Header
	provider.IndexPopularity[0] = append(provider.IndexPopularity[0], newPost.Header.Id)
	provider.IndexWriteMutex.Unlock()

	jsonData, _ := json.Marshal(newPost)
	err := streamFileUpload("posts/"+newPost.Header.Id+"/post.json", jsonData)

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Updates a post
func (provider *GCSProvider) UpdatePost(post Post, fileAttachments map[string][]byte) error {

	provider.IndexWriteMutex.Lock()
	header := provider.Index[post.Header.Id]
	header.Title = post.Header.Title
	header.Summary = post.Header.Summary
	header.Updated = post.Header.Updated
	provider.Index[post.Header.Id] = header
	provider.IndexWriteMutex.Unlock()

	post.Header = header

	jsonData, _ := json.Marshal(post)
	err := streamFileUpload("posts/"+post.Header.Id+"/post.json", jsonData)

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Increases the upvote count by 1 for a post identified by postId
func (provider *GCSProvider) UpvotePost(postId string, userEmail string) (*PostOverview, error) {

	post, ok := provider.Index[postId]

	if ok {
		provider.IndexWriteMutex.Lock()
		post.Upvotes++
		provider.Index[postId] = post

		// Remove item from previous popularity space
		for i, s := range provider.IndexPopularity[post.Upvotes-1] {
			if s == post.Id {
				// We found our post in the old spot, now remove
				provider.IndexPopularity[post.Upvotes-1][i] = provider.IndexPopularity[post.Upvotes-1][len(provider.IndexPopularity[post.Upvotes-1])-1] // Copy last element to index i.
				provider.IndexPopularity[post.Upvotes-1][len(provider.IndexPopularity[post.Upvotes-1])-1] = ""                                          // Erase last element (write zero value).
				provider.IndexPopularity[post.Upvotes-1] = provider.IndexPopularity[post.Upvotes-1][:len(provider.IndexPopularity[post.Upvotes-1])-1]   // Truncate slice.
			}
		}

		// Add to new popularity spot
		val, ok := provider.IndexPopularity[post.Upvotes]
		// If the key exists
		if ok {
			val = append(val, post.Id)
			provider.IndexPopularity[post.Upvotes] = val
		} else {
			provider.IndexPopularity[post.Upvotes] = []string{post.Id}
		}

		provider.IndexWriteMutex.Unlock()

		return &post, nil
	} else {
		return nil, errors.New("Post not found")
	}
}

// Adds a comment to the post identified by postId, and optionally nested under a parent comment
// identified by parentCommentId
func (provider *GCSProvider) CreateComment(postId string, parentCommentId string, newComment *PostComment) (*[]PostComment, error) {

	var postComments []PostComment = nil

	post, ok := provider.Index[postId]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Post %s not found!", postId))
	}

	dat, err := downloadFileIntoMemory("posts/" + postId + "/comments.json")

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
	err = streamFileUpload("posts/"+postId+"/comments.json", jsonData)

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

	dat, err := downloadFileIntoMemory("posts/" + postId + "/comments.json")

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

	dat, err := downloadFileIntoMemory("posts/" + postId + "/comments.json")

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
	err = streamFileUpload("posts/"+postId+"/comments.json", jsonData)

	if err != nil {
		return nil, err
	} else {
		return upvotedComment, nil
	}
}

// Deletes the post identified by postId
func (provider *GCSProvider) DeletePost(postId string) error {

	provider.IndexWriteMutex.Lock()

	// Remove item from popularity space
	for i, s := range provider.IndexPopularity[provider.Index[postId].Upvotes] {
		if s == postId {
			// We found our post in the old spot, now remove
			provider.IndexPopularity[provider.Index[postId].Upvotes][i] = provider.IndexPopularity[provider.Index[postId].Upvotes][len(provider.IndexPopularity[provider.Index[postId].Upvotes])-1] // Copy last element to index i.
			provider.IndexPopularity[provider.Index[postId].Upvotes][len(provider.IndexPopularity[provider.Index[postId].Upvotes])-1] = ""                                                          // Erase last element (write zero value).
			provider.IndexPopularity[provider.Index[postId].Upvotes] = provider.IndexPopularity[provider.Index[postId].Upvotes][:len(provider.IndexPopularity[provider.Index[postId].Upvotes])-1]   // Truncate slice.
		}
	}

	delete(provider.Index, postId)

	provider.IndexWriteMutex.Unlock()

	err := deleteObject("posts/" + postId)
	if err != nil {
		fmt.Printf("could not delete post %s: %s\n", postId, err)
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
