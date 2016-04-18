package client

import (
	"strings"
	"sync"
	"ultraman/lib/log"

	httpclient "ultraman/lib/conn/http"
	websocketclient "ultraman/lib/conn/websocket"

	"github.com/gorilla/websocket"
)

var wg sync.WaitGroup

func StartClient(authKey, webSocketAddr string) {
	wg.Add(1)
	buildWebSocketClient(&wg, authKey, webSocketAddr)
	wg.Wait()
}

func buildWebSocketClient(wg *(sync.WaitGroup), auth, addr string) {

	defer (*wg).Done()

	wsClient := &websocketclient.ClientClient{
		Addr:    addr,
		AuthKey: auth,
	}

	httpClient := &httpclient.ClientClient{}

	if wsClient.Dial() == true {
		if wsClient.Auth() == true {
			for {
				mt, msg, err := wsClient.Conn.ReadMessage()
				if err != nil {
					log.Warn("Failed to read websocket: %v", err)
					return
				}
				log.Debug("Websocket recv: %v,%s", mt, msg)
				// DATA
				if mt == websocket.BinaryMessage {
					headers := strings.Split(string(msg), "\n")
					id := headers[0]
					requestHeaders := headers[1:]
					message := []byte(strings.Join(requestHeaders, "\n"))

					resp := httpClient.OpenUrl(&message)
					newResp := string(id) + "\n" + string(resp)

					log.Debug("%s", newResp)

					err = wsClient.Conn.WriteMessage(mt, []byte(newResp))
					if err != nil {
						log.Warn("Failed to write websocket: %v", err)
						break
					}
				}
			}
		} else {
			log.Error("Auth failed")
		}
	} else {
		log.Error("Failed to connect server")
	}
}
