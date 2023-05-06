package content

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"sort"
	"sync"
	"time"

	"tyayers/go-cms/data"

	"github.com/blevesearch/bleve/v2"
)

var postsMutex = sync.Mutex{}
var index map[string]data.PostOverview
var index_time []string
var index_popularity map[int][]string
var index_tags map[string]map[int]string

var searchIndex bleve.Index
var tagIndex bleve.Index

// Local data provider, uncomment to test locally with files (in localdata dir)
// var dataProvider data.Provider = &data.LocalProvider{}

// Google Cloud Storage provider, uncomment to use GCS as storage provider
var dataProvider data.Provider = &data.GCSProvider{}

func Initialize(force bool) {
	index, index_time, index_popularity, index_tags = dataProvider.Initialize()

	// Initialize bleve search index, if it doesn't exist
	if _, err := os.Stat("./posts.bleve"); os.IsNotExist(err) || force {
		// Initialize bleve
		os.RemoveAll("./posts.bleve")
		mapping := bleve.NewIndexMapping()
		var err error
		searchIndex, err = bleve.New("posts.bleve", mapping)
		if err != nil {
			fmt.Println(err)
		}

		for k, v := range index {
			fmt.Printf("indexing key[%s] value[%s]\n", k, v.Id)
			searchIndex.Index(v.Id, v)
		}
	} else {
		log.Printf("Loading local bleve index..")
		searchIndex, err = bleve.Open("posts.bleve")
		log.Printf("Finished loading local bleve index..")
		if err != nil {
			fmt.Println(err)
		}
	}

	// Initialize bleve tag index, if it doesn't exist
	if _, err := os.Stat("./tags.bleve"); os.IsNotExist(err) || force {
		// Initialize bleve
		os.RemoveAll("./tags.bleve")
		mapping := bleve.NewIndexMapping()
		var err error
		tagIndex, err = bleve.New("tags.bleve", mapping)
		if err != nil {
			fmt.Println(err)
		}

		for k, _ := range index_tags {
			fmt.Printf("indexing key[%s] value[%s]\n", k, k)
			tagIndex.Index(k, k)
		}
	} else {
		log.Printf("Loading local bleve tag index..")
		tagIndex, err = bleve.Open("tags.bleve")
		log.Printf("Finished loading local bleve tag index..")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Finalize() {
	dataProvider.Finalize(index, index_time, index_popularity, index_tags)
}

func GetPosts(start int, limit int) []data.PostOverview {
	resultPosts := []data.PostOverview{}

	if len(index_time) > 0 {
		postIndex := len(index_time) - 1 - start

		for postIndex >= 0 && len(resultPosts) < limit {
			resultPosts = append(resultPosts, index[index_time[postIndex]])
			postIndex--
		}
	}

	return resultPosts
}

func GetPopularPosts(start int, limit int) []data.PostOverview {
	postsByPopularity := []data.PostOverview{}

	keys := make([]int, 0, len(index_popularity))
	for k := range index_popularity {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, v := range keys {
		for i := range index_popularity[v] {
			postsByPopularity = append(postsByPopularity, index[index_popularity[v][i]])

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

func GetTaggedPosts(tagName string, start int, limit int) []data.PostOverview {
	taggedPosts := []data.PostOverview{}

	posts, ok := index_tags[tagName]

	if ok {
		keys := make([]int, 0, len(posts))
		for k := range posts {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] > keys[j]
		})

		postIndex := start
		for postIndex < len(keys) && len(taggedPosts) < limit {
			post, ok := index[posts[keys[postIndex]]]
			if ok {
				taggedPosts = append(taggedPosts, post)
			}

			postIndex++
		}
	}

	return taggedPosts
}

func GetPost(postId string) *data.Post {
	var post = dataProvider.GetPost(postId)
	post.Header = index[postId]
	return post
}

func GetPostOverview(postId string) data.PostOverview {
	return index[postId]
}

func CreatePost(newPost *data.Post, attachments []multipart.FileHeader) error {

	var createTime = time.Now()

	newPost.Header.Id = time.Now().Format("20060102_") + RandomString(12)
	// newPost.Header.Id = time.Now().Format("20060102_150405.99_") + strings.Replace(strings.ToLower(newPost.Title), " ", "_", -1)
	// newPost.Header.Id = time.Now().Format("20060102_") + strings.Replace(strings.ToLower(newPost.Header.Title), " ", "_", -1)
	newPost.Files = []string{}
	newPost.Header.Created = createTime.Format("2006-01-02T15:04:05-0700")

	files := map[string][]byte{}
	for _, attachment := range attachments {
		src, _ := attachment.Open()
		defer src.Close()

		bytes, _ := ioutil.ReadAll(src)
		files[attachment.Filename] = bytes

		newPost.Files = append(newPost.Files, attachment.Filename)
	}

	newPost.Header.FileCount = len(newPost.Files)

	postsMutex.Lock()

	// Add to time index
	index_time = append(index_time, newPost.Header.Id)
	newPost.Header.Index = len(index_time) - 1

	// Add to id index
	index[newPost.Header.Id] = newPost.Header

	// Add to popularity index
	index_popularity[0] = append(index_popularity[0], newPost.Header.Id)

	// Add to tag index
	if newPost.Header.Tags != nil {
		for _, tag := range newPost.Header.Tags {

			if tag != "" {
				_, ok := index_tags[tag]

				if !ok {
					index_tags[tag] = map[int]string{}
					tagIndex.Index(tag, tag)
				}

				index_tags[tag][newPost.Header.Index] = newPost.Header.Id
			}
		}
	}

	postsMutex.Unlock()

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

	postsMutex.Lock()
	header := index[updatedPost.Header.Id]
	header.Title = updatedPost.Header.Title
	header.Summary = updatedPost.Header.Summary

	UpdateTags(header.Id, header.Index, header.Tags, updatedPost.Header.Tags)
	header.Tags = updatedPost.Header.Tags

	header.Updated = updatedPost.Header.Updated
	index[updatedPost.Header.Id] = header
	postsMutex.Unlock()

	updatedPost.Header = header

	err := dataProvider.UpdatePost(*updatedPost, files)

	if err != nil {
		return err
	} else {

		// Create successful, now index for search as well
		searchIndex.Index(updatedPost.Header.Id, updatedPost.Header)

		return nil
	}
}

func UpvotePost(postId string, userEmail string) (*data.PostOverview, error) {

	post, ok := index[postId]

	if ok {
		postsMutex.Lock()
		post.Upvotes++
		index[postId] = post

		// Remove item from previous popularity space
		for i, s := range index_popularity[post.Upvotes-1] {
			if s == post.Id {
				// We found our post in the old spot, now remove
				index_popularity[post.Upvotes-1][i] = index_popularity[post.Upvotes-1][len(index_popularity[post.Upvotes-1])-1] // Copy last element to index i.
				index_popularity[post.Upvotes-1][len(index_popularity[post.Upvotes-1])-1] = ""                                  // Erase last element (write zero value).
				index_popularity[post.Upvotes-1] = index_popularity[post.Upvotes-1][:len(index_popularity[post.Upvotes-1])-1]   // Truncate slice.
			}
		}

		// Add to new popularity spot
		val, ok := index_popularity[post.Upvotes]
		// If the key exists
		if ok {
			index_popularity[post.Upvotes] = append(val, post.Id)
		} else {
			index_popularity[post.Upvotes] = []string{post.Id}
		}

		postsMutex.Unlock()

		return &post, nil
	} else {
		return nil, errors.New("Post not found")
	}
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

	post, ok := index[postId]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Post %s not found!", postId))
	} else {
		result, err := dataProvider.CreateComment(postId, parentCommentId, newComment)

		if err == nil {
			postsMutex.Lock()
			post.CommentCount++
			index[postId] = post
			postsMutex.Unlock()

			return result, nil
		} else {
			return nil, err
		}
	}
}

// Gets all of the comments for a post
func GetComments(postId string) (*[]data.PostComment, error) {
	return dataProvider.GetComments(postId)
}

// Upvotes a specific comment
func UpvoteComment(postId string, commentId string, userEmail string) (*data.PostComment, error) {
	return dataProvider.UpvoteComment(postId, commentId, userEmail)
}

// Gets a file attachment for a post
func GetFileForPost(postId string, fileName string) ([]byte, error) {
	return dataProvider.GetFile(postId, fileName)
}

// Deletes a post
func DeletePost(postId string) error {

	postsMutex.Lock()
	delete(index, postId)
	// TODO delete from index_time and index_popularity
	postsMutex.Unlock()

	return dataProvider.DeletePost(postId)
}

// Searches posts
func SearchPosts(text string) ([]data.PostOverview, error) {
	query := bleve.NewMatchQuery(text)
	query.Fuzziness = 2

	query2 := bleve.NewPrefixQuery(text)

	query3 := bleve.NewDisjunctionQuery(query, query2)
	search := bleve.NewSearchRequest(query3)

	if searchIndex != nil {
		searchResults, err := searchIndex.Search(search)
		if err != nil {
			fmt.Println(err)
			return nil, err
		} else {
			results := []data.PostOverview{}
			fmt.Println(searchResults)
			dataMap := index
			for _, val := range searchResults.Hits {
				results = append(results, dataMap[val.ID])
			}

			return results, nil
		}
	} else {
		return nil, errors.New("search index is nil!")
	}

}

func SearchTags(text string) ([]data.SearchResult, error) {
	query := bleve.NewMatchQuery(text)
	query.Fuzziness = 2

	query2 := bleve.NewPrefixQuery(text)

	query3 := bleve.NewDisjunctionQuery(query, query2)

	//query := bleve.NewFuzzyQuery(text)
	search := bleve.NewSearchRequest(query3)

	if tagIndex != nil {
		searchResults, err := tagIndex.Search(search)
		if err != nil {
			fmt.Println(err)
			return nil, err
		} else {
			results := []data.SearchResult{}
			fmt.Println(searchResults)
			for _, val := range searchResults.Hits {
				tag, ok := index_tags[val.ID]
				if ok {
					newTagResult := data.SearchResult{
						Id:    val.ID,
						Title: val.ID,
						Count: len(tag),
					}

					results = append(results, newTagResult)
				}
			}

			return results, nil
		}
	} else {
		return nil, errors.New("tag search index is nil!")
	}
}

func UpdateTags(postId string, postIndex int, originalTagList []string, newTagList []string) {
	tagsToRemove, tagsToAdd := GetUpdatedTags(originalTagList, newTagList)

	for _, removeTag := range tagsToRemove {
		val, ok := index_tags[removeTag]
		// If the key exists
		if ok {
			delete(val, postIndex)
		} else {
			// Log error
		}
	}

	for _, addTag := range tagsToAdd {
		if addTag != "" {
			_, ok := index_tags[addTag]

			if !ok {
				index_tags[addTag] = map[int]string{}
				tagIndex.Index(addTag, addTag)
			}

			index_tags[addTag][postIndex] = postId
		}
	}
}

func GetUpdatedTags(originalTagList []string, newTagList []string) ([]string, []string) {
	tagsToRemove := []string{}
	tagsToAdd := []string{}

	for _, origTag := range originalTagList {
		if !ArrayContains(newTagList, origTag) {
			// Orig tag no longer in tag list
			tagsToRemove = append(tagsToRemove, origTag)
		}
	}

	for _, newTag := range newTagList {
		if !ArrayContains(originalTagList, newTag) {
			// Orig tag no longer in tag list
			tagsToAdd = append(tagsToAdd, newTag)
		}
	}

	return tagsToRemove, tagsToAdd
}

func ArrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
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
