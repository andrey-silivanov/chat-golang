FROM golang:1.17-alpine

RUN apk update && apk add bash

COPY ./ /app

WORKDIR /app

CMD go mod download

RUN go get github.com/githubnemo/CompileDaemon

#CMD make build

ENTRYPOINT CompileDaemon --build="go build -o ./bin/mychat cmd/myChat/main.go" --command=./bin/mychat