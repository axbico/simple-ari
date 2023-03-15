package command

import (
	"encoding/json"
	"errors"
	"pbx/aritask/asterisk"
	"pbx/aritask/cli"
	"strings"
)

/// ///

func List(conf string) *cli.Cmd {
	var cmd *cli.Cmd = new(cli.Cmd)

	cmd.Name("list")

	cmd.Conf(conf)

	cmd.Exec(list)

	return cmd
}

/// ///

func list(cl *cli.Cli, args ...string) (string, error) {

	if len(args) == 0 {
		return "", errors.New("configuration is not defined for this command")
	}
	var conf asterisk.AriApplicationConf
	if len(args) == 0 {
		return "", errors.New("configuration is not defined for this command")
	} else if err := json.Unmarshal([]byte(args[0]), &conf); err != nil {
		return "", err
	}

	var bridges []asterisk.Bridge
	var err error
	var display string = ""

	if bridges, err = asterisk.BridgeList(conf); err != nil {
		return "", errors.New(strings.ToLower(err.Error()))
	}

	if len(bridges) == 0 {
		display += "No active bridges"
	} else {
		for _, bridge := range bridges {
			display += "  -\n   id: " + bridge.ID + "\n   participants: ( "
			for _, channel := range bridge.Channels {
				display += strings.Split(channel, "-")[1] + " "
			}
			display += ")\n"
		}
	}

	cl.Print(display)

	return "", nil
}

/// ///
