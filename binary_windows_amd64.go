//go:build standalone && windows && amd64

package kookvoice

import _ "embed"

//go:embed ffmpeg/win32-x64
var memoryBinary []byte
