package command

import (
	"encoding/json"
	"errors"
	"pbx/aritask/asterisk"
	"pbx/aritask/cli"
	"strings"
	"time"
)

/// ///

func Dial(conf string) *cli.Cmd {
	var cmd *cli.Cmd = new(cli.Cmd)

	cmd.Name("dial")

	cmd.Usage("dial [extensions]...")
	cmd.Description("Initiate a call between endpoints\n(dial needs at least 2 extensions as argument)")
	cmd.AddExample("dial 100 101")

	cmd.Conf(conf)

	cmd.Exec(
		func(cl *cli.Cli, args ...string) (string, error) {

			var conf asterisk.AriApplicationConf
			if len(args) == 0 {
				return "", errors.New("configuration is not defined for this command")
			} else if err := json.Unmarshal([]byte(args[0]), &conf); err != nil {
				return "", err
			}

			if len(args[1:]) < 2 {
				return "", errors.New("invalid command usage, enter `help dial` to get usage information")
			}

			var channelsIds []string
			var bridge asterisk.Bridge
			var err error

			for _, endpoint := range args[1:] {
				channel, err := asterisk.ChannelCreate(conf, endpoint)
				if err != nil {
					return "", errors.New("issue creating channel for endpoint '" + endpoint + "' {" + err.Error() + "}")
				}

				if channel.ID != "" {
					channelsIds = append(channelsIds, channel.ID)
					cl.Print("Created channel " + channel.ID + " for endpoint " + endpoint)
				} else {
					cl.Print("!Failed to create channel for endpoint " + endpoint)
				}

			}

			time.Sleep(500 * time.Millisecond)

			if len(channelsIds) == 0 {
				return "Failed", errors.New("endpoints unreachable, stop initiating call")
			}
			if len(channelsIds) == 1 {
				if err := asterisk.ChannelDestroy(conf, channelsIds[0]); err != nil {
					return "", err
				}
				return "Failed", errors.New("only channel " + channelsIds[0] + " reachable, stop initiating call")
			}

			var bridgeType string = asterisk.BRIDGE_CALL
			if len(channelsIds) > 2 {
				bridgeType = asterisk.BRIDGE_CONFERENCE
			}

			if bridge, err = asterisk.BridgeCreate(conf, bridgeType); err != nil {
				return "", errors.New("issue creating bridge -" + bridgeType + " {" + err.Error() + "}")
			}

			cl.Print("Bridge " + bridge.ID + " created for connecting channels: " + strings.Join(channelsIds, ", "))

			if err = asterisk.BridgeAddChannel(conf, bridge.ID, channelsIds); err != nil {
				for _, channelId := range channelsIds {
					asterisk.ChannelDestroy(conf, channelId)
				}

				asterisk.BridgeDestroy(conf, bridge.ID)

				return "", errors.New("failed adding channels (" + strings.Join(channelsIds, ",") + ") to bridge " + bridge.ID + "{" + err.Error() + "}")
			}

			cl.Print("Success adding channels to bridge " + bridge.ID)

			for _, id := range channelsIds {
				if err = asterisk.ChannelDial(conf, id); err != nil {
					return "", errors.New("issue with request dialing channel " + id + "{" + err.Error() + "}")
				}
				cl.Print("Dialed channel " + id)
			}

			cl.Print("Dial complete")

			return "Complete", nil
		},
	)

	return cmd
}

/// ///
