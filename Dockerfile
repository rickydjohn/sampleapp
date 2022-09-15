FROM golang:1.15.7-buster

RUN apt-get update && \
    apt-get install -y git && \
    git clone https://github.com/rickydjohn/sampleapp.git

WORKDIR sampleapp

RUN echo "blah"
RUN go build -o sampleapp src/sample.go

ENTRYPOINT ["./sampleapp"]
