package asterisk

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pbx/aritask/websocket"
	"time"
)

/// ///

type __miniCache struct {
	Bridges map[string]string // map[bridgeId]callType(Call|Conference)
}

var Cache *__miniCache = func() *__miniCache {
	var mc *__miniCache = new(__miniCache)
	mc.Bridges = make(map[string]string)
	return mc
}()

/// ///

const (
	SCHEME  = "http"
	API_KEY = "api_key"
	APP     = "app"
)

const (
	EVENT_LISTEN_PATH = "/events"
)

/// ///

type AriApplicationConf struct {
	Application string `json:"app"`
	Host        string `json:"host"`
	Username    string `json:"user"`
	Password    string `json:"secret"`
	Port        int    `json:"port"`
	BasePath    string `json:"path"`
	ws          websocket.Websocket
}

/// ///

func EventListener(conf AriApplicationConf) (chan []byte, chan bool) {

	conf.ws = *websocket.Dial(conf.Host+":"+fmt.Sprint(conf.Port), conf.BasePath+EVENT_LISTEN_PATH, map[string]string{
		API_KEY: conf.Username + ":" + conf.Password,
		APP:     conf.Application,
	})

	return conf.ws.WaitReaderChannel(), conf.ws.WaitTerminatedConnectionChannel()
}

/// ///

func AriRequest(conf AriApplicationConf, method string, path string, values map[string]string) (string, int, error) {
	var baseUrl url.URL
	baseUrl.Scheme = SCHEME
	baseUrl.Host = conf.Host + ":" + fmt.Sprint(conf.Port)
	baseUrl.Path = conf.BasePath + path

	vars := ""
	if v, set := values["variables"]; set {
		vars = v
		delete(values, "variables")
	}

	params := url.Values{}
	for k, v := range values {
		params.Add(k, v)
	}

	var body io.Reader = nil
	if vars != "" {
		body = bytes.NewBuffer([]byte(vars))
	}

	params.Add(APP, conf.Application)

	baseUrl.RawQuery = params.Encode()

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, baseUrl.String(), body)
	if err != nil {
		return "", 0, err
	}

	req.SetBasicAuth(conf.Username, conf.Password)

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return string(rawBody), resp.StatusCode, nil
}

/// ///
