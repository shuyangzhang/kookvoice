package kookvoice

func Play(token string, channelId string, input string) {
	gatewayUrl := getGatewayUrl(token, channelId)
	connect, rtpUrl := initWebsocketClient(gatewayUrl)
	defer connect.Close()
	go keepWebsocketClientAlive(connect)
	go keepRecieveMessage(connect)
	streamAudio(rtpUrl, input)
}

func New(token string, channelId string) (*voiceInstance, error) {
	vi := voiceInstance{
		Token:     token,
		ChannelId: channelId,
	}
	err := vi.Init()
	if err != nil {
		return nil, err
	}
	return &vi, nil
}
