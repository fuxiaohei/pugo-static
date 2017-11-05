package build

import (
	"os"
	"path/filepath"

	"github.com/fuxiaohei/pugo-static/app/models"
	"github.com/fuxiaohei/pugo-static/app/utils/mylog"
	"github.com/fuxiaohei/pugo-static/app/vars"
)

// Clean clean destination directry with not-compiling file
func Clean(content *models.Content) {
	var count int
	err := filepath.Walk(vars.Config.DstDir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err

		}
		if !content.DstFiles[fpath] {
			os.Remove(fpath)
			mylog.Trace("clean dst file %s", fpath)
			count++
		}
		return nil
	})
	if err != nil {
		mylog.Error("clean contents error %s", err.Error())
	}
	mylog.Info("clean contents, remove %s file", count)
}
