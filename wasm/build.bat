SET GOARCH=wasm
SET GOOS=js
go build -o test.wasm main.go