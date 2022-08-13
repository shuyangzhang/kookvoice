package kookvoice

import (
	"fmt"
	"os/exec"
)

func streamAudio(rtpUrl string, audioSource string) {
	fmt.Println(">>>> start streaming <<<<")

	cmd := exec.Command(
		"ffmpeg",
		"-re",
		"-loglevel",
		"level+info",
		"-nostats",
		"-i",
		audioSource,
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
	cmd.Run()
}
