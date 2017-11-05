package build

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/fuxiaohei/pugo-static/app/utils/mylog"
	"github.com/fuxiaohei/pugo-static/app/vars"
)

// Server start http server with addr
func Server(addr string) {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		for _, file := range toServerFile(p) {
			if _, err := os.Stat(file); err == nil {
				http.ServeFile(rw, r, file)
				mylog.Trace("http serve file %s, %s", file, r.RemoteAddr)
				return
			}
		}
		http.NotFound(rw, r)
	})
	mylog.Info("http listen %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		mylog.Error("http listen %s error %s", addr, err.Error())
		return
	}
}

func toServerFile(p string) []string {
	var list []string
	list = append(list, filepath.ToSlash(path.Join(vars.Config.DstDir, p)))
	if filepath.Ext(p) == "" {
		list = append(list, filepath.ToSlash(path.Join(vars.Config.DstDir, p+".html")))
	}
	return list
}
