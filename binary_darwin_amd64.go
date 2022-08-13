//go:build standalone && darwin && amd64

package kookvoice

import _ "embed"

//go:embed ffmpeg/darwin-x64
var memoryBinary []byte
