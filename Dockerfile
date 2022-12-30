FROM golang

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o /jitD

EXPOSE 3000

CMD [ "/jitD" ]

