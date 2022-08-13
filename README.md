# kookvoice
kook voice golang library and cli tool.

## Statement
All APIs in this project are derived from packet crawling, and these APIs may be invalidated by version updates.  
You need to know that using this API violates [KOOK Voice Software License and Service Agreement](https://www.kookapp.cn/protocol.html) `3.2.3` or `3.2.5`.  
It also violates the terms of the [KOOK Developer Privacy Policy](https://developer.kookapp.cn/doc/privacy) `Data Information` or `Abuse`.  

## Go Library

### Dependencies

Go version >= 1.18

### Installation

```bash
go get github.com/shuyangzhang/kookvoice
```

### Usage

```go
package main

import (
    "github.com/shuyangzhang/kookvoice"
)

func main() {
    token := "1/MECxOTk=/zCX2VjWr6p+AmD84jL9asQ=="
    channel := "2559449076697969"
    input := "./test.mp3"   // Local audio path or network audio url are both valid.
    
    kookvoice.Play(token, channel, input)
}
```

## CLI Tool

### Download Binary
Go to [the release page](https://github.com/shuyangzhang/kookvoice/releases) to download the binary that matches your operating system

### Usage
use `-h` to get help message
```bash
./kookvoice-amd64-linux -h
```
```bash
Usage of ./kookvoice-amd64-linux:
  -c string
        channel id
  -i string
        input audio
  -t string
        bot token
```

binary file without `-standalone-` tag needs `ffmpeg` installed in your `PATH`
```bash
ffmpeg -version
```
```bash
./kookvoice-amd64-linux -t ${YOUR_TOKEN} -c ${YOUR_CHANNEL_ID} -i ${AUDIO_INPUT_URL_OR_PATH}
```

binary file with `-standalone-` tag can run without `ffmpeg`
```bash
./kookvoice-standalone-amd64-linux -t ${YOUR_TOKEN} -c ${YOUR_CHANNEL_ID} -i ${AUDIO_INPUT_URL_OR_PATH}
```

## License
This project is licensed under the terms of the [MIT license](./LICENSE).  
The binary release with `-standalone-` tag is licensed under terms of the GPL-3.0 license.

## Credits
This project is inspired by [kook-voice-API](https://github.com/hank9999/kook-voice-API).  
The standalone binaries are depends on [ffmpeg-static](https://github.com/eugeneware/ffmpeg-static).
