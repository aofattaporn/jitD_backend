FROM golang:1.19-alpine as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy all the files to gp app 
COPY . .

# Build the Go app
RUN go build -o ./out/app .

# Run stage 
FROM alpine:3.13
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app .

# port container is 3000 
EXPOSE 3000

CMD ["./out/app"]