package build

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/fuxiaohei/pugo-static/app/models"
	"github.com/fuxiaohei/pugo-static/app/utils/mylog"
	"github.com/fuxiaohei/pugo-static/app/vars"
)

// Compile compile data in content object to webpages
func Compile(content *models.Content) bool {
	tasks, err := prepareTasks(content, content.Meta.Meta.Title)
	if err != nil {
		mylog.Error("prepare compile task error %s", err.Error())
		return false
	}
	failCount := 0
	compileFn := createBuildFunc(content)
	for _, task := range tasks {
		fileByts, err := compileFn(task.TplFile, task.TplData)
		if err != nil {
			mylog.Error("compile file %s error %s", task.FromFile, err.Error())
			failCount++
			continue
		}
		fDir := filepath.Dir(task.ToFile)
		os.MkdirAll(fDir, os.ModePerm)
		if err := ioutil.WriteFile(task.ToFile, fileByts, os.ModePerm); err != nil {
			mylog.Error("compile file %s error %s", task.FromFile, err.Error())
			failCount++
			continue
		}
		content.DstFiles[task.ToFile] = true
		mylog.Trace("compile file %s from %s", task.ToFile, task.FromFile)
	}
	mylog.Info("compile files %s, fails %s", len(tasks)-failCount, failCount)
	return true
}

func toDstFile(fpath string) string {
	return filepath.Join(vars.Config.DstDir, fpath)
}

type compileTask struct {
	FromFile string
	ToFile   string
	TplFile  string
	TplData  map[string]interface{}
}

func prepareTasks(content *models.Content, sizeTitle string) ([]compileTask, error) {
	var tasks []compileTask
	// prepare posts
	for _, post := range content.Posts {
		tasks = append(tasks, compileTask{
			FromFile: post.FromFile(),
			ToFile:   toDstFile(post.ToFile()),
			TplFile:  "post.html",
			TplData: map[string]interface{}{
				"Post":  post,
				"Title": post.Title + " - " + sizeTitle,
				"Slug":  post.ToFile(),
			},
		})
	}
	// prepare post lists
	for _, postList := range content.Lists {
		slug := fmt.Sprintf("/posts/%d.html", postList.Pager.Current)
		task := compileTask{
			FromFile: fmt.Sprintf("post-list-%d", postList.Pager.Current),
			TplFile:  "posts.html",
			ToFile:   toDstFile(slug),
			TplData: map[string]interface{}{
				"Pager": postList.Pager,
				"Posts": postList.Posts,
				"Title": fmt.Sprintf("Page %d - %s", postList.Pager.Current, sizeTitle),
				"Slug":  slug,
			},
		}
		tasks = append(tasks, task)
	}
	// prepare index page
	indexTask := compileTask{
		TplFile: "posts.html",
		ToFile:  toDstFile("index.html"),
		TplData: map[string]interface{}{
			"Pager": content.Lists[0].Pager,
			"Posts": content.Lists[0].Posts,
			"Hover": "index",
			"Slug":  "index.html",
		},
		FromFile: "index",
	}
	if content.Theme.Template("index.html") != nil {
		indexTask.TplFile = "index.html"
	}
	tasks = append(tasks, indexTask)
	// prepare post tag list
	for tag, lists := range content.PostTagLists {
		for _, list := range lists {
			slug := fmt.Sprintf("/tag/%s/%d.html", tag, list.Pager.Current)
			task := compileTask{
				TplFile: "posts.html",
				ToFile:  toDstFile(slug),
				TplData: map[string]interface{}{
					"Tag":   list.Tag,
					"Posts": list.Posts,
					"Title": fmt.Sprintf("Tag %s - Page %d - %s", tag, list.Pager.Current, sizeTitle),
					"Slug":  slug,
				},
				FromFile: fmt.Sprintf("post-tag-%s-list-%d", tag, list.Pager.Current),
			}
			tasks = append(tasks, task)
		}
	}
	// prepare pages
	for _, p := range content.Pages {
		task := compileTask{
			TplFile: "page.html",
			ToFile:  toDstFile(p.ToFile()),
			TplData: map[string]interface{}{
				"Page":  p,
				"Hover": p.Hover,
				"Title": p.Title + " - " + sizeTitle,
				"Slug":  p.URL(),
			},
			FromFile: p.FromFile(),
		}
		if p.Template != "" {
			task.TplFile = p.Template
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

type compileFunc func(tplFile string, tplData map[string]interface{}) ([]byte, error)

func createBuildFunc(content *models.Content) compileFunc {
	viewData := map[string]interface{}{
		"Meta":      content.Meta.Meta,
		"Nav":       content.Meta.Nav,
		"Hover":     "",
		"Now":       time.Now(),
		"Version":   vars.Version,
		"Title":     content.Meta.Meta.Title + " - " + content.Meta.Meta.Subtitle,
		"Keyword":   content.Meta.Meta.Keyword,
		"Desc":      content.Meta.Meta.Desc,
		"Comment":   content.Meta.Comment,
		"Analytics": content.Meta.Analytics,
		"Base":      content.Meta.Meta.Base(),
		"Slug":      "",
	}
	return func(tplFile string, tplData map[string]interface{}) ([]byte, error) {
		for k, v := range tplData {
			viewData[k] = v
		}
		buf := bytes.NewBuffer(nil)
		if err := content.Theme.Execute(buf, tplFile, viewData); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}
