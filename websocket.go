package kookvoice

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
)

func getGatewayUrl(token string, channelId string) string {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(GATEWAY_URL, channelId),
		nil,
	)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bot %v", token))
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
	// fmt.Printf("gateway url is %v \n", gatewayUrl)
	return gatewayUrl
}

func initWebsocketClient(websocketHost string) (*websocket.Conn, string) {
	dialer := websocket.Dialer{}
	connect, _, err := dialer.Dial(websocketHost, nil)
	if err != nil {
		panic(err)
	}

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

	// fmt.Println("---- start recieve first shake message ----")
	recieveMessageOnce(connect)
	// fmt.Println("---- end recieve first shake message ----")

	secondShakeReqObj := secondShakeReq{
		Request: true,
		Id:      1000000,
		Method:  "join",
		Data: secondShakeReqData{
			DisplayName: "",
		},
	}

	secondShakeReqStr, err := json.Marshal(secondShakeReqObj)

	if err != nil {
		panic(err)
	}

	err = connect.WriteMessage(
		websocket.TextMessage,
		[]byte(secondShakeReqStr),
	)
	if err != nil {
		panic(err)
	}

	// fmt.Println("---- start recieve second shake message ----")
	recieveMessageOnce(connect)
	// fmt.Println("---- end recieve second shake message ----")

	thirdShakeReqObj := thirdShakeReq{
		Request: true,
		Id:      1000000,
		Method:  "createPlainTransport",
		Data: thirdShakeReqData{
			Comedia: true,
			RtcpMux: false,
			Type:    "plain",
		},
	}

	thirdShakeReqStr, err := json.Marshal(thirdShakeReqObj)

	if err != nil {
		panic(err)
	}

	err = connect.WriteMessage(
		websocket.TextMessage,
		[]byte(thirdShakeReqStr),
	)
	if err != nil {
		panic(err)
	}

	// fmt.Println("---- start recieve third shake message ----")
	thirdShakeRespStr := recieveMessageOnce(connect)
	// fmt.Println("---- end recieve third shake message ----")
	var thirdShakeRespObj thirdShakeResp
	err = json.Unmarshal(thirdShakeRespStr, &thirdShakeRespObj)
	if err != nil {
		panic(err)
	}
	transportId := thirdShakeRespObj.Data.Id
	ip := thirdShakeRespObj.Data.Ip
	port := thirdShakeRespObj.Data.Port
	rtcpPort := thirdShakeRespObj.Data.RtcpPort

	fourthShakeReqObj := fourthShakeReq{
		Request: true,
		Id:      1000000,
		Method:  "produce",
		Data: fourthShakeReqData{
			AppData: appData{},
			Kind:    "audio",
			PeerId:  "",
			RtpParameters: rtpParameters{
				Codecs: []codec{
					{
						Channels:  2,
						ClockRate: 48000,
						MimeType:  "audio/opus",
						Parameters: parameters{
							SpropStereo: 1,
						},
						PayloadType: 100,
					},
				},
				Encodings: []encoding{
					{
						Ssrc: 1357,
					},
				},
			},
			TransportId: transportId,
		},
	}

	fourthShakeReqStr, err := json.Marshal(fourthShakeReqObj)
	if err != nil {
		panic(err)
	}

	err = connect.WriteMessage(
		websocket.TextMessage,
		[]byte(fourthShakeReqStr),
	)
	if err != nil {
		panic(err)
	}

	// fmt.Println("---- start recieve fourth shake message ----")
	recieveMessageOnce(connect)
	// fmt.Println("---- end recieve fourth shake message ----")

	fmt.Println(">>>> shake hands succeed <<<<")
	//fmt.Printf("ssrc=1357 ffmpeg rtp url: \n  rtp://%v:%v?rtcpport=%v \n", ip, port, rtcpPort)

	rtpUrl := fmt.Sprintf("rtp://%v:%v?rtcpport=%v", ip, port, rtcpPort)

	return connect, rtpUrl
}

func recieveMessageOnce(connect *websocket.Conn) []byte {
	messageType, messageData, err := connect.ReadMessage()
	if err != nil {
		fmt.Println("failed to recieve message once")
		return nil
	}
	switch messageType {
	case websocket.TextMessage:
		//fmt.Println(string(messageData))
	case websocket.BinaryMessage:
		//fmt.Println(messageData)
	case websocket.CloseMessage:
		fmt.Println("recieved close message")
	case websocket.PingMessage:
		fmt.Println("recieved ping message")
	case websocket.PongMessage:
		fmt.Println("recieved pong message")
	default:
		fmt.Println("recieved unknown message")
	}
	return messageData
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
