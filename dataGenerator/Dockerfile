FROM golang:1.17-alpine

WORKDIR /app
RUN apk add bash
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /app
RUN go build -o /bin/dataGenerator cmd/main.go

ENTRYPOINT ["/bin/dataGenerator", "10", "20", "rawdata"]