

FROM golang:1.19.4-alpine3.17

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY main.go ./

RUN go build -o /jitD

EXPOSE 3000

CMD [ "/jitD" ]

