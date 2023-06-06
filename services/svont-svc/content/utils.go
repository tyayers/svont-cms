package content

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"tyayers/go-cms/data"
)

func InitializeProvider() data.PostIndex {

	log.Printf("Initializing Google Cloud Storage data provider.")

	dataProvider.Initialize()

	var index data.PostIndex = data.PostIndex{Index: map[string]data.PostHeader{}, IndexTime: []string{},
		IndexDrafts: map[string]int{}, IndexDeleted: map[string]int{}, IndexPopularityLikes: map[int][]string{}, IndexPopularityViews: map[int][]string{},
		IndexPopularityComments: map[int][]string{}, IndexTags: map[string]map[int]string{},
		IndexCountLikes: map[string]int{}, IndexCountComments: map[string]int{},
		IndexCountViews: map[string]int{}}

	postBytes, err := dataProvider.DownloadFile("index_headers.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.Index)
	}

	postBytes, err = dataProvider.DownloadFile("index_time.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexTime)
	}

	postBytes, err = dataProvider.DownloadFile("index_drafts.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexDrafts)
	}

	postBytes, err = dataProvider.DownloadFile("index_deleted.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexDeleted)
	}

	postBytes, err = dataProvider.DownloadFile("index_tags.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexTags)
	}

	postBytes, err = dataProvider.DownloadFile("index_popularity_likes.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexPopularityLikes)
	}

	if len(index.IndexPopularityLikes) == 0 {
		index.IndexPopularityLikes[0] = []string{}
	}

	postBytes, err = dataProvider.DownloadFile("index_popularity_comments.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexPopularityComments)
	}

	if len(index.IndexPopularityComments) == 0 {
		index.IndexPopularityComments[0] = []string{}
	}

	postBytes, err = dataProvider.DownloadFile("index_popularity_views.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexPopularityViews)
	}

	if len(index.IndexPopularityViews) == 0 {
		index.IndexPopularityViews[0] = []string{}
	}

	postBytes, err = dataProvider.DownloadFile("index_count_likes.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexCountLikes)
	}

	postBytes, err = dataProvider.DownloadFile("index_count_comments.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexCountComments)
	}

	postBytes, err = dataProvider.DownloadFile("index_count_views.json")

	if err == nil {
		json.Unmarshal(postBytes, &index.IndexCountViews)
	}

	return index
}

func FinalizeProvider(persistMode data.PersistMode, index data.PostIndex) {

	// Finalize provider
	dataProvider.Finalize(persistMode, index)

	// Persist header index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyHeaders {
		jsonData, err := json.Marshal(index.Index)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_headers.json", jsonData)
	}

	// Persist time index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyTime {
		jsonData, err := json.Marshal(index.IndexTime)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_time.json", jsonData)
	}

	// Persist drafts index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyDrafts {
		jsonData, err := json.Marshal(index.IndexDrafts)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_drafts.json", jsonData)
	}

	// Persist deleted index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyDeleted {
		jsonData, err := json.Marshal(index.IndexDeleted)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_deleted.json", jsonData)
	}

	// Persist tag index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyTags {
		jsonData, err := json.Marshal(index.IndexTags)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_tags.json", jsonData)
	}

	// Persist popularity likes index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyPopularityLikes {
		jsonData, err := json.Marshal(index.IndexPopularityLikes)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_popularity_likes.json", jsonData)
	}

	// Persist popularity comments index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyPopularityComments {
		jsonData, err := json.Marshal(index.IndexPopularityComments)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_popularity_comments.json", jsonData)
	}

	// Persist popularity views index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyPopularityViews {
		jsonData, err := json.Marshal(index.IndexPopularityViews)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_popularity_views.json", jsonData)
	}

	// Persist count likes index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyCountLikes {
		jsonData, err := json.Marshal(index.IndexCountLikes)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_count_likes.json", jsonData)
	}

	// Persist count comments index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyCountComments {
		jsonData, err := json.Marshal(index.IndexCountComments)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_count_comments.json", jsonData)
	}

	// Persist popularity likes index
	if persistMode == data.PersistAll || persistMode == data.PersistOnlyCountViews {
		jsonData, err := json.Marshal(index.IndexCountViews)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		dataProvider.UploadFile("index_count_views.json", jsonData)
	}
}

func GetPostFromProvider(postId string, draft bool) *data.Post {
	var post data.Post

	if draft {
		dat, err := dataProvider.DownloadFile("data/" + postId + "/post_draft.json")

		if err != nil {
			dat, err := dataProvider.DownloadFile("data/" + postId + "/post.json")

			if err == nil {
				json.Unmarshal(dat, &post)
			}
		} else {
			json.Unmarshal(dat, &post)
		}
	} else {
		dat, err := dataProvider.DownloadFile("data/" + postId + "/post.json")

		if err == nil {
			json.Unmarshal(dat, &post)
		}
	}

	
	

	return &post
}

func CreatePostForProvider(newPost data.Post, fileAttachments map[string][]byte) error {

	err := dataProvider.CreateDir("data/" + newPost.Header.Id)
	if err != nil {
		return err
	}

	jsonData, _ := json.Marshal(newPost)

	if newPost.Header.Draft {
		err = dataProvider.UploadFile("data/"+newPost.Header.Id+"/post_draft.json", jsonData)
	} else {
		err = dataProvider.UploadFile("data/"+newPost.Header.Id+"/post.json", jsonData)
	}

	for k, v := range fileAttachments {
		err = dataProvider.UploadFile("data/"+newPost.Header.Id+"/"+k, v)
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

func UpdatePostForProvider(post data.Post, fileAttachments map[string][]byte, draft bool) error {

	jsonData, err := json.Marshal(post)
	
	if draft {
		err = dataProvider.UploadFile("data/"+post.Header.Id+"/post_draft.json", jsonData)
	} else {
		err = dataProvider.UploadFile("data/"+post.Header.Id+"/post.json", jsonData)
	}

	for k, v := range fileAttachments {
		err = dataProvider.UploadFile("data/"+post.Header.Id+"/"+k, v)
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Adds a comment to the post identified by postId, and optionally nested under a parent comment
// identified by parentCommentId
func CreateComment(postId string, parentCommentId string, newComment *data.PostComment) (*[]data.PostComment, error) {

	var postComments []data.PostComment = nil

	dat, err := dataProvider.DownloadFile("data/" + postId + "/comments.json")

	if err != nil {
		postComments = *new([]data.PostComment)
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
	err = dataProvider.UploadFile("data/"+postId+"/comments.json", jsonData)

	if err != nil {
		return nil, err
	} else {
		return &postComments, nil
	}
}

// Gets the comments for a post identified by postId
func GetCommentsFromProvider(postId string) (*[]data.PostComment, error) {
	var postComments []data.PostComment = nil

	dat, err := dataProvider.DownloadFile("data/" + postId + "/comments.json")

	if err != nil {
		postComments = *new([]data.PostComment)
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
func UpvoteCommentForProvider(postId string, commentId string, userEmail string) (*data.PostComment, error) {

	var postComments []data.PostComment = nil

	dat, err := dataProvider.DownloadFile("data/" + postId + "/comments.json")

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Comments for post %s not found!", postId))
	} else {
		err = json.Unmarshal(dat, &postComments)

		if err != nil {
			return nil, err
		}
	}

	upvotedComment := DoUpvoteComment(&postComments, commentId)

	jsonData, _ := json.Marshal(postComments)
	err = dataProvider.UploadFile("data/"+postId+"/comments.json", jsonData)

	if err != nil {
		return nil, err
	} else {
		return upvotedComment, nil
	}
}

// Gets the file
func GetFile(postId string, fileName string) ([]byte, error) {

	dat, err := dataProvider.DownloadFile("data/" + postId + "/" + fileName)

	if err != nil {
		return nil, err
	} else {
		return dat, nil
	}
}

// Deletes the post identified by postId
func DeletePostForProvider(postId string) error {

	err := dataProvider.DeleteFile("data/" + postId)
	if err != nil {
		fmt.Printf("could not delete post %s: %s\n", postId, err)
		return err
	} else {
		return nil
	}
}

// Adds the comment to the correct parent
func AddCommentToParent(comments *[]data.PostComment, parentCommentId string, newComment *data.PostComment) bool {
	for i := 0; i < len(*comments); i++ {
		if (*comments)[i].Id == parentCommentId {
			// We found it!
			(*comments)[i].Children = append((*comments)[i].Children, *newComment)

			return true
		} else if len((*comments)[i].Children) > 0 {
			result := AddCommentToParent(&(*comments)[i].Children, parentCommentId, newComment)
			if result {
				return result
			}
		}
	}

	return false
}

// Adds the comment to the correct parent
func DoUpvoteComment(comments *[]data.PostComment, commentId string) *data.PostComment {
	for i := 0; i < len(*comments); i++ {
		if (*comments)[i].Id == commentId {
			// We found it!
			(*comments)[i].Upvotes++

			return &(*comments)[i]
		} else if len((*comments)[i].Children) > 0 {
			result := DoUpvoteComment(&(*comments)[i].Children, commentId)
			if result != nil {
				return result
			}
		}
	}

	return nil
}
