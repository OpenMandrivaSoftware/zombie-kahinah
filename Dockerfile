FROM golang:1.17-alpine as builder

RUN apk add sqlite-dev

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:latest  
RUN apk --no-cache add ca-certificates git

WORKDIR /

COPY --from=builder /go/bin/zombie-kahinah /
VOLUME ["/conf", "/data"]
CMD ["/zombie-kahinah"]
