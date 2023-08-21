package data

type PostIndex struct {
	Index                   map[string]PostHeader
	IndexTime               []string
	IndexDrafts             map[string]int
	IndexDeleted            map[string]int
	IndexTags               map[string]map[int]string
	IndexPopularityLikes    map[int][]string
	IndexPopularityViews    map[int][]string
	IndexPopularityComments map[int][]string
	IndexCountLikes         map[string]int
	IndexCountComments      map[string]int
	IndexCountViews         map[string]int
	IndexUsers              map[string]User
}

type PostHeader struct {
	Id                string   `json:"id"`
	Title             string   `json:"title"`
	Summary           string   `json:"summary"`
	Image             string   `json:"image"`
	Draft             bool     `json:"draft"`
	Deleted           bool     `json:"deleted"`
	Tags              []string `json:"tags"`
	AuthorId          string   `json:"authorId"`
	AuthorDisplayName string   `json:"authorDisplayName"`
	AuthorProfilePic  string   `json:"authorProfilePic"`
	Author            User     `json:"author"`
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
	Header  PostHeader `json:"header"`
	Content string     `json:"content"`
	Files   []string   `json:"files"`
}

type User struct {
	UID           string `json:"uid"`
	DisplayName   string `json:"displayName"`
	Email         string `json:"-"`
	ProfileText   string `json:"profileText"`
	PhotoURL      string `json:"photoURL"`
	ProviderId    string `json:"-"`
	EmailVerified bool   `json:"-"`
	IsAnonymous   bool   `json:"-"`
	Registered    string `json:"registered"`
	Followers     int    `json:"followers"`
	Following     int    `json:"following"`
	PostCount     int    `json:"postCount"`
}

type SearchResult struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Count int    `json:"count"`
}

type Metadata struct {
	PostCount    int `json:"postCount"`
	DeletedCount int `json:"deletedCount"`
	DraftCount   int `json:"draftCount"`
	UserCount    int `json:"userCount"`
}

type ImageUploadResult struct {
	Url string `json:"url"`
}

type Provider interface {
	Initialize()
	Finalize(peristMode PersistMode, index PostIndex)
	CreateDir(dirName string) error
	UploadFile(fileName string, content []byte) error
	DownloadFile(fileName string) ([]byte, error)
	DeleteFile(fileName string) error
}

type PersistMode int

const (
	PersistAll = iota
	PersistOnlyHeaders
	PersistOnlyTime
	PersistOnlyTags
	PersistOnlyPopularityLikes
	PersistOnlyPopularityComments
	PersistOnlyPopularityViews
	PersistOnlyCountLikes
	PersistOnlyCountComments
	PersistOnlyCountViews
	PersistOnlyDrafts
	PersistOnlyDeleted
	PersistOnlyUsers
)
