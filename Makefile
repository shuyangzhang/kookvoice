PROJECTNAME := kookvoice

.PHONY: windows linux macos win lin mac all

all: windows linux macos

windows:
	@GOOS=windows GOARCH=amd64 go build -o $(PROJECTNAME)-amd64.exe .
	@GOOS=windows GOARCH=arm64 go build -o $(PROJECTNAME)-arm64.exe .

win: windows

linux:
	@GOOS=linux GOARCH=amd64 go build -o $(PROJECTNAME)-amd64-linux .
	@GOOS=linux GOARCH=arm64 go build -o $(PROJECTNAME)-arm64-linux .

lin: linux

macos:
	@GOOS=darwin GOARCH=amd64 go build -o $(PROJECTNAME)-amd64-darwin .
	@GOOS=darwin GOARCH=arm64 go build -o $(PROJECTNAME)-arm64-darwin .

mac: macos
