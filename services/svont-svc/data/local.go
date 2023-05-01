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
	IndexWriteMutex sync.Mutex
}

func (provider *LocalProvider) Initialize() (map[string]PostOverview, []string, map[int][]string, map[string]map[int]string) {

	log.Printf("Initializing Local File data provider.")

	os.Mkdir("./localdata/", os.ModePerm)

	var index_main map[string]PostOverview
	var index_time []string
	var index_popularity map[int][]string
	var index_tags map[string]map[int]string

	dat, _ := os.ReadFile("./localdata/index.json")
	json.Unmarshal(dat, &index_main)

	if index_main == nil {
		index_main = map[string]PostOverview{}
	}

	dat, _ = os.ReadFile("./localdata/index_time.json")
	json.Unmarshal(dat, &index_time)

	if index_time == nil {
		index_time = []string{}
	}

	dat, _ = os.ReadFile("./localdata/index_popularity.json")
	json.Unmarshal(dat, &index_popularity)

	if index_popularity == nil {
		index_popularity = map[int][]string{}
		// initialize 0 popularity slot
		index_popularity[0] = []string{}
	}

	dat, _ = os.ReadFile("./localdata/index_tags.json")
	json.Unmarshal(dat, &index_tags)

	if index_tags == nil {
		index_tags = map[string]map[int]string{}
	}

	provider.IndexWriteMutex = sync.Mutex{}

	return index_main, index_time, index_popularity, index_tags
}

func (provider *LocalProvider) Finalize(index_main map[string]PostOverview, index_time []string, index_populary map[int][]string, index_tags map[string]map[int]string) {
	//_, err := os.Create("./localdata/index.json")

	jsonData, _ := json.Marshal(index_main)
	err := os.WriteFile("./localdata/index.json", jsonData, 0644)

	if err != nil {
		fmt.Printf("Could not write index: %s", err)
	} else {
		fmt.Printf("Successfully wrote index.")
	}

	jsonData, _ = json.Marshal(index_time)
	err = os.WriteFile("./localdata/index_time.json", jsonData, 0644)

	if err != nil {
		fmt.Printf("Could not write time index: %s", err)
	} else {
		fmt.Printf("Successfully wrote time index.")
	}

	jsonData, _ = json.Marshal(index_populary)
	err = os.WriteFile("./localdata/index_popularity.json", jsonData, 0644)

	if err != nil {
		fmt.Printf("Could not write popularity index: %s", err)
	} else {
		fmt.Printf("Successfully wrote popularity index.")
	}

	jsonData, _ = json.Marshal(index_tags)
	err = os.WriteFile("./localdata/index_tags.json", jsonData, 0644)

	if err != nil {
		fmt.Printf("Could not write tag index: %s", err)
	} else {
		fmt.Printf("Successfully wrote tag index.")
	}
}

// Returns the post specified by postId.
func (provider *LocalProvider) GetPost(postId string) *Post {

	dat, _ := os.ReadFile("./localdata/" + postId + "/post.json")

	var post Post
	json.Unmarshal(dat, &post)

	return &post
}

// Creates a new post.
func (provider *LocalProvider) CreatePost(newPost Post, fileAttachments map[string][]byte) error {
	if err := os.Mkdir("./localdata/"+newPost.Header.Id, os.ModePerm); err != nil {
		return err
	}

	jsonData, _ := json.Marshal(newPost)
	err := os.WriteFile("./localdata/"+newPost.Header.Id+"/post.json", jsonData, 0644)

	for k, v := range fileAttachments {
		err = os.WriteFile("./localdata/"+newPost.Header.Id+"/"+k, v, 0644)
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
	err := os.WriteFile("./localdata/"+post.Header.Id+"/post.json", jsonData, 0644)

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

func (provider *LocalProvider) GetFile(postId string, fileName string) ([]byte, error) {

	dat, err := os.ReadFile("./localdata/" + postId + "/" + fileName)

	if err != nil {
		return nil, err
	} else {
		return dat, nil
	}
}

// Deletes the post identified by postId
func (provider *LocalProvider) DeletePost(postId string) error {

	err := os.RemoveAll("./localdata/" + postId)

	if err != nil {
		fmt.Printf("Error deleting post: %s", err)
		return err
	} else {
		return nil
	}
}
