//go:build standalone

package kookvoice

import (
	_ "embed"
	"fmt"

	"github.com/amenzhinsky/go-memexec"
)

func streamAudio(rtpUrl string, audioSource string) {
	exe, err := memexec.New(memoryBinary)
	if err != nil {
		panic(err)
	}
	defer exe.Close()

	fmt.Println(">>>> start streaming <<<<")
	cmd := exe.Command(
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
