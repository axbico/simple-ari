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

type Channel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	State  string `json:"state"`
	Caller struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	} `json:"caller"`
}

/// ///

func ChannelList(conf AriApplicationConf) ([]Channel, error) {
	var channels []Channel
	var body string
	var err error

	if body, _, err = AriRequest(
		conf,
		http.MethodGet,
		"/channels",
		nil,
	); err != nil {
		return channels, err
	}

	if err := json.Unmarshal([]byte(body), &channels); err != nil {
		return channels, err
	}

	return channels, nil
}

/// ///

func ChannelGet(conf AriApplicationConf, channelId string) (Channel, error) {
	var channel Channel
	var err error
	var body string
	var statusCode int
	var result map[string]string

	if body, statusCode, err = AriRequest(
		conf,
		http.MethodGet,
		strings.ReplaceAll("/channels/{channelId}", "{channelId}", channelId),
		nil,
	); err != nil {
		return channel, err
	}

	switch statusCode {
	case 404:
		if err = json.Unmarshal([]byte(body), &result); err == nil {
			err = errors.New(result["message"])
		}
	default:
		err = json.Unmarshal([]byte(body), &channel)
	}

	return channel, err
}

/// ///

func ChannelCreate(conf AriApplicationConf, endpoint string) (Channel, error) {
	var channel Channel
	var err error
	var body string
	var statusCode int
	var result map[string]string

	if body, statusCode, err = AriRequest(
		conf,
		http.MethodPost,
		"/channels/create",
		map[string]string{
			"endpoint":  "PJSIP/" + endpoint,
			"channelId": fmt.Sprint(rand.Uint32()) + "-" + endpoint,
		},
	); err != nil {
		return channel, err
	}

	switch statusCode {
	case 409:
		if err = json.Unmarshal([]byte(body), &result); err == nil {
			err = errors.New(result["message"])
		}
	default:
		err = json.Unmarshal([]byte(body), &channel)
	}

	return channel, err
}

/// ///

func ChannelDestroy(conf AriApplicationConf, channelId string) error {
	var err error
	var body string
	var statusCode int
	var result map[string]string

	if body, statusCode, err = AriRequest(
		conf,
		http.MethodDelete,
		strings.ReplaceAll("/channels/{channelId}", "{channelId}", channelId),
		nil,
	); err != nil {
		return err
	}

	switch statusCode {
	case 400, 404:
		if err = json.Unmarshal([]byte(body), &result); err == nil {
			err = errors.New(result["message"])
		}
	default:
		err = nil
	}

	return err
}

/// ///

func ChannelDial(conf AriApplicationConf, channelId string) error {
	var err error
	var body string
	var statusCode int
	var result map[string]string

	body, statusCode, err = AriRequest(
		conf,
		http.MethodPost,
		strings.ReplaceAll("/channels/{channelId}/dial", "{channelId}", channelId),
		map[string]string{
			"timeout": "20",
		},
	)

	switch statusCode {
	case 404, 409:
		if err = json.Unmarshal([]byte(body), &result); err == nil {
			err = errors.New(result["message"])
		}
	}

	return err
}

/// ///
