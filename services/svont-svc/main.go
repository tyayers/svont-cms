package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"tyayers/go-cms/content"
	"tyayers/go-cms/data"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
)

var app *firebase.App

func getPosts(c *gin.Context) {
	start, err := strconv.Atoi(c.Query("start"))
	if err != nil {
		start = 0
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}

	c.IndentedJSON(http.StatusOK, content.GetPosts(start, limit))
}

func getPopularPosts(c *gin.Context) {
	start, err := strconv.Atoi(c.Query("start"))
	if err != nil {
		start = 0
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	c.IndentedJSON(http.StatusOK, content.GetPopularPosts(start, limit))
}

func getPost(c *gin.Context) {
	postId := c.Param("id")

	post := content.GetPost(postId)
	c.IndentedJSON(http.StatusOK, post)
}

func createPost(c *gin.Context) {
	newPost := data.Post{}

	var files []multipart.FileHeader

	form, _ := c.MultipartForm()
	attachments := form.File["files"]

	for _, file := range attachments {
		files = append(files, *file)
	}

	for key, value := range c.Request.PostForm {
		fmt.Printf("%v = %v \n", key, value)
		switch key {
		case "title":
			newPost.Header.Title = value[0]
		case "content":
			newPost.Content = value[0]
		case "summary":
			newPost.Header.Summary = value[0]
		case "authorId":
			newPost.Header.AuthorId = value[0]
		case "authorDisplayName":
			newPost.Header.AuthorDisplayName = value[0]
		case "authorProfilePic":
			newPost.Header.AuthorProfilePic = value[0]
		default:
			fmt.Println("No handler found for form item " + key)
		}
	}

	err := content.CreatePost(&newPost, files)

	if err != nil {
		c.String(500, fmt.Sprintf("Post could not be created! More info: %s", err.Error()))
	} else {
		fmt.Println(fmt.Sprintf("Post '%s' created.", newPost.Header.Id))

		c.IndentedJSON(http.StatusCreated, newPost)
	}

}

func updatePost(c *gin.Context) {
	postId := c.Param("id")
	user_id := c.GetString("user_id")

	post := content.GetPostOverview(postId)

	if post.AuthorId != user_id {
		c.String(401, fmt.Sprintf("User not authorized to delete post."))
	} else {
		updatedPost := data.Post{}
		updatedPost.Header.Id = postId

		var files []multipart.FileHeader

		form, _ := c.MultipartForm()
		attachments := form.File["files"]

		for _, file := range attachments {
			files = append(files, *file)
		}

		for key, value := range c.Request.PostForm {
			fmt.Printf("%v = %v \n", key, value)
			switch key {
			case "title":
				updatedPost.Header.Title = value[0]
			case "content":
				updatedPost.Content = value[0]
			case "summary":
				updatedPost.Header.Summary = value[0]
			default:
				fmt.Println("No handler found for form item " + key)
			}
		}

		content.UpdatePost(&updatedPost, files)

		c.IndentedJSON(http.StatusOK, updatedPost)
	}
}

func upvotePost(c *gin.Context) {
	postId := c.Param("id")
	userEmail := c.GetString("userEmail")

	post, err := content.UpvotePost(postId, userEmail)

	if err != nil {
		c.String(500, fmt.Sprintf("Could not upvote post, more info: %s", err.Error()))
	} else {
		fmt.Println(fmt.Sprintf("Post '%s' upvoted.", postId))
		c.IndentedJSON(http.StatusCreated, post)
	}
}

// Add a comment to a post, either at the end or under parentCommentId.
func addCommentToPost(c *gin.Context) {

	postId := c.Param("id")

	var newContent, authorId, authorDisplayName, authorProfilePic, parentCommentId string

	form, _ := c.MultipartForm()
	for key, value := range form.Value {
		fmt.Printf("%v = %v \n", key, value)
		switch key {
		case "authorId":
			authorId = value[0]
		case "authorDisplayName":
			authorDisplayName = value[0]
		case "authorProfilePic":
			authorProfilePic = value[0]
		case "content":
			newContent = value[0]
		case "parentCommentId":
			parentCommentId = value[0]
		default:
			fmt.Println("No handler found for form item " + key)
		}
	}

	newComment, err := content.AddCommentToPost(postId, parentCommentId, authorId, authorDisplayName, authorProfilePic, newContent)

	if err != nil {
		fmt.Println(fmt.Sprintf("Comment could not be created! More info: %s", err.Error()))
		c.String(500, fmt.Sprintf("Comment could not be created! More info: %s", err.Error()))
	} else {
		fmt.Println(fmt.Sprintf("Post '%s' comment created.", postId))

		c.IndentedJSON(http.StatusCreated, newComment)
	}

}

func getCommentsForPost(c *gin.Context) {
	postId := c.Param("id")

	comments, err := content.GetComments(postId)

	if err != nil {
		c.String(500, fmt.Sprintf("Could not get comments! More info: %s", err.Error()))
	} else {
		c.IndentedJSON(http.StatusOK, comments)
	}
}

func upvoteComment(c *gin.Context) {
	postId := c.Param("id")
	commentId := c.Param("commentId")
	userEmail := c.GetString("userEmail")

	post, err := content.UpvoteComment(postId, commentId, userEmail)

	if err != nil {
		c.String(500, fmt.Sprintf("Could not upvote post, more info: %s", err.Error()))
	} else {
		fmt.Println(fmt.Sprintf("Post '%s' upvoted.", postId))
		c.IndentedJSON(http.StatusCreated, post)
	}
}

// Attach a file to a post.
func attachFileToPost(c *gin.Context) {

	postId := c.Param("id")
	user_id := c.GetString("user_id")

	post := content.GetPostOverview(postId)

	if post.AuthorId != user_id {
		c.String(401, fmt.Sprintf("User not authorized to delete post."))
	} else {
		file, err := c.FormFile("file")
		// The file cannot be received.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}
		src, _ := file.Open()
		defer src.Close()

		byteContainer, err := ioutil.ReadAll(src)

		content.AttachFileToPost(postId, file.Filename, byteContainer)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	}
}

// Search posts
func searchPosts(c *gin.Context) {
	text := c.Query("q")

	res, err := content.SearchPosts(text)

	if err == nil {
		c.IndentedJSON(http.StatusOK, res)
	} else {
		c.String(500, fmt.Sprintf("Error: %s", err))
	}
}

// Delete a post
func deletePost(c *gin.Context) {

	postId := c.Param("id")
	user_id := c.GetString("user_id")

	post := content.GetPostOverview(postId)

	if post.AuthorId != user_id {
		// Reject, user's can delete other user's posts
		c.String(401, fmt.Sprintf("User not authorized to delete post."))
	} else {
		err := content.DeletePost(postId)

		if err == nil {
			c.String(http.StatusOK, fmt.Sprintf("'%s' deleted!", postId))
		} else {
			c.String(500, fmt.Sprintf("'%s' could not be deleted!", postId))
		}
	}
}

func jwtValidation() gin.HandlerFunc {
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	return func(c *gin.Context) {

		var idToken = c.Request.Header["Authorization"]

		if len(idToken) > 0 {
			cleanedToken := strings.ReplaceAll(idToken[0], "Bearer ", "")
			token, err := client.VerifyIDToken(context.Background(), cleanedToken)
			if err != nil {
				log.Printf("Error verifying ID token: %v, rejecting.\n", err)
				c.AbortWithStatus(401)
			} else {
				//log.Printf("token claims %v", token)
				c.Set("userEmail", token.Claims["email"])
				c.Set("user_id", token.Claims["user_id"])
			}
		} else {
			log.Printf("No id token found, rejecting.")
			c.AbortWithStatus(401)
		}

		c.Next()
	}
}

func main() {

	content.Initialize()

	app, _ = firebase.NewApp(context.Background(), nil)

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			fmt.Println("Interrupt received, persisting index and closing.")
			content.Finalize()
			os.Exit(1)
		case syscall.SIGKILL:
			fmt.Println("SIGINT received, persisting index and closing.")
			content.Finalize()
			os.Exit(1)
		case syscall.SIGTERM:
			fmt.Println("SIGTERM received, persisting index and closing.")
			content.Finalize()
			os.Exit(1)
		}
	}()

	router := gin.Default()

	router.Use(CORSMiddleware())
	//router.Use(jwtValidation())

	router.GET("/posts", getPosts)
	router.GET("/posts/popular", getPopularPosts)
	router.GET("/posts/search", searchPosts)
	router.GET("/posts/:id", getPost)
	router.POST("/posts", jwtValidation(), createPost)
	router.PUT("/posts/:id", jwtValidation(), updatePost)
	router.POST("/posts/:id/upvote", jwtValidation(), upvotePost)
	router.POST("/posts/:id/files", jwtValidation(), attachFileToPost)
	router.GET("/posts/:id/comments", getCommentsForPost)
	router.POST("/posts/:id/comments", jwtValidation(), addCommentToPost)
	router.POST("/posts/:id/comments/:commentId/upvote", jwtValidation(), upvoteComment)
	router.DELETE("/posts/:id", jwtValidation(), deletePost)

	router.Run("0.0.0.0:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
