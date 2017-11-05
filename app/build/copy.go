package build

import (
	"io"
	"os"
	"path/filepath"

	"github.com/fuxiaohei/pugo-static/app/models"
	"github.com/fuxiaohei/pugo-static/app/utils/mylog"
	"github.com/fuxiaohei/pugo-static/app/vars"
)

// Copy copy static files
func Copy(content *models.Content) bool {
	var (
		copyCount int
		keepCount int
	)
	err := filepath.Walk(vars.Config.MediaDir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		toFile := filepath.Join(vars.Config.DstDir, fpath)
		toFileInfo, _ := os.Stat(toFile)
		if toFileInfo != nil && toFileInfo.ModTime().Unix() == info.ModTime().Unix() {
			mylog.Trace("copy file %s but exist", fpath)
			keepCount++
			content.DstFiles[toFile] = true
			return nil
		}
		if err = copyFile(fpath, toFile); err != nil {
			mylog.Warn("copy file %s from %s error %s", toFile, fpath, err.Error())
			return nil
		}
		copyCount++
		content.DstFiles[toFile] = true
		mylog.Trace("copy file %s from %s", toFile, fpath)
		return nil
	})
	if err != nil {
		return false
	}
	staticDir := content.Theme.StaticDir()
	err = filepath.Walk(staticDir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		relpath, _ := filepath.Rel(staticDir, fpath)
		toFile := filepath.Join(vars.Config.DstDir, relpath)
		toFileInfo, _ := os.Stat(toFile)
		if toFileInfo != nil && toFileInfo.ModTime().Unix() == info.ModTime().Unix() {
			mylog.Trace("copy file %s but exist", fpath)
			keepCount++
			content.DstFiles[toFile] = true
			return nil
		}
		if err = copyFile(fpath, toFile); err != nil {
			mylog.Warn("copy file %s from %s error %s", toFile, fpath, err.Error())
			return nil
		}
		copyCount++
		content.DstFiles[toFile] = true
		mylog.Trace("copy file %s from %s", toFile, fpath)
		return nil
	})
	if err != nil {
		return false
	}
	mylog.Info("copy static files %s, not modified %s", copyCount, keepCount)
	return true
}

func copyFile(from, to string) error {
	if err := os.MkdirAll(filepath.Dir(to), os.ModePerm); err != nil {
		return err
	}
	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	toFile, err := os.OpenFile(to, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer toFile.Close()
	_, err = io.Copy(toFile, fromFile)
	if err != nil {
		return err
	}
	stat, _ := fromFile.Stat()
	return os.Chtimes(to, stat.ModTime(), stat.ModTime()) // set copied file's modification time is same to old file
}
