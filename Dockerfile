FROM golang:1.19

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go env -w GOPROXY=https://proxy.golang.com.cn,direct \
    & go build -v -o /usr/local/bin/app ./main.go

CMD ["app"]