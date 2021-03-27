windows:
	GOOS=windows \
	GOARCH=amd64 \
	go build -o ./cvpn_windows_amd64.exe ./cmd/cvpn/main.go

linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build -o ./cvpn_linux_amd64 ./cmd/cvpn/main.go 

darwin:
	GOOS=darwin \
	GOARCH=amd64 \
	go build -o ./cvpn_darwin_amd64 ./cmd/cvpn/main.go

