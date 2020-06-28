FROM golang:1.14

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build

EXPOSE 8080

RUN ls
CMD ["./realtimedashboard"]