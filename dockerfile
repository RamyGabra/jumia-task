FROM golang:1.19 AS golang
RUN go version

FROM ubuntu:latest as build-env
RUN apt-get -qq update && apt-get -qq install -y ca-certificates curl git gcc build-essential 

# Set-up go
COPY --from=golang /usr/local/go/ /usr/local/go/
ENV PATH /usr/local/go/bin:$PATH
ENV GO111MODULE=on

RUN mkdir /src
WORKDIR /src

# Add modules in docker layer
COPY go.mod .
COPY go.sum .
RUN go mod download

# copy files
COPY . . 

RUN go build -o /main

EXPOSE 3333

CMD ["/main"] 