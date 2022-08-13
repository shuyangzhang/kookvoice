PROJECTNAME := kookvoice

.PHONY: windows linux macos win lin mac all

all: windows linux macos

windows:
	@GOOS=windows GOARCH=amd64 go build -o output/$(PROJECTNAME)-amd64.exe cmd/kook-voice/main.go 
	@GOOS=windows GOARCH=arm64 go build -o output/$(PROJECTNAME)-arm64.exe cmd/kook-voice/main.go
	@GOOS=windows GOARCH=386 go build -o output/$(PROJECTNAME)-i386.exe cmd/kook-voice/main.go
	@GOOS=windows GOARCH=amd64 go build -tags standalone -o output/$(PROJECTNAME)-standalone-amd64.exe cmd/kook-voice/main.go
	@GOOS=windows GOARCH=386 go build -tags standalone -o output/$(PROJECTNAME)-standalone-i386.exe cmd/kook-voice/main.go

win: windows

linux:
	@GOOS=linux GOARCH=amd64 go build -o output/$(PROJECTNAME)-amd64-linux cmd/kook-voice/main.go
	@GOOS=linux GOARCH=arm64 go build -o output/$(PROJECTNAME)-arm64-linux cmd/kook-voice/main.go
	@GOOS=linux GOARCH=amd64 go build -tags standalone -o output/$(PROJECTNAME)-standalone-amd64-linux cmd/kook-voice/main.go
	@GOOS=linux GOARCH=arm64 go build -tags standalone -o output/$(PROJECTNAME)-standalone-arm64-linux cmd/kook-voice/main.go

lin: linux

macos:
	@GOOS=darwin GOARCH=amd64 go build -o output/$(PROJECTNAME)-amd64-darwin cmd/kook-voice/main.go
	@GOOS=darwin GOARCH=arm64 go build -o output/$(PROJECTNAME)-arm64-darwin cmd/kook-voice/main.go
	@GOOS=darwin GOARCH=amd64 go build -tags standalone -o output/$(PROJECTNAME)-standalone-amd64-darwin cmd/kook-voice/main.go
	@GOOS=darwin GOARCH=arm64 go build -tags standalone -o output/$(PROJECTNAME)-standalone-arm64-darwin cmd/kook-voice/main.go

mac: macos
