package cmd

import (
	"github.com/fuxiaohei/pugo-static/app/build"
	"github.com/fuxiaohei/pugo-static/app/utils/mylog"
	"github.com/urfave/cli"
)

var Server = cli.Command{
	Name:      "server",
	Usage:     "run http server to serve compiled files",
	ShortName: "srv",
	Flags: append(commonFlags, cli.BoolFlag{
		Name:  "static",
		Usage: "only serve static directory, do not build",
	}, cli.BoolFlag{
		Name:  "clean",
		Usage: "clean files in destination dir but not compiled by PuGo",
	}, cli.StringFlag{
		Name:  "addr",
		Usage: "run http server on this address",
		Value: "0.0.0.0:9899",
	}),
	Action: func(ctx *cli.Context) error {
		if !ctx.Bool("static") {
			if ctx.Bool("debug") {
				mylog.EnableTrace = true
			}
			if !build.Build(ctx.Bool("clean")) {
				return nil
			}
		}
		build.Server(ctx.String("addr"))
		return nil
	},
}
