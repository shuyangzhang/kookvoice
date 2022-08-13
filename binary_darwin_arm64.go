//go:build standalone && darwin && arm64

package kookvoice

import _ "embed"

//go:embed ffmpeg/darwin-arm64
var memoryBinary []byte
