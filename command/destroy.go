package command

import (
	"encoding/json"
	"errors"
	"pbx/aritask/asterisk"
	"pbx/aritask/cli"
	"strings"
)

/// ///

func Destroy(conf string) *cli.Cmd {
	var cmd *cli.Cmd = new(cli.Cmd)

	cmd.Name("destroy")

	cmd.Usage("destroy [bridge]")
	cmd.Description("Shut down a bridge")
	cmd.AddExample("destroy 2348234")

	cmd.Conf(conf)

	cmd.Exec(
		func(c *cli.Cli, args ...string) (string, error) {
			var conf asterisk.AriApplicationConf
			if len(args) == 0 {
				return "", errors.New("configuration is not defined for this command")
			} else if err := json.Unmarshal([]byte(args[0]), &conf); err != nil {
				return "", err
			}

			if len(args) < 2 {
				return "", errors.New("invalid command usage, enter `help destroy` to get usage information")
			}

			var err error
			var bridge asterisk.Bridge
			if bridge, err = asterisk.BridgeGet(conf, args[1]); err != nil {
				return strings.ToLower(err.Error()), errors.New("failed fetching bridge " + args[1])
			}

			for _, channelId := range bridge.Channels {
				if err := asterisk.ChannelDestroy(conf, channelId); err != nil {
					return "", err
				}
			}

			if err = asterisk.BridgeDestroy(conf, args[1]); err != nil {
				return strings.ToLower(err.Error()), errors.New("failed destroying bridge " + args[1])
			}

			return "Bridge " + args[1] + " shut down", nil
		},
	)

	return cmd
}

/// ///
