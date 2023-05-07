package data

type PostOverview struct {
	Id                string   `json:"id"`
	Title             string   `json:"title"`
	Summary           string   `json:"summary"`
	Image             string   `json:"image"`
	Draft             bool     `json:"draft"`
	Deleted           bool     `jason:"deleted"`
	Tags              []string `json:"tags"`
	AuthorId          string   `json:"authorId"`
	AuthorDisplayName string   `json:"authorDisplayName"`
	AuthorProfilePic  string   `json:"authorProfilePic"`
	Created           string   `json:"created"`
	Index             int      `json:"index"`
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

type SearchResult struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Count int    `json:"count"`
}

type ImageUploadResult struct {
	Url string `json:"url"`
}

type Provider interface {
	Initialize() (map[string]PostOverview, []string, map[int][]string, map[string]map[int]string)
	Finalize(index_main map[string]PostOverview, index_time []string, index_populary map[int][]string, index_tags map[string]map[int]string)

	GetPost(postId string) *Post
	CreatePost(newPost Post, fileAttachments map[string][]byte) error
	UpdatePost(post Post, fileAttachments map[string][]byte) error
	CreateComment(postId string, parentCommentId string, postComment *PostComment) (*[]PostComment, error)
	GetComments(postId string) (*[]PostComment, error)
	UpvoteComment(postId string, commentId string, userEmail string) (*PostComment, error)

	GetFile(postId string, fileName string) ([]byte, error)
	DeletePost(postId string) error
}
