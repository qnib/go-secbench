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
	c := make(chan interface{}, 1)
	s := secbench.NewStore(c)
	done := make(chan error)
	go s.Loop(done)
	g := secbench.NewGrok()
	sb, err := secbench.NewSecBenc(g, c)
	if err != nil {
		log.Panicf(err.Error())
	}
	sb.Run(cfg)
	c <- "Run complete!"
	err = <- done
	if err != nil {
		log.Fatal(err.Error())
	}
	cMap, _ := cfg.Settings()
	s.Eval(cMap)
}

func main() {
	app := cli.NewApp()
	app.Name = "Run and parse docker security benchmark"
	app.Usage = "go-secbench [options]"
	app.Version = secbench.VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "modes-ignore",
			Value:  "",
			Usage:  "Modes to ignore during evaluation",
		 	EnvVar: "SECB_MODES_IGNORE",
		}, cli.BoolFlag{
			Name:  "skip-pull",
			Usage: "Skip pulling dockerbench image",
			EnvVar: "SECB_SKIP_PULL",
		}, cli.BoolFlag{
			Name:  "quiet",
			Usage: "Supress log output",
			EnvVar: "SECB_QUIET",
		}, cli.BoolFlag{
			Name:  "piped",
			Usage: "Read Benchmark output from sdtin, instead of running a container.",
			EnvVar: "SECB_PIPED",
		}, cli.BoolFlag{
			Name:  "skip-empty-rules",
			Usage: "Skip rules which do not have instances to it.",
			EnvVar: "SECB_SKIP_EMPTY_RULES",
		}, cli.BoolFlag{
			Name:  "fail-on-warn",
			Usage: "Return non-zweo errorcode when warnings are present.",
			EnvVar: "SECB_FAIL_ON_WARN",
		}, cli.StringFlag{
			Name:  "rule-numbers-skip",
			Value: "",
			Usage: "Comma separated list of rule numbers to skip",
			EnvVar: "SECB_RULES_NUMBERS_SKIP",
		}, cli.StringFlag{
			Name:  "rule-numbers-skip-regex",
			Value: "",
			Usage: "Regular expression of rule numbers to skip",
			EnvVar: "SECB_RULES_NUMBERS_SKIP_REGEX",
		}, cli.StringFlag{
			Name:  "rule-numbers-only",
			Value: "",
			Usage: "Comma separated list of rule numbers to include",
			EnvVar: "SECB_RULES_NUMBERS_ONLY",
		}, cli.StringFlag{
			Name:  "rule-numbers-only-regex",
			Value: "",
			Usage: "Regulary expression of rule numbers to include",
			EnvVar: "SECB_RULES_NUMBERS_ONLY_REGEX",
		},
	}
	app.Action = Run
	app.Run(os.Args)
}
