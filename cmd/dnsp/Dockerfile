FROM golang:latest

RUN apt-get update
RUN apt-get install dnsutils -y
RUN mkdir /app
ADD *.go /app/
WORKDIR /app
RUN go get -u github.com/gophergala/dnsp/...
RUN go build -o /app/dnsp

EXPOSE 53

ENTRYPOINT ["/app/dnsp"]
