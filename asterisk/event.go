package asterisk

/// ///

const (
	CHANNEL_DESTROYED   = "ChannelDestroyed"
	CHANNEL_LEFT_BRIDGE = "ChannelLeftBridge"
)

/// ///

type Event struct {
	Type    string  `json:"type"`
	Channel Channel `json:"channel"`
	Bridge  Bridge  `json:"bridge"`
}

/// ///

var EventHandle map[string]func(Event, AriApplicationConf) string = map[string]func(Event, AriApplicationConf) string{
	CHANNEL_DESTROYED:   channelDestroyed,
	CHANNEL_LEFT_BRIDGE: channelLeftBridge,
}

/// ///

func channelDestroyed(event Event, conf AriApplicationConf) string {
	return ""
}

/// ///

func channelLeftBridge(event Event, conf AriApplicationConf) string {
	if len(event.Bridge.Channels) == 1 {

		if Cache.Bridges[event.Bridge.ID] == BRIDGE_CALL {
			if err := BridgeRemoveChannel(conf, event.Bridge.ID, event.Bridge.Channels); err != nil {
				return err.Error()
			}
			if err := ChannelDestroy(conf, event.Bridge.Channels[0]); err != nil {
				return err.Error()
			}
			if err := BridgeDestroy(conf, event.Bridge.ID); err != nil {
				return err.Error()
			}

			return "Call ended"
		} else {
			return "Play MOH, only one participant in conference"
		}
	}

	if len(event.Bridge.Channels) == 0 {
		if err := BridgeDestroy(conf, event.Bridge.ID); err != nil {
			return err.Error()
		}

		return "Conference ended"
	}

	return ""
}

/// ///
