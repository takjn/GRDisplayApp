build:
	go build -o bin/linux/displayapp cmd/displayapp/main.go

build_windows:
	GOOS=windows GOARCH=amd64 go build -o bin/win/displayapp.exe cmd/displayapp/main.go
