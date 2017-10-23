package cmd

import (
	"github.com/fuxiaohei/pugo-static/app/build"
	"github.com/fuxiaohei/pugo-static/app/utils/mylog"
	"github.com/urfave/cli"
)

var commonFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "debug",
		Usage: "print debug info",
	},
}

// Build is build command ,
// use to build contents to webpages
var Build = cli.Command{
	Name:      "build",
	ShortName: "b",
	Usage:     "build pugo.static website",
	Flags: append(commonFlags, cli.BoolFlag{
		Name:  "clean",
		Usage: "clean files in destination dir but not compiled by PuGo",
	}),
	Action: func(ctx *cli.Context) error {
		if ctx.Bool("debug") {
			mylog.EnableTrace = true
		}
		build.Build(ctx.Bool("clean"))
		return nil
	},
}
