package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

type LocalProvider struct {
	Index           map[string]PostOverview
	IndexWriteMutex sync.Mutex
}

func (provider *LocalProvider) Initialize() {

	log.Printf("Initializing Local File data provider.")

	os.Mkdir("./localdata/", os.ModePerm)

	dat, _ := os.ReadFile("./localdata/index.json")
	json.Unmarshal(dat, &provider.Index)

	if provider.Index == nil {
		provider.Index = map[string]PostOverview{}
	}

	provider.IndexWriteMutex = sync.Mutex{}
}

func (provider *LocalProvider) Finalize() {
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
func (provider *LocalProvider) GetIndex() map[string]PostOverview {
	return provider.Index
}

// Returns paginated posts array
func (provider *LocalProvider) GetPosts(start int, limit int) []PostOverview {
	return []PostOverview{}
}

// Returns the post specified by postId.
func (provider *LocalProvider) GetPost(postId string) *Post {

	dat, _ := os.ReadFile("./localdata/" + postId + "/post.json")

	var post Post
	json.Unmarshal(dat, &post)

	post.Header = provider.Index[postId]

	return &post
}

// Returns the post overview
func (provider *LocalProvider) GetPostOverview(postId string) *PostOverview {

	post, ok := provider.Index[postId]

	if ok {
		return &post
	} else {
		return nil
	}
}

// Creates a new post.
func (provider *LocalProvider) CreatePost(newPost Post, fileAttachments map[string][]byte) error {
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
func (provider *LocalProvider) UpdatePost(post Post, fileAttachments map[string][]byte) error {

	provider.IndexWriteMutex.Lock()
	header := provider.Index[post.Header.Id]
	header.Title = post.Header.Title
	header.Summary = post.Header.Summary
	header.Updated = post.Header.Updated
	provider.Index[post.Header.Id] = header
	provider.IndexWriteMutex.Unlock()

	post.Header = header

	jsonData, _ := json.Marshal(post)
	err := os.WriteFile("./localdata/"+post.Header.Id+"/post.json", jsonData, 0644)

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Increases the upvote count by 1 for a post identified by postId
func (provider *LocalProvider) UpvotePost(postId string, userEmail string) (*PostOverview, error) {

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
func (provider *LocalProvider) CreateComment(postId string, parentCommentId string, newComment *PostComment) (*[]PostComment, error) {

	var postComments []PostComment = nil

	post, ok := provider.Index[postId]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Post %s not found!", postId))
	}

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
			return nil, errors.New(fmt.Sprintf("Could not find comment '%s' to add comment to.", parentCommentId))
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
func (provider *LocalProvider) GetComments(postId string) (*[]PostComment, error) {
	var postComments []PostComment = nil

	dat, err := os.ReadFile("./localdata/" + postId + "/comments.json")

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
func (provider *LocalProvider) UpvoteComment(postId string, commentId string, userEmail string) (*PostComment, error) {

	var postComments []PostComment = nil

	dat, err := os.ReadFile("./localdata/" + postId + "/comments.json")

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
	err = os.WriteFile("./localdata/"+postId+"/comments.json", jsonData, 0644)

	if err != nil {
		return nil, err
	} else {
		return upvotedComment, nil
	}
}

// Deletes the post identified by postId
func (provider *LocalProvider) DeletePost(postId string) error {

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
