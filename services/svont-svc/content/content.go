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
	"strings"
	"sync"
	"time"

	"tyayers/go-cms/data"

	"github.com/blevesearch/bleve/v2"
)

var postsMutex = sync.Mutex{}
var index data.PostIndex
var searchIndex bleve.Index
var tagIndex bleve.Index

// Local data provider, uncomment to test locally with files (in localdata dir)
var dataProvider data.Provider = &data.LocalProvider{}

// Google Cloud Storage provider, uncomment to use GCS as storage provider
// var dataProvider data.Provider = &data.GCSProvider{}

func Initialize(force bool) {
	fmt.Println("Starting loading indexes...")
	start := time.Now()
	index = InitializeProvider()

	elapsed := time.Since(start)
	fmt.Printf("Finished loading indexes in {%s}\n", elapsed)

	// Initialize bleve search index, if it doesn't exist
	if _, err := os.Stat("./posts.bleve"); os.IsNotExist(err) || force {
		// Initialize bleve
		fmt.Println("Starting building bleve index...")
		start = time.Now()
		os.RemoveAll("./posts.bleve")
		mapping := bleve.NewIndexMapping()
		var err error
		searchIndex, err = bleve.New("posts.bleve", mapping)
		if err != nil {
			fmt.Println(err)
		}

		// Loop and index posts in a separate thread
		go InitializeBleveIndex()
	} else {
		log.Printf("Loading local bleve index..\n")
		start = time.Now()
		searchIndex, err = bleve.Open("posts.bleve")

		elapsed = time.Since(start)
		fmt.Printf("Finished loading bleve index in {%s}\n", elapsed)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Initialize bleve tag index, if it doesn't exist
	if _, err := os.Stat("./tags.bleve"); os.IsNotExist(err) || force {

		os.RemoveAll("./tags.bleve")
		mapping := bleve.NewIndexMapping()
		var err error
		tagIndex, err = bleve.New("tags.bleve", mapping)
		if err != nil {
			fmt.Println(err)
		}

		// Loop and index tags in a separate thread
		go InitializeBleveTags()
	} else {
		log.Println("Loading local bleve tag index..")
		start = time.Now()
		tagIndex, err = bleve.Open("tags.bleve")
		elapsed = time.Since(start)
		fmt.Printf("Finished loading bleve tag index in {%s}\n", elapsed)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func InitializeBleveIndex() {
	log.Printf("Starting loading local bleve post index..\n")
	start := time.Now()
	count := 0
	for _, v := range index.Index {
		//fmt.Printf("Indexing key[%s] and index [%d]\n", k, count)
		searchIndex.Index(v.Id, v)
		count++
	}

	elapsed := time.Since(start)
	fmt.Printf("Finished indexing bleve posts in {%s}\n", elapsed)
}

func InitializeBleveTags() {
	log.Printf("Starting loading local bleve tag index..\n")
	start := time.Now()
	count := 0
	for k := range index.IndexTags {
		//fmt.Printf("Indexing tag key[%s] and index [%d]\n", k, count)
		tagIndex.Index(k, k)
		count++
	}

	elapsed := time.Since(start)
	fmt.Printf("Finished indexing bleve tags in {%s}\n", elapsed)
}

func Finalize(persistMode data.PersistMode) {
	fmt.Printf("Starting finalizing indexes for persist mode %d.\n", persistMode)
	start := time.Now()

	FinalizeProvider(persistMode, index)

	elapsed := time.Since(start)
	fmt.Printf("Finished finalizing indexes with persist mode %d in {%s}\n", persistMode, elapsed)
}

func GetData() data.Metadata {
	result := data.Metadata{PostCount: len(index.IndexTime), DraftCount: len(index.IndexDrafts), DeletedCount: len(index.IndexDeleted)}
	return result
}

func GetPosts(start int, limit int) []data.PostHeader {
	resultPosts := []data.PostHeader{}

	if len(index.IndexTime) > 0 {
		postIndex := len(index.IndexTime) - 1 - start

		for postIndex >= 0 && len(resultPosts) < limit {
			if !index.Index[index.IndexTime[postIndex]].Deleted && !index.Index[index.IndexTime[postIndex]].Draft {

				var post = index.Index[index.IndexTime[postIndex]]
				post.Upvotes = index.IndexCountLikes[post.Id]
				post.CommentCount = index.IndexCountComments[post.Id]
				resultPosts = append(resultPosts, post)
			}

			postIndex--
		}
	}

	return resultPosts
}

func GetPopularPosts(start int, limit int) []data.PostHeader {
	postsByPopularity := []data.PostHeader{}

	keys := make([]int, 0, len(index.IndexPopularityLikes))
	for k := range index.IndexPopularityLikes {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, v := range keys {
		for i := range index.IndexPopularityLikes[v] {
			postsByPopularity = append(postsByPopularity, index.Index[index.IndexPopularityLikes[v][i]])

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

func GetTaggedPosts(tagName string, start int, limit int) []data.PostHeader {
	taggedPosts := []data.PostHeader{}

	posts, ok := index.IndexTags[tagName]

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
			post, ok := index.Index[posts[keys[postIndex]]]
			if ok {
				taggedPosts = append(taggedPosts, post)
			}

			postIndex++
		}
	}

	return taggedPosts
}

func GetPost(postId string, draft bool) *data.Post {
	var post = GetPostFromProvider(postId, draft)
	post.Header = index.Index[postId]
	post.Header.Upvotes = index.IndexCountLikes[postId]
	post.Header.CommentCount = index.IndexCountComments[postId]
	return post
}

func GetPostOverview(postId string) data.PostHeader {
	return index.Index[postId]
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
	index.IndexTime = append(index.IndexTime, newPost.Header.Id)
	newPost.Header.Index = len(index.IndexTime) - 1

	// Add to draft index
	if newPost.Header.Draft {
		index.IndexDrafts[newPost.Header.Id] = newPost.Header.Index
	} else {
		// Add to popularity indexes
		index.IndexPopularityLikes[0] = append(index.IndexPopularityLikes[0], newPost.Header.Id)
		index.IndexPopularityComments[0] = append(index.IndexPopularityComments[0], newPost.Header.Id)
		index.IndexPopularityViews[0] = append(index.IndexPopularityViews[0], newPost.Header.Id)

		// Add to tag index
		if newPost.Header.Tags != nil {
			for _, tag := range newPost.Header.Tags {

				if tag != "" {
					_, ok := index.IndexTags[tag]

					if !ok {
						index.IndexTags[tag] = map[int]string{}
						tagIndex.Index(tag, tag)
					}

					index.IndexTags[tag][newPost.Header.Index] = newPost.Header.Id
				}
			}
		}
	}

	// Add to id index
	index.Index[newPost.Header.Id] = newPost.Header

	postsMutex.Unlock()
	// Persist changes to storage in the background
	go Finalize(data.PersistAll)

	err := CreatePostForProvider(*newPost, files)

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

		updatedPost.Files = append(updatedPost.Files, attachment.Filename)

		if updatedPost.Header.Image == "" && (strings.HasSuffix(strings.ToLower(attachment.Filename), "png") || strings.HasSuffix(strings.ToLower(attachment.Filename), "jpg")) {
			updatedPost.Header.Image = attachment.Filename
		}
	}

	updatedPost.Header.FileCount = len(updatedPost.Files)

	postsMutex.Lock()
	header := index.Index[updatedPost.Header.Id]
	header.Title = updatedPost.Header.Title

	if !updatedPost.Header.Draft {
		header.Summary = updatedPost.Header.Summary
	}

	if header.Image == "" {
		header.Image = updatedPost.Header.Image
	}

	if header.Draft && !updatedPost.Header.Draft {
		// post is no longer in draft, remove from draft index
		delete(index.IndexDrafts, header.Id)
		header.Draft = updatedPost.Header.Draft

		// delete draft post file
		dataProvider.DeleteFile("data/" + header.Id + "/post_draft.json")

		// Add to popularity indexes
		index.IndexPopularityLikes[0] = append(index.IndexPopularityLikes[0], updatedPost.Header.Id)
		index.IndexPopularityComments[0] = append(index.IndexPopularityComments[0], updatedPost.Header.Id)
		index.IndexPopularityViews[0] = append(index.IndexPopularityViews[0], updatedPost.Header.Id)

		// Add to tag index
		if updatedPost.Header.Tags != nil {
			for _, tag := range updatedPost.Header.Tags {

				if tag != "" {
					_, ok := index.IndexTags[tag]

					if !ok {
						index.IndexTags[tag] = map[int]string{}
						tagIndex.Index(tag, tag)
					}

					index.IndexTags[tag][updatedPost.Header.Index] = updatedPost.Header.Id
				}
			}
		}

		// Persist changes to storage in the background
		go Finalize(data.PersistAll)
	} else {
		UpdateTags(header.Id, header.Index, header.Tags, updatedPost.Header.Tags)
		header.Tags = updatedPost.Header.Tags
		header.FileCount = updatedPost.Header.FileCount
		header.Updated = updatedPost.Header.Updated
	}

	index.Index[updatedPost.Header.Id] = header
	postsMutex.Unlock()

	isDraft := updatedPost.Header.Draft
	updatedPost.Header = header

	err := UpdatePostForProvider(*updatedPost, files, isDraft)

	if err != nil {
		return err
	} else {

		// Create successful, now index for search as well
		searchIndex.Index(updatedPost.Header.Id, updatedPost.Header)

		return nil
	}
}

func UpvotePost(postId string, userEmail string) (*data.PostHeader, error) {

	post, ok := index.Index[postId]

	if ok {
		postsMutex.Lock()

		index.IndexCountLikes[postId]++
		voteCount := index.IndexCountLikes[postId]

		// Remove item from previous popularity space
		for i, s := range index.IndexPopularityLikes[voteCount-1] {
			if s == post.Id {
				// We found our post in the old spot, now remove
				index.IndexPopularityLikes[voteCount-1][i] = index.IndexPopularityLikes[voteCount-1][len(index.IndexPopularityLikes[voteCount-1])-1] // Copy last element to index i.
				index.IndexPopularityLikes[voteCount-1][len(index.IndexPopularityLikes[voteCount-1])-1] = ""                                         // Erase last element (write zero value).
				index.IndexPopularityLikes[voteCount-1] = index.IndexPopularityLikes[voteCount-1][:len(index.IndexPopularityLikes[voteCount-1])-1]   // Truncate slice.
			}
		}

		// Add to new popularity spot
		val, ok := index.IndexPopularityLikes[voteCount]
		// If the key exists
		if ok {
			index.IndexPopularityLikes[voteCount] = append(val, post.Id)
		} else {
			index.IndexPopularityLikes[voteCount] = []string{post.Id}
		}

		postsMutex.Unlock()
		// Persist changes to storage in the background
		go Finalize(data.PersistOnlyCountLikes)
		go Finalize(data.PersistOnlyPopularityLikes)

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

	post, ok := index.Index[postId]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Post %s not found!", postId))
	} else {
		result, err := CreateComment(postId, parentCommentId, newComment)

		if err == nil {
			postsMutex.Lock()
			post.CommentCount++
			index.IndexCountComments[postId]++
			index.Index[postId] = post
			postsMutex.Unlock()
			// Persist changes to storage in the background
			go Finalize(data.PersistOnlyCountComments)

			return result, nil
		} else {
			return nil, err
		}
	}
}

// Gets all of the comments for a post
func GetComments(postId string) (*[]data.PostComment, error) {
	return GetCommentsFromProvider(postId)
}

// Upvotes a specific comment
func UpvoteComment(postId string, commentId string, userEmail string) (*data.PostComment, error) {
	return UpvoteCommentForProvider(postId, commentId, userEmail)
}

// Gets a file attachment for a post
func GetFileForPost(postId string, fileName string) ([]byte, error) {
	return GetFile(postId, fileName)
}

// Deletes a post
func DeletePost(postId string) error {

	postsMutex.Lock()
	post, ok := index.Index[postId]

	if ok {
		post.Deleted = true
		index.Index[postId] = post
		voteCount := index.IndexCountLikes[postId]

		// Remove item from popularity spaces
		for i, s := range index.IndexPopularityLikes[voteCount] {
			if s == post.Id {
				// We found our post in the old spot, now remove
				index.IndexPopularityLikes[voteCount][i] = index.IndexPopularityLikes[voteCount][len(index.IndexPopularityLikes[voteCount])-1] // Copy last element to index i.
				index.IndexPopularityLikes[voteCount][len(index.IndexPopularityLikes[voteCount])-1] = ""                                       // Erase last element (write zero value).
				index.IndexPopularityLikes[voteCount] = index.IndexPopularityLikes[voteCount][:len(index.IndexPopularityLikes[voteCount])-1]   // Truncate slice.
			}
		}

		// Remove from tags collection
		for _, removeTag := range post.Tags {
			val, ok := index.IndexTags[removeTag]
			// If the key exists
			if ok {
				delete(val, post.Index)
			} else {
				// Log error
			}
		}

		// Remove from search
		if searchIndex != nil {
			searchIndex.Delete(postId)
		}

		// Remove from drafts
		delete(index.IndexDrafts, postId)

		// Add to deleted index, if it wasn't a draft
		if !post.Draft {
			index.IndexDeleted[postId] = post.Index
		}
	}

	postsMutex.Unlock()
	// Persist changes to storage in the background
	go Finalize(data.PersistAll)

	return DeletePostForProvider(postId)
}

// Searches posts
func SearchPosts(text string) ([]data.PostHeader, error) {
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
			results := []data.PostHeader{}
			fmt.Println(searchResults)
			dataMap := index.Index
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
				tag, ok := index.IndexTags[val.ID]
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
		val, ok := index.IndexTags[removeTag]
		// If the key exists
		if ok {
			delete(val, postIndex)
		} else {
			// Log error
		}
	}

	for _, addTag := range tagsToAdd {
		if addTag != "" {
			_, ok := index.IndexTags[addTag]

			if !ok {
				index.IndexTags[addTag] = map[int]string{}
				tagIndex.Index(addTag, addTag)
			}

			index.IndexTags[addTag][postIndex] = postId
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
