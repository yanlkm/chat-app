FROM golang:1.18.0

WORKDIR /go/src/app

COPY . .

RUN go build -o chat-app .

EXPOSE 8080

CMD ["./chat-app"]

LABEL authors="yan"
