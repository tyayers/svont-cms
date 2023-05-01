package data

import (
	"math/rand"
	"time"
)

// Adds the comment to the correct parent
func AddCommentToParent(comments *[]PostComment, parentCommentId string, newComment *PostComment) bool {
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
func UpvoteComment(comments *[]PostComment, commentId string) *PostComment {
	for i := 0; i < len(*comments); i++ {
		if (*comments)[i].Id == commentId {
			// We found it!
			(*comments)[i].Upvotes++

			return &(*comments)[i]
		} else if len((*comments)[i].Children) > 0 {
			result := UpvoteComment(&(*comments)[i].Children, commentId)
			if result != nil {
				return result
			}
		}
	}

	return nil
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
