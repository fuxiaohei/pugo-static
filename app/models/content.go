package models

import (
	"github.com/fuxiaohei/pugo-static/app/theme"
)

// Content is all data collection
type Content struct {
	Meta         *MetaData
	Posts        []*Post
	Lists        []*PostList
	Tags         map[string]string
	PostTagLists map[string][]*PostList
	Pages        []*Page
	Theme        *theme.Theme
	DstFiles     map[string]bool
}

// NewContent return new content object with proper default values
func NewContent() *Content {
	return &Content{
		DstFiles: make(map[string]bool),
	}
}
