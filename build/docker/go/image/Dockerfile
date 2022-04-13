FROM golang:1.18-alpine

RUN apk add git gcc g++
RUN apk update && apk add bash

COPY ./ /app

WORKDIR /app

CMD go mod download

RUN go install github.com/githubnemo/CompileDaemon@latest

#CMD make build

ENTRYPOINT CompileDaemon --build="go build -o ./bin/mychat cmd/myChat/main.go" --command=./bin/mychat