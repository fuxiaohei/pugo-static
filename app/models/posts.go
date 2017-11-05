package models

import (
	"strings"
)

// PostList is list of posts
type PostList struct {
	Tag   string
	Posts []*Post
	Pager *Pager
}

// NewPostLists build lists of posts with size
func NewPostLists(posts []*Post, size int) []*PostList {
	cursor := NewPagerCursor(size, len(posts))
	var lists []*PostList
	for _, pg := range cursor.Pages() {
		list := &PostList{
			Tag:   "",
			Pager: pg,
			Posts: posts[pg.Begin:pg.End],
		}
		lists = append(lists, list)
	}
	return lists
}

// NewPostTags build post tags and post lists of tag
func NewPostTags(posts []*Post, size int) (map[string]string, map[string][]*PostList) {
	m := make(map[string]string)
	tmp := make(map[string][]*Post)
	for _, p := range posts {
		for _, tag := range p.TagList {
			tag = strings.TrimSpace(tag)
			m[tag] = tag
			tmp[tag] = append(tmp[tag], p)
		}
	}
	lists := make(map[string][]*PostList)
	for tag, posts := range tmp {
		lists[tag] = NewPostLists(posts, size)
	}
	return m, lists
}

// SizeOfPostLists get all posts count from post lists
func SizeOfPostLists(lists []*PostList) int {
	if len(lists) == 0 {
		return 0
	}
	return lists[0].Pager.All
}
