package vars

const (
	// Version is version number of pugo.static
	Version = "1.0.0"
)

var (
	// MetaFile is filename of metadata
	MetaFile = "meta.toml"

	// FrontMetaSeperator is seperator of frontmeta and post content
	FrontMetaSeperator = "```"

	// FrontMetaTypes contains types and prefixes of front meta
	FrontMetaTypes = map[string]string{
		"toml": "toml",
	}

	// ContentTimeLayouts is time layouts for posts and pages
	ContentTimeLayouts = []string{
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
	}

	// ContentBriefSeperator is seperator of brief and all content
	ContentBriefSeperator = "<!--more-->"
)

var (
	// ThemeMetadataFile is default theme's metadata filename
	ThemeMetadataFile = "theme.toml"
)

var (
	// Config is default global config
	Config = ConfigData{
		PageDir:             "page",
		PostDir:             "post",
		ThemeDir:            "theme",
		LangDir:             "lang",
		MediaDir:            "media",
		DstDir:              "dest",
		PostSizePerPage:     5,
		PostTemplateFile:    "post.html",
		PageTemplateFile:    "page.html",
		IndexTemplateFile:   "index.html",
		ArchiveTemplateFile: "archive.html",
	}
)

// ConfigData is configuration data
type ConfigData struct {
	PostDir             string
	PageDir             string
	ThemeDir            string
	LangDir             string
	MediaDir            string
	PostSizePerPage     int
	PostTemplateFile    string
	PageTemplateFile    string
	IndexTemplateFile   string
	ArchiveTemplateFile string
	DstDir              string
}
