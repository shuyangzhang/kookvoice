package main

import (
	"flag"

	"github.com/shuyangzhang/kookvoice"
)

var (
	TOKEN      = flag.String("t", "", "bot token")
	CHANNEL_ID = flag.String("c", "", "channel id")
	INPUT      = flag.String("i", "", "input audio")
)

func init() {
	flag.Parse()
}

func main() {
	kookvoice.Play(*TOKEN, *CHANNEL_ID, *INPUT)
}
