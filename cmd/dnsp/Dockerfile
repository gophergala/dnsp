FROM golang:latest

RUN mkdir /app
ADD *.go /app/
WORKDIR /app
RUN go get -u github.com/gophergala/dnsp/...
RUN go build -o dnsp

EXPOSE 53

ENTRYPOINT ["/app/dnsp"]
