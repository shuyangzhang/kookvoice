package kookvoice

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
)

type voiceInstance struct {
	Token         string
	ChannelId     string
	wsConnect     *websocket.Conn
	streamProcess *os.Process
	sourceProcess *os.Process
}

func (i *voiceInstance) Init() error {
	makeFifoCmd := exec.Command(
		"mkfifo",
		"streampipe",
	)
	err := makeFifoCmd.Run()
	if err != nil {
		return err
	}

	keepFifoOpenCmd := exec.Command(
		"bash",
		"-c",
		"exec 7<>streampipe",
	)
	err = keepFifoOpenCmd.Run()
	if err != nil {
		return err
	}

	silentSourceCmd := exec.Command(
		"ffmpeg",
		"-f",
		"lavfi",
		"-i",
		"anullsrc",
		"-f",
		"wav",
		"-c:a",
		"pcm_s16le",
		"-b:a",
		"1411200",
		"-ar",
		"44.1k",
		"-ac",
		"2",
		"pipe:1",
		">",
		"streampipe",
	)
	err = silentSourceCmd.Start()
	if err != nil {
		return err
	}
	i.sourceProcess = silentSourceCmd.Process

	gatewayUrl := getGatewayUrl(i.Token, i.ChannelId)
	connect, rtpUrl := initWebsocketClient(gatewayUrl)

	go keepWebsocketClientAlive(connect)
	go keepRecieveMessage(connect)

	i.wsConnect = connect

	streamCmd := exec.Command(
		"ffmpeg",
		"-re",
		"-loglevel",
		"level+info",
		"-nostats",
		"-i",
		"streampipe",
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
		fmt.Sprintf("[select=a:f=rtp:ssrc=1357:payload_type=100]%v", rtpUrl),
	)
	err = streamCmd.Start()
	if err != nil {
		return err
	}
	i.streamProcess = streamCmd.Process

	return nil
}

func (i *voiceInstance) PlayMusic(input string) error {
	if err := i.sourceProcess.Kill(); err != nil {
		return err
	}

	musicSourceCmd := exec.Command(
		"ffmpeg",
		"-re",
		"-i",
		input,
		"-f",
		"s16le",
		"-c:a",
		"pcm_s16le",
		"-b:a",
		"1411200",
		"-ar",
		"44.1k",
		"-ac",
		"2",
		"pipe:1",
		">",
		"streampipe",
	)
	err := musicSourceCmd.Start()
	if err != nil {
		return err
	}
	i.sourceProcess = musicSourceCmd.Process

	err = musicSourceCmd.Wait()
	if err != nil {
		return err
	}

	silentSourceCmd := exec.Command(
		"ffmpeg",
		"-f",
		"lavfi",
		"-i",
		"anullsrc",
		"-f",
		"wav",
		"-c:a",
		"pcm_s16le",
		"-b:a",
		"1411200",
		"-ar",
		"44.1k",
		"-ac",
		"2",
		"pipe:1",
		">",
		"streampipe",
	)
	err = silentSourceCmd.Start()
	if err != nil {
		return err
	}
	i.sourceProcess = silentSourceCmd.Process

	return nil
}
