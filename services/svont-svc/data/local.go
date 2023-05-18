package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

type LocalProvider struct {
	RootDirectory string
}

// Initialize loads persisted data structures from storage, if available
func (provider *LocalProvider) Initialize() PostIndex {

	log.Printf("Initializing Google Cloud Storage data provider.")

	provider.RootDirectory = os.Getenv("ROOT_DIR")
	if provider.RootDirectory == "" {
		provider.RootDirectory = "./localdata/"
	}

	var index PostIndex = PostIndex{Index: map[string]PostHeader{}, IndexTime: []string{},
		IndexDrafts: map[string]int{}, IndexDeleted: map[string]int{}, IndexPopularityLikes: map[int][]string{}, IndexPopularityViews: map[int][]string{},
		IndexPopularityComments: map[int][]string{}, IndexTags: map[string]map[int]string{},
		IndexCountLikes: map[string]int{}, IndexCountComments: map[string]int{},
		IndexCountViews: map[string]int{}}

	postBytes, err := provider.DownloadFile("index_headers.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.Index)
	}

	postBytes, err = provider.DownloadFile("index_time.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexTime)
	}

	postBytes, err = provider.DownloadFile("index_drafts.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexDrafts)
	}

	postBytes, err = provider.DownloadFile("index_deleted.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexDeleted)
	}

	postBytes, err = provider.DownloadFile("index_tags.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexTags)
	}

	postBytes, err = provider.DownloadFile("index_popularity_likes.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexPopularityLikes)
	}

	if len(index.IndexPopularityLikes) == 0 {
		index.IndexPopularityLikes[0] = []string{}
	}

	postBytes, err = provider.DownloadFile("index_popularity_comments.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexPopularityComments)
	}

	if len(index.IndexPopularityComments) == 0 {
		index.IndexPopularityComments[0] = []string{}
	}

	postBytes, err = provider.DownloadFile("index_popularity_views.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexPopularityViews)
	}

	if len(index.IndexPopularityViews) == 0 {
		index.IndexPopularityViews[0] = []string{}
	}

	postBytes, err = provider.DownloadFile("index_count_likes.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexCountLikes)
	}

	postBytes, err = provider.DownloadFile("index_count_comments.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexCountComments)
	}

	postBytes, err = provider.DownloadFile("index_count_views.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexCountViews)
	}

	return index
}

// Finalize writes the data structures to storage
func (provider *LocalProvider) Finalize(persistMode PersistMode, index PostIndex) {

	// Persist header index
	if persistMode == PersistAll || persistMode == PersistOnlyHeaders {
		jsonData, err := json.Marshal(index.Index)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_headers.json", jsonData)
	}

	// Persist time index
	if persistMode == PersistAll || persistMode == PersistOnlyTime {
		jsonData, err := json.Marshal(index.IndexTime)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_time.json", jsonData)
	}

	// Persist drafts index
	if persistMode == PersistAll || persistMode == PersistOnlyDrafts {
		jsonData, err := json.Marshal(index.IndexDrafts)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_drafts.json", jsonData)
	}

	// Persist deleted index
	if persistMode == PersistAll || persistMode == PersistOnlyDeleted {
		jsonData, err := json.Marshal(index.IndexDeleted)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_deleted.json", jsonData)
	}

	// Persist tag index
	if persistMode == PersistAll || persistMode == PersistOnlyTags {
		jsonData, err := json.Marshal(index.IndexTags)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_tags.json", jsonData)
	}

	// Persist popularity likes index
	if persistMode == PersistAll || persistMode == PersistOnlyPopularityLikes {
		jsonData, err := json.Marshal(index.IndexPopularityLikes)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_popularity_likes.json", jsonData)
	}

	// Persist popularity comments index
	if persistMode == PersistAll || persistMode == PersistOnlyPopularityComments {
		jsonData, err := json.Marshal(index.IndexPopularityComments)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_popularity_comments.json", jsonData)
	}

	// Persist popularity views index
	if persistMode == PersistAll || persistMode == PersistOnlyPopularityViews {
		jsonData, err := json.Marshal(index.IndexPopularityViews)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_popularity_views.json", jsonData)
	}

	// Persist count likes index
	if persistMode == PersistAll || persistMode == PersistOnlyCountLikes {
		jsonData, err := json.Marshal(index.IndexCountLikes)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_count_likes.json", jsonData)
	}

	// Persist count comments index
	if persistMode == PersistAll || persistMode == PersistOnlyCountComments {
		jsonData, err := json.Marshal(index.IndexCountComments)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_count_comments.json", jsonData)
	}

	// Persist popularity likes index
	if persistMode == PersistAll || persistMode == PersistOnlyCountViews {
		jsonData, err := json.Marshal(index.IndexCountViews)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		provider.UploadFile("index_count_views.json", jsonData)
	}
}

// Returns the post specified by postId.
func (provider *LocalProvider) GetPost(postId string) *Post {

	dat, _ := provider.DownloadFile("data/" + postId + "/post.json")

	var post Post
	json.Unmarshal(dat, &post)

	return &post
}

// Creates a new post.
func (provider *LocalProvider) CreatePost(newPost Post, fileAttachments map[string][]byte) error {

	jsonData, _ := json.Marshal(newPost)
	err := provider.UploadFile("data/"+newPost.Header.Id+"/post.json", jsonData)

	for k, v := range fileAttachments {
		err = provider.UploadFile("data/"+newPost.Header.Id+"/"+k, v)
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Updates a post
func (provider *LocalProvider) UpdatePost(post Post, fileAttachments map[string][]byte) error {

	jsonData, _ := json.Marshal(post)
	err := provider.UploadFile("data/"+post.Header.Id+"/post.json", jsonData)

	for k, v := range fileAttachments {
		err = provider.UploadFile("data/"+post.Header.Id+"/"+k, v)
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Adds a comment to the post identified by postId, and optionally nested under a parent comment
// identified by parentCommentId
func (provider *LocalProvider) CreateComment(postId string, parentCommentId string, newComment *PostComment) (*[]PostComment, error) {

	var postComments []PostComment = nil

	dat, err := provider.DownloadFile("data/" + postId + "/comments.json")

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
	err = provider.UploadFile("data/"+postId+"/comments.json", jsonData)

	if err != nil {
		return nil, err
	} else {
		return &postComments, nil
	}
}

// Gets the comments for a post identified by postId
func (provider *LocalProvider) GetComments(postId string) (*[]PostComment, error) {
	var postComments []PostComment = nil

	dat, err := provider.DownloadFile("data/" + postId + "/comments.json")

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

	dat, err := provider.DownloadFile("data/" + postId + "/comments.json")

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
	err = provider.UploadFile("data/"+postId+"/comments.json", jsonData)

	if err != nil {
		return nil, err
	} else {
		return upvotedComment, nil
	}
}

// Gets the file
func (provider *LocalProvider) GetFile(postId string, fileName string) ([]byte, error) {

	dat, err := provider.DownloadFile("data/" + postId + "/" + fileName)

	if err != nil {
		return nil, err
	} else {
		return dat, nil
	}
}

// Deletes the post identified by postId
func (provider *LocalProvider) DeletePost(postId string) error {

	err := provider.DeleteFile("data/" + postId)
	if err != nil {
		fmt.Printf("could not delete post %s: %s\n", postId, err)
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
