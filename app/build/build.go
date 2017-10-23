package build

import (
	"time"

	"github.com/fuxiaohei/pugo-static/app/models"
	"github.com/fuxiaohei/pugo-static/app/utils/mylog"
)

// Build run build process
func Build(isClean bool) {
	var (
		content = models.NewContent()
		st      = time.Now()
	)
	// read all contents and generate all cached temporary data
	if !Read(content) {
		return
	}
	mylog.Trace("read all contents done, %s ms", time.Since(st).Nanoseconds()/1e6)
	st = time.Now()

	if !Compile(content) {
		return
	}
	mylog.Trace("compile all webpages done, %s ms", time.Since(st).Nanoseconds()/1e6)
	st = time.Now()

	if !Copy(content) {
		return
	}
	mylog.Trace("copy static files done, %s ms", time.Since(st).Nanoseconds()/1e6)
	st = time.Now()
	if isClean {
		Clean(content)
		mylog.Trace("clean destination directory, %s ms", time.Since(st).Nanoseconds()/1e6)
	}
}
