FROM golang:1.17 as builder

RUN apt-get update && apt-get install -y \
  build-essential \
  libsqlite3-dev \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM ubuntu:latest  
RUN apt-get update && apt-get install -y \
  git \
  openssh-client \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*
RUN echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config
RUN echo "UserKnownHostsFile /dev/null" >> /etc/ssh/ssh_config
RUN echo "IdentityFile /.ssh.identity" >> /etc/ssh/ssh_config

WORKDIR /

COPY --from=builder /go/bin/zombie-kahinah /
COPY --from=builder /go/src/app/static /static
COPY --from=builder /go/src/app/views /views

ENV HOME=/

VOLUME ["/conf", "/data", "/.ssh.identity", "/news.txt"]
CMD ["/zombie-kahinah"]
