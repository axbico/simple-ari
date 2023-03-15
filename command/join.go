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

func Join(conf string) *cli.Cmd {
	var cmd *cli.Cmd = new(cli.Cmd)

	cmd.Name("join")

	cmd.Conf(conf)

	cmd.Exec(
		func(cl *cli.Cli, args ...string) (string, error) {
			var conf asterisk.AriApplicationConf
			if len(args) > 0 {
				if err := json.Unmarshal([]byte(args[0]), &conf); err != nil {
					return "", err
				}
			} else {
				return "", errors.New("configuration is not defined for this command")
			}

			args = args[1:]

			if len(args) < 2 {
				return "", errors.New("invalid command usage, enter `help join` to get usage information")
			}

			var bridgeId string = args[0]

			args = args[1:]

			var endpoints map[string]byte = make(map[string]byte)
			for _, arg := range args {
				endpoints[arg] = 1
			}

			var channelIds []string = []string{}
			var err error
			var bridge asterisk.Bridge

			if bridge, err = asterisk.BridgeGet(conf, bridgeId); err != nil {
				return "", errors.New(strings.ToLower(err.Error()))
			}

			for endpoint := range endpoints {
				for _, channel := range bridge.Channels {
					if endpoint == strings.Split(channel, "-")[1] {
						delete(endpoints, endpoint)
						cl.Print("Endpoint " + endpoint + " already in this call")
					}
				}
			}

			for endpoint := range endpoints {
				channel, err := asterisk.ChannelCreate(conf, endpoint)
				if err != nil {
					return "", errors.New("issue creating channel for endpoint '" + endpoint + "' {" + err.Error() + "}")
				}

				if channel.ID != "" {
					channelIds = append(channelIds, channel.ID)
					cl.Print("Created channel " + channel.ID + " for endpoint " + endpoint)
				} else {
					cl.Print("!Failed to create channel for endpoint " + endpoint)
				}
			}

			time.Sleep(500 * time.Millisecond)

			if len(channelIds) == 0 {
				return "No channels to add", nil
			}

			if asterisk.Cache.Bridges[bridgeId] == asterisk.BRIDGE_CALL {
				asterisk.Cache.Bridges[bridgeId] = asterisk.BRIDGE_CONFERENCE
				cl.Print("Bridge " + bridgeId + " promoted to conference")
			}

			if err = asterisk.BridgeAddChannel(conf, bridgeId, channelIds); err != nil {
				for _, channelId := range channelIds {
					asterisk.ChannelDestroy(conf, channelId)
				}

				return strings.ToLower(err.Error()), errors.New("failed adding channels (" + strings.Join(channelIds, ", ") + ") to bridge " + bridgeId)
			}

			cl.Print("Success adding channels to bridge " + bridgeId)

			for _, id := range channelIds {
				if err = asterisk.ChannelDial(conf, id); err != nil {
					return strings.ToLower(err.Error()), errors.New("issue with request dialing channel " + id)
				}
				cl.Print("Dialed channel " + id)
			}

			return "Complete", nil
		},
	)

	return cmd
}

/// ///
