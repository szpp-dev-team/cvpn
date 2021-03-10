windows:
	GOOS=windows \
	GOARCH=amd64 \
	go build ./cmd/cvpn/main.go -o ./cvpn_windows_amd64.exe

linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build ./cmd/cvpn/main.go -o ./cvpn_linux_amd64