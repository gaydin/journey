package structure

import (
	"time"
)

type Post struct {
	Id              int64
	Uuid            string
	Title           string
	Slug            string
	Markdown        []byte
	Html            []byte
	IsFeatured      bool
	IsPage          bool
	IsPublished     bool
	Date            *time.Time
	Tags            []Tag
	Author          *User
	MetaDescription string
	Image           string
}
