package models

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/BurntSushi/toml"
)

type (
	// MetaData is metadata from meta file
	MetaData struct {
		Meta      *MetaInfo `toml:"meta"`
		Nav       NavList   `toml:"nav"`
		Authors   []Author  `toml:"author"`
		Comment   Comment   `toml:"comment"`
		Analytics Analytics `toml:"analytics"`
		srcFile   string
	}
	// MetaInfo is basic info from metadata to describe the website
	MetaInfo struct {
		Title    string `toml:"title"`
		Subtitle string `toml:"subtitle"`
		Keyword  string `toml:"keyword"`
		Desc     string `toml:"desc"`
		Root     string `toml:"root"`
		Lang     string `toml:"lang"`

		urlObject *url.URL
	}
)

// FindAuthor find author by name of nickname
// if not found, return first author
func (md *MetaData) FindAuthor(name string) *Author {
	if name == "" {
		return &(md.Authors[0])
	}
	for _, au := range md.Authors {
		if au.Name == name {
			return &au
		}
		if au.Nick == name {
			return &au
		}
	}
	return &(md.Authors[0])
}

// Base return base path of the root url
// use for subdirectory url
func (minfo *MetaInfo) Base() string {
	if minfo.urlObject == nil {
		return ""
	}
	return minfo.urlObject.Path
}

// ReadMetadata read metadata from bytes
// metadata's src file is blank
func ReadMetadata(data []byte) (*MetaData, error) {
	m := new(MetaData)
	if err := toml.Unmarshal(data, m); err != nil {
		return nil, err
	}
	urlObject, err := url.Parse(m.Meta.Root)
	if err != nil {
		return nil, fmt.Errorf("meta's root is wrong, %s", err.Error())
	}
	m.Meta.urlObject = urlObject
	return m, nil
}

// ReadMetadataFile read metadata from file
// metadata's src file is set
func ReadMetadataFile(file string) (*MetaData, error) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	m, err := ReadMetadata(fileBytes)
	if err == nil {
		m.srcFile = file
	}
	return m, err
}

type (
	// Nav is an item for a navigator link
	Nav struct {
		Link    string `toml:"link"`
		Title   string `toml:"title"`
		Hover   string `toml:"hover"`
		I18n    string `toml:"i18n"`
		IsBlank bool   `toml:"is_blank"`
		Icon    string `toml:"icon"`
	}
	// NavList is list of Nav items
	NavList []Nav
)

type (
	// Author is author item
	Author struct {
		Name  string `toml:"name"`
		Nick  string `toml:"nick"`
		Email string `toml:"email"`
		URL   string `toml:"url"`
	}
)

type (
	// Comment is setting of third-party comment system
	Comment struct {
		Disqus  string `toml:"disqus"`
		Duoshuo string `toml:"duoshuo"`
	}
	// Analytics is setting of third-party website analytics tool
	Analytics struct {
		Google  string `toml:"google"`
		Baidu   string `toml:"baidu"`
		Tencent string `toml:"tencent"`
		Cnzz    string `toml:"cnzz"`
	}
)
