FROM golang:1.26.4-alpine3.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o ./app-bin cmd/app/main.go

CMD [ "./app-bin" ]