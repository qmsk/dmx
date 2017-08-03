FROM ubuntu:xenial

RUN apt-get update && apt-get install -y \
  git \
  golang-go \
  nodejs nodejs-legacy npm

ENV GOPATH=/go
RUN mkdir -p /go/src/github.com/qmsk/dmx
RUN go get -u github.com/kardianos/govendor

ADD web/package.json /go/src/github.com/qmsk/dmx/web/
WORKDIR /go/src/github.com/qmsk/dmx/web
RUN npm install

ADD vendor/vendor.json /go/src/github.com/qmsk/dmx/vendor/
WORKDIR /go/src/github.com/qmsk/dmx
RUN /go/bin/govendor sync -v

ADD . /go/src/github.com/qmsk/dmx
RUN go install -v ./cmd/...

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
