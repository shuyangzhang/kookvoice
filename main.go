package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	GATEWAY_URL = "https://www.kaiheila.cn/api/v3/gateway/voice?channel_id=%v"
	CHANNEL_ID  = "7553950861219535"
	TOKEN       = ""
)

type gatewayResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    gatewayRespData `json:"data"`
}

type gatewayRespData struct {
	GatewayUrl string `json:"gateway_url"`
}

func main() {
	gatewayUrl := getGatewayUrl()
	startWebsocketClient(gatewayUrl)
}

func getGatewayUrl() string {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(GATEWAY_URL, CHANNEL_ID),
		nil,
	)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bot %v", TOKEN))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	respObj := gatewayResp{}
	if err := json.Unmarshal(respBytes, &respObj); err != nil {
		panic(err)
	}
	gatewayUrl := respObj.Data.GatewayUrl
	fmt.Printf("gateway url is %v \n", gatewayUrl)
	return gatewayUrl
}

func startWebsocketClient(websocketHost string) {
	dialer := websocket.Dialer{}
	connect, _, err := dialer.Dial(websocketHost, nil)
	if err != nil {
		panic(err)
	}
	defer connect.Close()

	go keepWebsocketClientAlive(connect)

	firstShakeReqObj := firstShakeReq{
		Request: true,
		Id:      1000000,
		Method:  "getRouterRtpCapabilities",
	}

	firstShakeReqStr, err := json.Marshal(firstShakeReqObj)

	if err != nil {
		panic(err)
	}

	err = connect.WriteMessage(
		websocket.TextMessage,
		[]byte(firstShakeReqStr),
	)
	if err != nil {
		panic(err)
	}

	keepRecieveMessage(connect)

}

func keepRecieveMessage(connect *websocket.Conn) {
	for {
		messageType, messageData, err := connect.ReadMessage()
		if err != nil {
			fmt.Println("failed to read message")
			break
		}
		switch messageType {
		case websocket.TextMessage:
			fmt.Println(string(messageData))
		case websocket.BinaryMessage:
			fmt.Println(messageData)
		case websocket.CloseMessage:
			fmt.Println("recieved close message")
		case websocket.PingMessage:
			fmt.Println("recieved ping message")
		case websocket.PongMessage:
			fmt.Println("recieved pong message")
		default:
			fmt.Println("recieved unknown message")
		}

	}
}

func keepWebsocketClientAlive(connect *websocket.Conn) {
	for {
		time.Sleep(30 * time.Second)
		err := connect.WriteMessage(
			websocket.PingMessage,
			[]byte{},
		)
		if err != nil {
			fmt.Println("heart beat failed")
			break
		}
		fmt.Println("heart beat succeed, ws client is alive")
	}
}

type firstShakeReq struct {
	Request bool              `json:"request"`
	Id      int               `json:"id"`
	Method  string            `json:"method"`
	Data    firstShakeReqData `json:"data"`
}

type firstShakeReqData struct {
}
