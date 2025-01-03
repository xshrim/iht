bin:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o iht main.go

bin-win:
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o iht.exe main.go

dev:
	env GOL_LEVEL=debug go run main.go

run:
	go run main.go