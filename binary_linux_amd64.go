//go:build standalone && linux && amd64

package kookvoice

import _ "embed"

//go:embed ffmpeg/linux-x64
var memoryBinary []byte
