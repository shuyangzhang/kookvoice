//go:build standalone && windows && 386

package kookvoice

import _ "embed"

//go:embed ffmpeg/win32-ia32
var memoryBinary []byte
