FROM golang:1.15.7-buster

RUN mkdir /sampleapp
COPY . /sampleapp/
RUN go build -o /sampleapp/main /sampleapp/src/sample.go

ENTRYPOINT ["/sampleapp/main"]
