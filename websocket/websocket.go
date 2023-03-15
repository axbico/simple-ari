package websocket

import (
	"net/url"

	"github.com/gorilla/websocket"
)

/// ///

type Websocket struct {
	connection                      *websocket.Conn
	waitTerminatedConnectionChannel chan bool
	readerChannel                   chan []byte
}

/// ///

func Dial(host string, path string, query map[string]string) *Websocket {
	var ws *Websocket = new(Websocket)

	baseUrl, _ := url.Parse("ws://" + host)
	baseUrl.Path = path

	params := url.Values{}
	for k, v := range query {
		params.Add(k, v)
	}

	baseUrl.RawQuery = params.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(baseUrl.String(), nil)

	if err != nil {
		panic(err.Error())
	}

	ws.connection = conn
	ws.waitTerminatedConnectionChannel = make(chan bool)
	ws.readerChannel = make(chan []byte)

	ws.listen()

	return ws
}

/// ///

func (ws *Websocket) CloseConnection() {
	ws.connection.Close()
	ws.waitTerminatedConnectionChannel <- true
}

/// ///

func (ws *Websocket) Send(output []byte) {
	if err := ws.connection.WriteMessage(websocket.TextMessage, output); err != nil {
		ws.CloseConnection()
	}
}

/// ///

func (ws *Websocket) listen() {
	go func() {
		for {
			if x, body, _ := ws.connection.ReadMessage(); x == -1 {
				ws.CloseConnection()
			} else {
				ws.readerChannel <- body
			}
		}
	}()
}

/// ///

func (ws *Websocket) WaitTerminatedConnectionChannel() chan bool {
	return ws.waitTerminatedConnectionChannel
}

func (ws *Websocket) WaitReaderChannel() chan []byte {
	return ws.readerChannel
}

/// ///
