package main

import (
	"log"
	"github.com/qnib/go-secbench"
	"os"
	"github.com/codegangsta/cli"
	"github.com/zpatrick/go-config"

)

func Run(ctx *cli.Context) {
	cfg := config.NewConfig([]config.Provider{config.NewCLI(ctx, true)})
	g := secbench.NewGrok()
	sb, err := secbench.NewSecBenc(g)
	if err != nil {
		log.Panicf(err.Error())
	}
	sb.Run(cfg)
}

func main() {
	app := cli.NewApp()
	app.Name = "Run and parse docker security benchmark"
	app.Usage = "go-secbench [options]"
	app.Version = "0.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "modes-show",
			Value: "NOTE,INFO,PASS,WARN",
			Usage: "Modes to show in output",
			EnvVar: "SECB_MODES_SHOW",
		},
	}
	app.Action = Run
	app.Run(os.Args)
}
