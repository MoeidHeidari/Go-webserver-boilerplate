package ws

import (
	"fmt"
	"main/lib"
	"net/http"
	"time"

	"main/api/kubes"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Ws struct {
	upgrader websocket.Upgrader
	logger   lib.Logger
}

func NewWs(logger lib.Logger) Ws {
	wsupgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return Ws{
		upgrader: wsupgrader,
		logger:   logger,
	}
}

func (w Ws) MessageHandler(c *gin.Context) {
	connection, err := w.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		w.logger.Fatal(err.Error())
	}
	defer connection.Close()
	_, pod_name, err := connection.ReadMessage()
	if err != nil {
		w.logger.Fatal(err.Error())
	}

	if err != nil {
		w.logger.Fatal(err.Error())
	}
	k := kubes.NewKubeRequest(w.logger)
	for {

		response := k.GetCurrentPodStatusRequest(string(pod_name))
		fmt.Println(string(response))
		err = connection.WriteMessage(websocket.TextMessage, response)
		time.Sleep(time.Second * 10)
		if err != nil {
			w.logger.Fatal(err.Error())
			break
		}
		defer connection.Close()
	}

}
