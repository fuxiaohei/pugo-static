package build

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/fuxiaohei/pugo-static/app/models"
	"github.com/fuxiaohei/pugo-static/app/theme"
	"github.com/fuxiaohei/pugo-static/app/utils/mylog"
	"github.com/fuxiaohei/pugo-static/app/vars"
)

// Read read all data into content object
// if reading something wrongly, return false
func Read(content *models.Content) bool {
	// read metadata file first
	meta, err := models.ReadMetadataFile(vars.MetaFile)
	if err != nil {
		mylog.Error("read metadata file %s error %s", vars.MetaFile, err)
		return false
	}
	content.Meta = meta
	// read contents
	mylog.Info("read metadata file %s", vars.MetaFile)
	if err = readPosts(content); err != nil {
		mylog.Error("read posts error %s", err.Error())
		return false
	}
	if err = readPages(content); err != nil {
		mylog.Error("read pages error %s", err.Error())
		return false
	}
	// read theme
	if err = readTheme(content); err != nil {
		mylog.Error("read theme error %s", err.Error())
		return false
	}
	return true
}

func readPosts(cnt *models.Content) error {
	var posts []*models.Post
	err := filepath.Walk(vars.Config.PostDir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(fpath) != ".md" {
			return nil
		}
		post, err := models.NewPostFromFile(fpath)
		if err != nil {
			mylog.Warn("read post file %s error %s", fpath, err.Error())
			return nil
		}
		if post != nil {
			if post.IsDraft {
				mylog.Warn("read post draft %s", fpath)
				return nil
			}
			post.SetAuthor(cnt.Meta.FindAuthor(post.AuthorName))
			posts = append(posts, post)
			mylog.Trace("read post file %s", fpath)
		}
		return nil
	})
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Created().Unix() > posts[j].Created().Unix()
	})
	if err != nil {
		return err
	}
	cnt.Posts = posts
	mylog.Info("read posts %s", len(posts))
	cnt.Lists = models.NewPostLists(posts, vars.Config.PostSizePerPage, "")
	mylog.Info("build post lists %s", len(cnt.Lists))
	cnt.Tags, cnt.PostTagLists = models.NewPostTags(posts, vars.Config.PostSizePerPage)
	mylog.Info("build post tags %s", len(cnt.Tags))
	for t, list := range cnt.PostTagLists {
		mylog.Trace("build tag %s post list %s, posts %s", t, len(list), models.SizeOfPostLists(list))
	}
	cnt.Archives = models.NewArchives(posts)
	mylog.Info("build archives %s", len(cnt.Archives))
	for _, a := range cnt.Archives {
		mylog.Trace("build archive %s list, posts %s", a.Year, len(a.Posts))
	}
	return nil
}

func readPages(cnt *models.Content) error {
	var pages []*models.Page
	err := filepath.Walk(vars.Config.PageDir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(fpath) != ".md" {
			return nil
		}
		rel, err := filepath.Rel(vars.Config.PageDir, fpath)
		if err != nil {
			return err
		}
		page, err := models.NewPageFromFile(fpath, rel)
		if err != nil {
			mylog.Warn("read page file %s error %s", fpath, err.Error())
			return nil
		}
		if page != nil {
			if page.IsDraft {
				mylog.Warn("read page draft %s", fpath)
				return nil
			}
			page.SetAuthor(cnt.Meta.FindAuthor(page.AuthorName))
			pages = append(pages, page)
			mylog.Trace("read page file %s", fpath)
		}
		return nil
	})
	if err != nil {
		return err
	}
	cnt.Pages = pages
	mylog.Info("read pages %s", len(pages))
	return nil
}

func readTheme(cnt *models.Content) error {
	th, err := theme.New(vars.Config.ThemeDir)
	if err != nil {
		return err
	}
	cnt.Theme = th
	return nil
}
