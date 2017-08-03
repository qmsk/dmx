FROM ubuntu:xenial

RUN apt-get update && apt-get install -y \
  git \
  golang-go \
  nodejs nodejs-legacy npm

ENV GOPATH=/go
RUN mkdir -p /go/src/github.com/qmsk/dmx
ADD web/package.json /go/src/github.com/qmsk/dmx/web/

WORKDIR /go/src/github.com/qmsk/dmx/web
RUN npm install

ADD . /go/src/github.com/qmsk/dmx
WORKDIR /go/src/github.com/qmsk/dmx
RUN go get -d ./cmd/... && go install -v ./cmd/...

WORKDIR /go/src/github.com/qmsk/dmx/web
RUN ./node_modules/typescript/bin/tsc

WORKDIR /go/src/github.com/qmsk/dmx
ENV ARTNET_DISCOVERY=2.255.255.255
CMD ["/go/bin/qmsk-dmx", \
  "--log=info", \
  "--http-listen=:8000", "--http-static=web/", \
  "--heads-library=library", "config/" \
]

EXPOSE 8000
