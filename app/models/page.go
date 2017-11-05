package models

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/fuxiaohei/pugo-static/app/vars"
)

// Page define common web page
type Page struct {
	Title      string `toml:"title"`
	Desc       string `toml:"desc"`
	Date       string `toml:"date"`
	UpdateDate string `toml:"update_date"`
	AuthorName string `toml:"author"`
	IsDraft    bool   `toml:"draft"`
	Hover      string `toml:"hover"`
	I18n       string `toml:"i18n"`
	Template   string `toml:"template"`

	created    time.Time
	updated    time.Time
	rawContent []byte
	rawHTML    []byte
	author     *Author

	srcFile string
	relPath string
}

// NewPageFromFile read page object from file
func NewPageFromFile(file, rel string) (*Page, error) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	dataSlice := bytes.SplitN(fileBytes, []byte(vars.FrontMetaSeperator), 3)
	if len(dataSlice) != 3 {
		return nil, ErrorFrontMetaParseFail
	}
	page := new(Page)
	if err = parseFrontMeta(dataSlice[1], page); err != nil {
		return nil, err
	}
	page.relPath = rel
	page.srcFile = file
	page.rawContent = bytes.TrimSpace(dataSlice[2])
	return page, page.formatFrontMeta()
}

func (p *Page) formatFrontMeta() error {
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
	return nil
}

// Created return create time of page
func (p *Page) Created() time.Time {
	return p.created
}

// Updated return update time of page
// if not set, same time to Created()
func (p *Page) Updated() time.Time {
	return p.updated
}

// Content return page content bytes
func (p *Page) Content() []byte {
	return p.rawContent
}

// ContentHTML return page content bytes as template.HTML
func (p *Page) ContentHTML() template.HTML {
	if len(p.rawHTML) == 0 {
		p.rawHTML = Markdown(p.rawContent)
	}
	return template.HTML(p.rawHTML)
}

// FromFile return source code file of the post
func (p *Page) FromFile() string {
	return p.srcFile
}

// ToFile return destination webpage file path
func (p *Page) ToFile() string {
	if p.IsDraft {
		return ""
	}
	sfx := filepath.Ext(p.relPath)
	return strings.TrimSuffix(p.relPath, sfx) + ".html"
}

// URL return visit link of the page
func (p *Page) URL() string {
	return p.ToFile()
}

// SetAuthor set author object to the page
func (p *Page) SetAuthor(author *Author) {
	p.author = author
}

// Author return page's author
func (p *Page) Author() *Author {
	return p.author
}
