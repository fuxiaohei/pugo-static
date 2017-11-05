package models

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/fuxiaohei/pugo-static/app/vars"
)

var (
	// ErrorFrontMetaParseFail means front-meta is parsed fail
	ErrorFrontMetaParseFail = errors.New("front-meta parse fail")
	// ErrorContentTimeLayoutFail means time string in content is parsed fail
	ErrorContentTimeLayoutFail = errors.New("content time parse fail")
)

// Post define blog post
type Post struct {
	Title      string   `toml:"title"`
	Slug       string   `toml:"slug"`
	Desc       string   `toml:"desc"`
	Date       string   `toml:"date"`
	UpdateDate string   `toml:"update_date"`
	AuthorName string   `toml:"author"`
	TagList    []string `toml:"tags"`
	IsDraft    bool     `toml:"draft"`

	created      time.Time
	updated      time.Time
	rawContent   []byte
	rawHTML      []byte
	briefContent []byte
	briefHTML    []byte
	tags         map[string]*PostTag
	author       *Author

	srcFile string
}

// NewPostFromFile read post from file
func NewPostFromFile(file string) (*Post, error) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	dataSlice := bytes.SplitN(fileBytes, []byte(vars.FrontMetaSeperator), 3)
	if len(dataSlice) != 3 {
		return nil, ErrorFrontMetaParseFail
	}
	post := new(Post)
	if err = parseFrontMeta(dataSlice[1], post); err != nil {
		return nil, err
	}
	post.srcFile = file
	post.rawContent = bytes.TrimSpace(dataSlice[2])
	post.tags = make(map[string]*PostTag)
	for _, tag := range post.TagList {
		tag = strings.TrimSpace(tag)
		post.tags[tag] = &PostTag{
			Name: tag,
			URL:  "/tag/" + tag + "/1.html",
		}
	}
	return post, post.formatFrontMeta()
}

func (p *Post) formatFrontMeta() error {
	var err error
	if p.created, err = parseContentTime(p.Date); err != nil {
		return err
	}
	if p.UpdateDate != "" {
		if p.updated, err = parseContentTime(p.UpdateDate); err != nil {
			return err
		}
	} else {
		p.updated = p.created
	}
	contentSlice := bytes.SplitN(p.rawContent, []byte(vars.ContentBriefSeperator), 2)
	if len(contentSlice) == 2 {
		p.briefContent = contentSlice[0]
	} else {
		p.briefContent = p.rawContent
	}
	return nil
}

func parseFrontMeta(data []byte, value interface{}) error {
	data = bytes.TrimSpace(data)
	for typeName, prefix := range vars.FrontMetaTypes {
		if !bytes.HasPrefix(data, []byte(prefix)) {
			continue
		}
		data = bytes.TrimPrefix(data, []byte(prefix))
		switch typeName {
		case "toml":
			return toml.Unmarshal(data, value)
		}
	}
	return ErrorFrontMetaParseFail
}

func parseContentTime(timeStr string) (time.Time, error) {
	for _, layout := range vars.ContentTimeLayouts {
		t, err := time.Parse(layout, timeStr)
		if err != nil {
			continue
		}
		return t, nil
	}
	return time.Unix(0, 0), ErrorContentTimeLayoutFail
}

// Created return create time of post
func (p *Post) Created() time.Time {
	return p.created
}

// Updated return update time of post
// if not set, same time to Updated()
func (p *Post) Updated() time.Time {
	return p.updated
}

// Content return post content bytes
func (p *Post) Content() []byte {
	return p.rawContent
}

// ContentHTML return post content bytes as template.HTML
func (p *Post) ContentHTML() template.HTML {
	if len(p.rawHTML) == 0 {
		p.rawHTML = Markdown(p.rawContent)
	}
	return template.HTML(p.rawHTML)
}

// Brief return brief content bytes
func (p *Post) Brief() []byte {
	return p.briefContent
}

// BriefHTML return brief content as template.HTML
func (p *Post) BriefHTML() template.HTML {
	if len(p.briefHTML) == 0 {
		p.briefHTML = Markdown(p.briefContent)
	}
	return template.HTML(p.briefHTML)
}

// FromFile return source code file of the post
func (p *Post) FromFile() string {
	return p.srcFile
}

// ToFile return destination webpage file path
func (p *Post) ToFile() string {
	if p.IsDraft {
		return ""
	}
	return fmt.Sprintf("%s/%s.html", p.created.Format("2006/1/2"), p.Slug)
}

// URL return visit link of the post
// same to ToFile()
func (p *Post) URL() string {
	return p.ToFile()
}

// PostTag is tag of a post
type PostTag struct {
	Name string
	URL  string
}

// Tags return tags of the post
func (p *Post) Tags() map[string]*PostTag {
	return p.tags
}

// SetAuthor set author object to the post
func (p *Post) SetAuthor(author *Author) {
	p.author = author
}

// Author return post's author
func (p *Post) Author() *Author {
	return p.author
}
