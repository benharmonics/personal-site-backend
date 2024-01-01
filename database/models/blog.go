package models

import (
	"time"

	"github.com/benharmonics/personal-site-backend/api/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	BlogPost struct {
		ID          primitive.ObjectID `json:"id" bson:"_id"`
		Title       string             `json:"title" bson:"title"`
		Subtitle    *string            `json:"subtitle" bson:"subtitle,omitempty"`
		Author      string             `json:"author" bson:"author"`
		Content     string             `json:"content,omitempty" bson:"content"`
		Tags        []string           `json:"tags,omitempty" bson:"tags,omitempty"`
		DateCreated time.Time          `json:"dateCreated" bson:"dateCreated"`
		LastUpdate  time.Time          `json:"lastUpdate" bson:"lastUpdate"`
	}

	BlogPostOption func(*BlogPost)
)

func NewBlogPost(opts ...BlogPostOption) BlogPost {
	now := time.Now()
	post := BlogPost{
		ID:          primitive.NewObjectID(),
		Title:       "unknown",
		Author:      "unknown",
		Content:     "unknown",
		DateCreated: now,
		LastUpdate:  now,
	}
	for _, optFunc := range opts {
		optFunc(&post)
	}
	return post
}

func FromRequest(req requests.NewBlogPost) BlogPostOption {
	return func(post *BlogPost) {
		post.Author = req.Author
		post.Title = req.Title
		post.Subtitle = req.Subtitle
		post.Content = req.Content
	}
}

func WithTitle(title string) BlogPostOption {
	return func(post *BlogPost) {
		post.Title = title
	}
}

func WithSubtitle(subtitle string) BlogPostOption {
	return func(post *BlogPost) {
		post.Subtitle = &subtitle
	}
}

func WithAuthor(author string) BlogPostOption {
	return func(post *BlogPost) {
		post.Author = author
	}
}

func WithContent(content string) BlogPostOption {
	return func(post *BlogPost) {
		post.Content = content
	}
}

func WithTags(tags []string) BlogPostOption {
	return func(post *BlogPost) {
		post.Tags = tags
	}
}
