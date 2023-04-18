package content

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"os"
	"sync"
	"time"

	"tyayers/go-cms/data"

	"github.com/blevesearch/bleve/v2"
)

var postsMutex = &sync.Mutex{}
var postsMap = map[string]data.PostOverview{}

var searchIndex bleve.Index

// Local data provider, uncomment to test locally with files (in localdata dir)
var dataProvider data.Provider = &data.LocalProvider{}

// Google Cloud Storage provider, uncomment to use GCS as storage provider
// var dataProvider data.Provider = &data.GCSProvider{}

func Initialize() {
	dataProvider.Initialize()

	// Initialize bleve search index, if it doesn't exist
	if _, err := os.Stat("./posts.bleve"); os.IsNotExist(err) {
		// Initialize bleve
		mapping := bleve.NewIndexMapping()
		var err error
		searchIndex, err = bleve.New("posts.bleve", mapping)
		if err != nil {
			fmt.Println(err)
		}

		for k, v := range dataProvider.GetIndex() {
			fmt.Printf("indexing key[%s] value[%s]\n", k, v.Id)
			searchIndex.Index(v.Id, v)
		}
	} else {
		searchIndex, err = bleve.Open("posts.bleve")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Finalize() {
	dataProvider.Finalize()
}

func GetAllPosts(start int, limit int) map[string]data.PostOverview {
	return dataProvider.GetIndex()
}

func GetPopularPosts(start int, limit int) []data.PostOverview {
	return dataProvider.GetPopularPosts(start, limit)
}

func GetPost(postId string) *data.Post {
	return dataProvider.GetPost(postId)
}

func GetPostOverview(postId string) *data.PostOverview {
	return dataProvider.GetPostOverview(postId)
}

func CreatePost(newPost *data.Post, attachments []multipart.FileHeader) error {

	newPost.Header.Id = time.Now().Format("20060102_") + RandomString(12)
	// newPost.Header.Id = time.Now().Format("20060102_150405.99_") + strings.Replace(strings.ToLower(newPost.Title), " ", "_", -1)
	// newPost.Header.Id = time.Now().Format("20060102_") + strings.Replace(strings.ToLower(newPost.Header.Title), " ", "_", -1)
	newPost.Files = []string{}
	newPost.Header.Created = time.Now().Format("2006-01-02T15:04:05-0700")

	files := map[string][]byte{}
	for _, attachment := range attachments {
		src, _ := attachment.Open()
		defer src.Close()

		bytes, _ := ioutil.ReadAll(src)
		files[attachment.Filename] = bytes
		// streamFileUpload("posts/"+newPost.Header.Id+"/"+attachment.Filename, bytes)

		newPost.Files = append(newPost.Files, attachment.Filename)
	}

	newPost.Header.FileCount = len(newPost.Files)

	err := dataProvider.CreatePost(*newPost, files)

	if err != nil {
		return err
	} else {

		// Create successful, now index for search as well
		searchIndex.Index(newPost.Header.Id, newPost.Header)

		return nil
	}
}

func UpdatePost(updatedPost *data.Post, attachments []multipart.FileHeader) error {

	updatedPost.Header.Updated = time.Now().Format("2006-01-02T15:04:05-0700")

	files := map[string][]byte{}
	for _, attachment := range attachments {
		src, _ := attachment.Open()
		defer src.Close()

		bytes, _ := ioutil.ReadAll(src)
		files[attachment.Filename] = bytes
		// streamFileUpload("posts/"+newPost.Header.Id+"/"+attachment.Filename, bytes)

		updatedPost.Files = append(updatedPost.Files, attachment.Filename)
	}

	updatedPost.Header.FileCount = len(updatedPost.Files)

	err := dataProvider.UpdatePost(*updatedPost, files)

	if err != nil {
		return err
	} else {

		// Create successful, now index for search as well
		searchIndex.Index(updatedPost.Header.Id, updatedPost.Header)

		return nil
	}

	// if entry, ok := postsMap[updatedPost.Header.Id]; ok {
	// 	entry.Title = updatedPost.Header.Title
	// 	postsMutex.Lock()
	// 	postsMap[updatedPost.Header.Id] = entry
	// 	postsMutex.Unlock()
	// }

	// if writeToStorage {
	// 	jsonData, err := json.Marshal(updatedPost)
	// 	if err != nil {
	// 		fmt.Printf("could not marshal json: %s\n", err)
	// 		return
	// 	}

	// 	streamFileUpload("posts/"+updatedPost.Header.Id+"/post.json", jsonData)
	// }
}

func UpvotePost(postId string, userEmail string) (*data.PostOverview, error) {
	return dataProvider.UpvotePost(postId, userEmail)
}

func AddCommentToPost(postId string, parentCommentId string, authorId string, authorDisplayName string, authorProfilePic string, content string) (*[]data.PostComment, error) {

	var newComment = new(data.PostComment)
	newComment.Id = time.Now().Format("20060102_150405.99_") + RandomString(12)
	newComment.Created = time.Now().Format("2006-01-02T15:04:05-0700")
	newComment.AuthorId = authorId
	newComment.AuthorDisplayName = authorDisplayName
	newComment.AuthorProfilePic = authorProfilePic
	newComment.Children = []data.PostComment{}
	newComment.Content = content

	return dataProvider.CreateComment(postId, parentCommentId, newComment)
}

// func AddCommentToPost(postId string, parentCommentId string, newComment *data.PostComment) (error) {
// 	newComment.Id = time.Now().Format("20060102_150405.99_") + RandomString(12)
// 	return dataProvider.CreateComment(postId, parentCommentId, newComment)
// }

// Gets all of the comments for a post
func GetComments(postId string) (*[]data.PostComment, error) {
	return dataProvider.GetComments(postId)
}

// Upvotes a specific comment
func UpvoteComment(postId string, commentId string, userEmail string) (*data.PostComment, error) {
	return dataProvider.UpvoteComment(postId, commentId, userEmail)
}

// Attaches a file to a post
func AttachFileToPost(postId string, fileName string, file []byte) {

	// postBytes, _ := downloadFileIntoMemory("posts/" + postId + "/post.json")
	// var updatedPost data.Post
	// json.Unmarshal(postBytes, &updatedPost)

	// streamFileUpload("posts/"+postId+"/"+fileName, file)

	// updatedPost.Files = append(updatedPost.Files, fileName)
	// jsonData, _ := json.Marshal(updatedPost)
	// streamFileUpload("posts/"+updatedPost.Header.Id+"/post.json", jsonData)
}

// Deletes a post
func DeletePost(postId string) error {
	// err := deleteObject("posts/" + postId)
	// if err != nil {
	// 	fmt.Printf("could not delete post %s: %s\n", postId, err)
	// 	return err
	// }

	// postsMutex.Lock()
	// delete(postsMap, postId)
	// postsMutex.Unlock()

	// return nil
	return dataProvider.DeletePost(postId)
}

// Searches posts
func SearchPosts(text string) ([]data.PostOverview, error) {
	//query := bleve.NewMatchQuery(text)
	query := bleve.NewFuzzyQuery(text)
	search := bleve.NewSearchRequest(query)

	if searchIndex != nil {
		searchResults, err := searchIndex.Search(search)
		if err != nil {
			fmt.Println(err)
			return nil, err
		} else {
			results := []data.PostOverview{}
			fmt.Println(searchResults)
			dataMap := dataProvider.GetIndex()
			for _, val := range searchResults.Hits {
				results = append(results, dataMap[val.ID])
			}

			return results, nil
		}
	} else {
		return nil, errors.New("search index is nil!")
	}

}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}
