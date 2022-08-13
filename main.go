package kookvoice

func Play(token string, channelId string, input string) {
	gatewayUrl := getGatewayUrl(token, channelId)
	connect, rtpUrl := initWebsocketClient(gatewayUrl)
	defer connect.Close()
	go keepWebsocketClientAlive(connect)
	go keepRecieveMessage(connect)
	streamAudio(rtpUrl, input)
}
