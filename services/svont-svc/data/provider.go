package data

type PostOverview struct {
	Id                string   `json:"id"`
	Title             string   `json:"title"`
	Summary           string   `json:"summary"`
	Tags              []string `json:"tags"`
	AuthorId          string   `json:"authorId"`
	AuthorDisplayName string   `json:"authorDisplayName"`
	AuthorProfilePic  string   `json:"authorProfilePic"`
	Created           string   `json:"created"`
	Updated           string   `json:"updated"`
	Upvotes           int      `json:"upvotes"`
	CommentCount      int      `json:"commentCount"`
	FileCount         int      `json:"fileCount"`
}

type PostComment struct {
	Id                string        `json:"id"`
	Created           string        `json:"created"`
	Updated           string        `json:"updated"`
	AuthorId          string        `json:"authorId"`
	AuthorDisplayName string        `json:"authorDisplayName"`
	AuthorProfilePic  string        `json:"authorProfilePic"`
	Content           string        `json:"content"`
	Upvotes           int           `json:"upvotes"`
	Children          []PostComment `json:"children"`
}

type Post struct {
	Header  PostOverview `json:"header"`
	Content string       `json:"content"`
	Files   []string     `json:"files"`
}

type Provider interface {
	Initialize()
	Finalize()
	GetIndex() map[string]PostOverview
	GetPosts(start int, limit int) []PostOverview
	GetPost(postId string) *Post
	GetPostOverview(postId string) *PostOverview
	CreatePost(newPost Post, fileAttachments map[string][]byte) error
	UpdatePost(post Post, fileAttachments map[string][]byte) error
	UpvotePost(postId string, userEmail string) (*PostOverview, error)
	CreateComment(postId string, parentCommentId string, postComment *PostComment) (*[]PostComment, error)
	GetComments(postId string) (*[]PostComment, error)
	UpvoteComment(postId string, commentId string, userEmail string) (*PostComment, error)
	DeletePost(postId string) error
}
