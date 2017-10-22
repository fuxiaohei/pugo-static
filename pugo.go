package main

import "github.com/fuxiaohei/pugo-static/app/build"
import "github.com/fuxiaohei/pugo-static/app/utils/mylog"

func main() {
	mylog.EnableTrace = true
	build.Build()
}
