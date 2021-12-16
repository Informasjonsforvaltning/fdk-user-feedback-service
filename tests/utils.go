package tests

import (
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
)

func DeepEqualsPost(a *model.Post, b *model.Post) bool {
	if a == b {
		return true
	}

	if a.PostId != b.PostId && a.PostId != nil && *a.PostId != *b.PostId {
		return false
	}
	if a.UserId != b.UserId && a.UserId != nil && *a.UserId != *b.UserId {
		return false
	}
	if a.ThreadId != b.ThreadId && a.ThreadId != nil && *a.ThreadId != *b.ThreadId {
		return false
	}
	if a.Content != b.Content && a.Content != nil && *a.Content != *b.Content {
		return false
	}
	if a.ToPostId != b.ToPostId && a.ToPostId != nil && *a.ToPostId != *b.ToPostId {
		return false
	}
	if a.Timestamp != b.Timestamp && a.Timestamp != nil && *a.Timestamp != *b.Timestamp {
		return false
	}
	if a.Deleted != b.Deleted && a.Deleted != nil && *a.Deleted != *b.Deleted {
		return false
	}
	if a.UserInfo != b.UserInfo && a.UserInfo != nil && *a.UserInfo != *b.UserInfo {
		return false
	}

	return true
}
