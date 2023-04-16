go build -o container/api.bin cmd/main.go
docker build -t blocopad:v001 container/
