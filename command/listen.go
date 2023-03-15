package command

import (
	"encoding/json"
	"errors"
	"pbx/aritask/asterisk"
	"pbx/aritask/cli"
)

/// ///

func ApplicationListen(conf string) *cli.Cmd {
	var cmd *cli.Cmd = new(cli.Cmd)

	cmd.Exec(listen)

	cmd.Conf(conf)

	return cmd
}

/// ///

func listen(cl *cli.Cli, args ...string) (string, error) {
	var conf asterisk.AriApplicationConf
	if err := json.Unmarshal([]byte(args[0]), &conf); err != nil {
		return "", errors.New("issue with json event listener configuration")
	}

	readChannel, terminateChannel := asterisk.EventListener(conf)

	go func(cl *cli.Cli, rec chan []byte, conf asterisk.AriApplicationConf) {
		for {
			go func(body []byte) {
				var event asterisk.Event
				if err := json.Unmarshal(body, &event); err != nil {
					cl.Print("EventListener issue extracting json body")
					return
				}
				if handler, ok := asterisk.EventHandle[event.Type]; ok {
					if result := handler(event, conf); result != "" {
						cl.Print("<event-message> " + result)
					}
				}
			}(<-rec)
		}
	}(cl, readChannel, conf)

	<-terminateChannel

	return "", nil
}

/// ///

func printEvent(cl *cli.Cli, event asterisk.Event) {
	var display string = "Event<" + string(event.Type) + ">"
	display += "bridge[" + event.Bridge.ID + "] "
	display += "channel[" + event.Channel.ID + "]"

	if event.Type == asterisk.CHANNEL_LEFT_BRIDGE {
		display += asterisk.CHANNEL_LEFT_BRIDGE
	}

	cl.Print(display)
}

/// ///
