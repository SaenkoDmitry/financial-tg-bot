FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

RUN ls -al

CMD ["./bin/bot"]
