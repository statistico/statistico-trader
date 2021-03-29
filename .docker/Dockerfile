# Step 1
FROM golang:1.13.9-alpine AS builder
RUN apk update && apk add git gcc g++
ENV GO111MODULE=on
WORKDIR /app
COPY . .

RUN go get -u github.com/pressly/goose/cmd/goose && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/grpc

# Step 2
FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /opt
COPY ./bin ./bin
COPY ./database ./database
COPY ./.docker/wait-for-it.sh .
COPY --from=builder /go/bin/goose /usr/local/bin
COPY --from=builder /app/grpc .

CMD ["/bin/sh"]