FROM golang:latest
MAINTAINER angdev(me@angdev.io)

WORKDIR /go/src/github.com/angdev/chocolat

# Install Godep
RUN go get github.com/tools/godep

# Install Dependencies
ADD ./Godeps/Godeps.json ./Godeps/Godeps.json
RUN godep restore

# Build
ADD . .
RUN go get .
RUN go build

# Install
RUN go install
