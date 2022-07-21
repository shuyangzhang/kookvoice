package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

const (
	GATEWAY_URL = "https://www.kaiheila.cn/api/v3/gateway/voice?channel_id=%v"
)

var (
	CHANNEL_ID = flag.String("c", "", "channel id")
	TOKEN      = flag.String("t", "", "bot token")
	INPUT      = flag.String("i", "", "input audio")
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
	flag.Parse()
	gatewayUrl := getGatewayUrl()
	startWebsocketClient(gatewayUrl)
}

func getGatewayUrl() string {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(GATEWAY_URL, *CHANNEL_ID),
		nil,
	)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bot %v", *TOKEN))
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
	fmt.Println(">>>> start streaming <<<<")

	cmd := exec.Command(
		"ffmpeg",
		"-re",
		"-loglevel",
		"level+info",
		"-nostats",
		"-i",
		*INPUT,
		"-map",
		"0:a:0",
		"-acodec",
		"libopus",
		"-ab",
		"128k",
		"-filter:a",
		"volume=0.8",
		"-ac",
		"2",
		"-ar",
		"48000",
		"-f",
		"tee",
		fmt.Sprintf("[select=a:f=rtp:ssrc=1357:payload_type=100]rtp://%v:%v?rtcpport=%v", ip, port, rtcpPort),
	)
	cmd.Run()

	fmt.Println("---- start keep recieve message ----")
	keepRecieveMessage(connect)

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

type firstShakeReq struct {
	Request bool              `json:"request"`
	Id      int               `json:"id"`
	Method  string            `json:"method"`
	Data    firstShakeReqData `json:"data"`
}

type firstShakeReqData struct {
}

type secondShakeReq struct {
	Request bool               `json:"request"`
	Id      int                `json:"id"`
	Method  string             `json:"method"`
	Data    secondShakeReqData `json:"data"`
}

type secondShakeReqData struct {
	DisplayName string `json:"displayName"`
}

type BaseShakeReq struct {
	Request bool   `json:"request"`
	Id      int    `json:"id"`
	Method  string `json:"method"`
}

type thirdShakeReq struct {
	Request bool              `json:"request"`
	Id      int               `json:"id"`
	Method  string            `json:"method"`
	Data    thirdShakeReqData `json:"data"`
}

type thirdShakeReqData struct {
	Comedia bool   `json:"comedia"`
	RtcpMux bool   `json:"rtcpMux"`
	Type    string `json:"type"`
}

type thirdShakeResp struct {
	Response bool `json:"response"`
	Id       int  `json:"id"`
	Ok       bool `json:"ok"`
	Data     thirdShakeRespData
}

type thirdShakeRespData struct {
	Id       string `json:"id"`
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	RtcpPort int    `json:"rtcpPort"`
}

type fourthShakeReq struct {
	Request bool               `json:"request"`
	Id      int                `json:"id"`
	Method  string             `json:"method"`
	Data    fourthShakeReqData `json:"data"`
}

type fourthShakeReqData struct {
	AppData       appData       `json:"appData"`
	Kind          string        `json:"kind"`
	PeerId        string        `json:"peerId"`
	RtpParameters rtpParameters `json:"rtpParameters"`
	TransportId   string        `json:"transportId"`
}

type appData struct {
}

type rtpParameters struct {
	Codecs    []codec    `json:"codecs"`
	Encodings []encoding `json:"encodings"`
}

type codec struct {
	Channels    int        `json:"channels"`
	ClockRate   int        `json:"clockRate"`
	MimeType    string     `json:"mimeType"`
	Parameters  parameters `json:"parameters"`
	PayloadType int        `json:"payloadType"`
}

type parameters struct {
	SpropStereo int `json:"sprop-stereo"`
}

type encoding struct {
	Ssrc int `json:"ssrc"`
}
