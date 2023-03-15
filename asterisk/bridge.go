package asterisk

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

/// ///

type BaseBridgeId string

const (
	BRIDGE_CALL       = "Call"
	BRIDGE_CONFERENCE = "Conference"
)

/// ///

type Bridge struct {
	ID         string   `json:"id"`
	Technology string   `json:"technology"`
	Type       string   `json:"bridge_type"`
	Class      string   `json:"bridge_class"`
	Creator    string   `json:"creator"`
	Name       string   `json:"name"`
	Channels   []string `json:"channels"`
}

/// ///

func BridgeList(conf AriApplicationConf) ([]Bridge, error) {
	var bridges []Bridge
	var body string
	var err error

	if body, _, err = AriRequest(
		conf,
		http.MethodGet,
		"/bridges",
		nil,
	); err != nil {
		return bridges, err
	}

	if err := json.Unmarshal([]byte(body), &bridges); err != nil {
		return bridges, err
	}

	return bridges, nil
}

/// ///

func BridgeGet(conf AriApplicationConf, bridgeId string) (Bridge, error) {
	var bridge Bridge
	var err error
	var body string
	var statusCode int
	var result map[string]string

	if body, statusCode, err = AriRequest(
		conf,
		http.MethodGet,
		strings.ReplaceAll("/bridges/{bridgeId}", "{bridgeId}", bridgeId),
		nil,
	); err != nil {
		return bridge, err
	}

	switch statusCode {
	case 404:
		if err = json.Unmarshal([]byte(body), &result); err == nil {
			err = errors.New(result["message"])
		}
	default:
		err = json.Unmarshal([]byte(body), &bridge)
	}

	return bridge, err
}

/// ///

func BridgeDestroy(conf AriApplicationConf, bridgeId string) error {
	var err error
	var body string
	var statusCode int
	var result map[string]string

	if body, statusCode, err = AriRequest(
		conf,
		http.MethodDelete,
		strings.ReplaceAll("/bridges/{bridgeId}", "{bridgeId}", bridgeId),
		nil,
	); err != nil {
		return err
	}

	switch statusCode {
	case 404:
		if err = json.Unmarshal([]byte(body), &result); err == nil {
			err = errors.New(result["message"])
		}
	default:
		err = nil
	}

	return err
}

/// ///

func BridgeUpdate(conf AriApplicationConf, bridgeId string, query map[string]string) (Bridge, error) {
	var bridge Bridge
	var err error
	var body string
	var statusCode int

	if body, statusCode, err = AriRequest(
		conf,
		http.MethodPost,
		strings.ReplaceAll("/bridges/{bridgeId}", "{bridgeId}", bridgeId),
		query,
	); err != nil {
		return bridge, err
	}

	fmt.Println(statusCode)
	fmt.Println(body)

	err = json.Unmarshal([]byte(body), &bridge)

	return bridge, err
}

/// ///

func BridgeCreate(conf AriApplicationConf, bridgeType string) (Bridge, error) {
	var bridge Bridge
	var randId string

	for {
		randId = fmt.Sprint(rand.Uint32())
		if _, ok := Cache.Bridges[randId]; !ok {
			break
		}
	}

	var body string
	var err error

	if body, _, err = AriRequest(
		conf,
		http.MethodPost,
		"/bridges",
		map[string]string{
			"type":     "mixing",
			"bridgeId": randId,
			"name":     randId,
		},
	); err != nil {
		return bridge, err
	}

	if err := json.Unmarshal([]byte(body), &bridge); err != nil {
		return bridge, err
	}

	Cache.Bridges[bridge.ID] = bridgeType

	return bridge, nil
}

/// ///

func bridgeChannel(conf AriApplicationConf, method string, bridgeId string, channels []string) error {
	var err error
	var result map[string]string
	var statusCode int
	var body string

	if body, statusCode, err = AriRequest(
		conf,
		method,
		strings.ReplaceAll("/bridges/{bridgeId}/addChannel", "{bridgeId}", bridgeId),
		map[string]string{
			"channel": strings.Join(channels, ","),
		},
	); err != nil {
		return err
	}

	switch statusCode {
	case 400, 404, 409, 422:
		if err = json.Unmarshal([]byte(body), &result); err == nil {
			err = errors.New(result["message"])
		}
	default:
		err = nil
	}

	return err
}

/// ///

func BridgeAddChannel(conf AriApplicationConf, bridgeId string, channels []string) error {
	return bridgeChannel(
		conf,
		http.MethodPost,
		bridgeId,
		channels,
	)
}

/// ///

func BridgeRemoveChannel(conf AriApplicationConf, bridgeId string, channels []string) error {
	return bridgeChannel(
		conf,
		http.MethodDelete,
		bridgeId,
		channels,
	)
}

/// ///
