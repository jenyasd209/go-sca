FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

COPY . .

RUN go build -o main ./src/cmd

EXPOSE 8080

CMD ["./main"]
