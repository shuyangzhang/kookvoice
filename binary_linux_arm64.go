//go:build standalone && linux && arm64

package kookvoice

import _ "embed"

//go:embed ffmpeg/linux-arm64
var memoryBinary []byte
