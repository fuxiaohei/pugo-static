package main

import (
	"github.com/fuxiaohei/pugo-static/app/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "PuGo.Static"
	app.Usage = "Static Site Generator"
	app.Description = "a very simple static site generator"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		cmd.Build,
	}
	app.RunAndExitOnError()
}
