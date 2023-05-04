package main

import (
	"encoding/json"
	"flag"
	"pbx/aritask/asterisk"
	"pbx/aritask/cli"
	"pbx/aritask/command"
)

/// ///

const (
	ALL_APP       = "aritest"
	ALL_HOST      = ""
	ALL_PORT      = 8088
	ALL_USER      = "asterisk"
	ALL_SECRET    = "vrjwbpviudkvbwoeibvpiufbsfdvlkjqepr"
	ALL_BASE_PATH = "/ari"
)

/// ///

func main() {
	var verbose bool

	flag.BoolVar(&verbose, "v", false, "enable verbosity")
	flag.Parse()

	var conf asterisk.AriApplicationConf = asterisk.AriApplicationConf{
		Application: ALL_APP,
		Host:        ALL_HOST,
		Port:        ALL_PORT,
		Username:    ALL_USER,
		Password:    ALL_SECRET,
		BasePath:    ALL_BASE_PATH,
	}

	var err error

	jsonrawconf, err := json.Marshal(conf)

	if err != nil {
		panic("JSON configuration:" + err.Error())
	}

	var appConf string = string(jsonrawconf)

	var CLI *cli.Cli = new(cli.Cli)
	CLI.
		AddCommand(
			command.Dial(appConf),
		).
		AddCommand(
			command.List(appConf),
		).
		AddCommand(
			command.Join(appConf),
		).
		AddCommand(
			command.Destroy(appConf),
		).
		AddCommand(
			cli.Help(),
		).
		AddCommand(
			cli.Exit(),
		).
		AddBackgroundCommand(
			command.ApplicationListen(appConf),
		).
		Title(
			"CLI-ARI",
		).
		Verbose(verbose)

	cli.Run(CLI)
}

/// just testing github codespaces

/// ///
